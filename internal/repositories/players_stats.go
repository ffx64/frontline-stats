package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/ffx64/gamestats-backend/internal/entities"
	"github.com/ffx64/gamestats-backend/internal/helpers"
	"gorm.io/gorm"
)

type PlayersStatsRepository interface {
	Save(ctx context.Context, stats *entities.PlayersStats) error
	UpdatePlayerStats(ctx context.Context) error
	FindByPlayerID(ctx context.Context, playerID string) (*entities.PlayersStats, error)
	GetLeaderboard(ctx context.Context) ([]*entities.PlayersStats, error)
	GetLeaderboardHeadshots(ctx context.Context) ([]*entities.PlayersStats, error)
	GetLeaderboardVehicle(ctx context.Context) ([]*entities.PlayersStats, error)
}

type playersStatsRepository struct {
	db *gorm.DB
}

func NewPlayersStatsRepository(db *gorm.DB) PlayersStatsRepository {
	return &playersStatsRepository{db: db}
}

func (r *playersStatsRepository) FindByPlayerID(ctx context.Context, playerID string) (*entities.PlayersStats, error) {
	var stats entities.PlayersStats
	if err := r.db.WithContext(ctx).Where("player_id = ?", playerID).First(&stats).Error; err != nil {
		log.Printf("[repository:stats] erro ao buscar estatísticas do jogador %s: %v", playerID, err)
		return nil, fmt.Errorf("erro ao buscar estatísticas do jogador: %w", err)
	}
	log.Printf("[repository:stats] estatísticas recuperadas para jogador: %s", playerID)
	return &stats, nil
}

func (r *playersStatsRepository) Save(ctx context.Context, stats *entities.PlayersStats) error {
	if err := r.db.WithContext(ctx).Create(stats).Error; err != nil {
		log.Printf("[repository:stats] erro ao criar estatística do jogador: %v", err)
		return fmt.Errorf("erro ao criar estátistica do jogador: %w", err)
	}
	log.Printf("[repository:stats] estatística criada para jogador: %v", stats.PlayerID)
	return nil
}

func (r *playersStatsRepository) UpdatePlayerStats(ctx context.Context) error {
	tx, err := helpers.AdvisoryLock(ctx, r.db, "update_player_stats_key")

	if err != nil {
		log.Printf("[repository:stats] erro ao tentar adquirir lock: %v", err)
		return err
	}

	if tx == nil {
		log.Printf("[repository:stats] outra instância já está executando")
		return err
	}

	defer helpers.AdvisoryUnlock(ctx, tx, "update_player_stats_key")

	query := `
WITH killer_stats AS (
  SELECT
    k.killer_id AS player_id,
    COUNT(*) FILTER (WHERE NOT k.is_friendly) AS kills,
    COUNT(*) FILTER (WHERE k.is_friendly) AS friendly_fire_kills,
    COUNT(*) FILTER (WHERE k.is_vehicle AND NOT k.is_friendly) AS vehicle_kills,
    MAX(k.distance) FILTER (WHERE NOT k.is_friendly AND NOT k.is_vehicle) AS longest_kill_distance,
    AVG(k.distance) FILTER (WHERE NOT k.is_friendly AND NOT k.is_vehicle) AS average_kill_distance,
    COUNT(*) FILTER (
      WHERE k.killer_weapon_name ILIKE ANY (ARRAY['%grenade%', '%stun grenade%'])
    ) AS grenades_thrown,
    COUNT(*) FILTER (WHERE k.is_headshot AND NOT k.is_friendly) AS headshots_made
  FROM kills k
  WHERE k.killer_id IN (SELECT player_id FROM players_stats)
  GROUP BY k.killer_id
),
victim_stats AS (
  SELECT
    k.victim_id AS player_id,
    COUNT(*) AS deaths,
    COUNT(*) FILTER (WHERE k.is_friendly) AS friendly_fire_deaths,
    COUNT(*) FILTER (WHERE k.is_vehicle) AS vehicle_deaths,
    COUNT(*) FILTER (WHERE k.is_headshot) AS headshots_taken,
    AVG(k.distance) AS average_death_distance
  FROM kills k
  WHERE k.victim_id IN (SELECT player_id FROM players_stats)
  GROUP BY k.victim_id
)
UPDATE players_stats ps
SET
  kills = COALESCE(ks.kills, 0),
  deaths = COALESCE(vs.deaths, 0),
  friendly_fire_kills = COALESCE(ks.friendly_fire_kills, 0),
  friendly_fire_deaths = COALESCE(vs.friendly_fire_deaths, 0),
  headshots_made = COALESCE(ks.headshots_made, 0),
  headshots_taken = COALESCE(vs.headshots_taken, 0),
  vehicle_kills = COALESCE(ks.vehicle_kills, 0),
  vehicle_deaths = COALESCE(vs.vehicle_deaths, 0),
  grenades_thrown = COALESCE(ks.grenades_thrown, 0),
  longest_kill_distance = COALESCE(ks.longest_kill_distance, 0),
  average_kill_distance = COALESCE(ks.average_kill_distance, 0),
  average_death_distance = COALESCE(vs.average_death_distance, 0),
  updated_at = NOW()
FROM killer_stats ks
FULL OUTER JOIN victim_stats vs ON ks.player_id = vs.player_id
WHERE ps.player_id = COALESCE(ks.player_id, vs.player_id);
`

	if err := r.db.WithContext(ctx).Exec(query).Error; err != nil {
		log.Printf("[repository:stats] erro ao atualizar estatísticas principais: %v", err)
		return err
	}

	query = `
UPDATE players_stats
SET
  ratio_kdr = CASE WHEN deaths > 0 THEN kills::decimal / deaths ELSE kills END,
  ratio_headshot = CASE WHEN kills > 0 THEN headshots_made::decimal / kills ELSE 0 END,
  ratio_friendly_fire = CASE WHEN kills > 0 THEN friendly_fire_kills::decimal / kills ELSE 0 END,
  ratio_vehicle = CASE WHEN vehicle_deaths > 0 THEN vehicle_kills::decimal / vehicle_deaths ELSE 0 END
WHERE TRUE;
`
	// Calcula os ratios separadamente
	if err := r.db.WithContext(ctx).Exec(query).Error; err != nil {
		log.Printf("[repository:stats] erro ao atualizar ratios: %v", err)
		return err
	}

	queries := []string{
		`UPDATE players_stats ps SET weapons_most_used = sub.weapon
			FROM (
			  SELECT DISTINCT ON (killer_id)
			    killer_id,
			    killer_weapon_name AS weapon,
			    COUNT(*) OVER (PARTITION BY killer_id, killer_weapon_name) AS total
			  FROM kills
			  WHERE killer_weapon_name IS NOT NULL
			    AND NOT is_vehicle
				AND NOT is_friendly
				AND killer_weapon_name NOT ILIKE ANY (ARRAY[
				'%unknown%', '%vehicle%', '%none%', '%unarmed%'])
		  	  ORDER BY killer_id, total DESC, killer_weapon_name ASC
			) AS sub
			WHERE ps.player_id = sub.killer_id;`,

		`UPDATE players_stats ps SET vehicle_most_used = sub.vehicle
			FROM (
			  SELECT DISTINCT ON (killer_id)
			    killer_id,
			    killer_weapon_name AS vehicle,
			    COUNT(*) OVER (PARTITION BY killer_id, killer_weapon_name) AS total
			  FROM kills
			  WHERE is_vehicle
			  	AND NOT is_friendly
			  ORDER BY killer_id, total DESC, killer_weapon_name ASC
			) AS sub
			WHERE ps.player_id = sub.killer_id;`,

		`UPDATE players_stats ps SET hit_zones_most_killed = sub.zone
			FROM (
			  SELECT DISTINCT ON (killer_id)
			    killer_id,
			    hit_zone AS zone,
			    COUNT(*) OVER (PARTITION BY killer_id, hit_zone) AS total
			  FROM kills
			  WHERE killer_id IS NOT NULL
			    AND hit_zone IS NOT NULL
				AND NOT is_friendly
			  ORDER BY killer_id, total DESC, hit_zone ASC
			) AS sub
			WHERE ps.player_id = sub.killer_id;`,

		`UPDATE players_stats ps SET hit_zones_most_died = sub.zone
			FROM (
			  SELECT DISTINCT ON (victim_id)
			    victim_id,
			    hit_zone AS zone,
			    COUNT(*) OVER (PARTITION BY victim_id, hit_zone) AS total
			  FROM kills
			  WHERE victim_id IS NOT NULL
			    AND hit_zone IS NOT NULL
				AND NOT is_friendly
			  ORDER BY victim_id, total DESC, hit_zone ASC
			) AS sub
			WHERE ps.player_id = sub.victim_id;`,
	}

	for _, q := range queries {
		if err := r.db.WithContext(ctx).Exec(q).Error; err != nil {
			log.Printf("[repository:stats] erro ao atualizar mosts: %v", err)
			return err
		}
	}

	log.Println("[repository:stats] estatísticas dos jogadores atualizadas com sucesso")
	return nil
}

func (r *playersStatsRepository) GetLeaderboard(ctx context.Context) ([]*entities.PlayersStats, error) {
	var stats []*entities.PlayersStats

	if err := r.db.WithContext(ctx).Table("players_stats").Order("(kills - vehicle_kills) DESC").Limit(20).Find(&stats).Error; err != nil {
		log.Printf("[repository:stats] erro ao obter leaderboard: %v", err)
		return nil, err
	}

	return stats, nil
}

func (r *playersStatsRepository) GetLeaderboardHeadshots(ctx context.Context) ([]*entities.PlayersStats, error) {
	var stats []*entities.PlayersStats

	if err := r.db.WithContext(ctx).Table("players_stats").Order("headshots_made DESC").Limit(20).Find(&stats).Error; err != nil {
		log.Printf("[repository:stats] erro ao obter leaderboard de headshots: %v", err)
		return nil, err
	}

	return stats, nil
}

func (r *playersStatsRepository) GetLeaderboardVehicle(ctx context.Context) ([]*entities.PlayersStats, error) {
	var stats []*entities.PlayersStats

	if err := r.db.WithContext(ctx).Table("players_stats").Order("vehicle_kills DESC").Limit(20).Find(&stats).Error; err != nil {
		log.Printf("[repository:stats] erro ao obter leaderboard de veículos destruídos: %v", err)
		return nil, err
	}

	return stats, nil
}
