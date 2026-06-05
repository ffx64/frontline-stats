package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ffx64/frontline-stats/internal/entities"
	"github.com/ffx64/frontline-stats/internal/helpers"
)

type RoundsStatsRepository interface {
	Save(ctx context.Context, stats *entities.RoundsStats) (*entities.RoundsStats, error)
	FindScoreboardByRoundID(ctx context.Context, roundId uuid.UUID) ([]entities.RoundsStats, error)
	UpdateRoundsStats(ctx context.Context) error
	Update(ctx context.Context, stats *entities.RoundsStats, roundId uuid.UUID, playerId uuid.UUID) (*entities.RoundsStats, error)
}

type roundsStatsRepository struct {
	db *gorm.DB
}

func NewRoundsStatsRepository(db *gorm.DB) RoundsStatsRepository {
	return &roundsStatsRepository{db: db}
}

func (r *roundsStatsRepository) Save(ctx context.Context, stats *entities.RoundsStats) (*entities.RoundsStats, error) {
	tx := r.db.WithContext(ctx).Create(stats)
	if tx.Error != nil {
		log.Printf("[repository:rounds_stats] failed to create round stats: %v", tx.Error)
		return nil, fmt.Errorf("failed to create round stats: %w", tx.Error)
	}
	log.Printf("[repository:rounds_stats] round stats created: %v", stats.ID)
	return stats, nil
}

func (r *roundsStatsRepository) FindScoreboardByRoundID(ctx context.Context, roundId uuid.UUID) ([]entities.RoundsStats, error) {
	var stats []entities.RoundsStats
	tx := r.db.WithContext(ctx).Where("round_id = ?", roundId).Find(&stats)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			log.Printf("[repository:rounds_stats] no round stats found for round_id: %v", roundId)
			return nil, nil
		}
		log.Printf("[repository:rounds_stats] failed to find round stats by round_id %v: %v", roundId, tx.Error)
		return nil, fmt.Errorf("failed to find round stats by round_id: %w", tx.Error)
	}
	log.Printf("[repository:rounds_stats] round stats retrieved: %v", roundId)
	return stats, nil
}

func (r *roundsStatsRepository) UpdateRoundsStats(ctx context.Context) error {
	tx, err := helpers.AdvisoryLock(ctx, r.db, "update_rounds_stats_key")
	if err != nil {
		log.Printf("[repository:rounds_stats] failed to acquire advisory lock: %v", err)
		return err
	}
	if tx == nil {
		log.Printf("[repository:rounds_stats] another instance is already running, skipping update")
		return nil
	}

	query := `
WITH active_rounds AS (
    SELECT
        id AS round_id,
        server_id
    FROM rounds
    WHERE status = 'in_progress'
),

round_kills AS (
    SELECT
        k.*,
        ar.round_id AS r_id,
        ar.server_id AS s_id
    FROM kills k
    INNER JOIN active_rounds ar
        ON ar.round_id = k.round_id
        AND ar.server_id = k.server_id
),

killer_stats AS (
    SELECT
        killer_id AS player_id,
        ar.round_id AS r_id,
        ar.server_id AS s_id,
        COUNT(*) FILTER (WHERE NOT is_friendly AND killer_id != victim_id) AS total_kills,
        COUNT(*) FILTER (WHERE is_headshot = TRUE AND killer_id != victim_id) AS total_headshots,
        COUNT(*) FILTER (WHERE is_vehicle = TRUE AND killer_id != victim_id) AS total_vehicle_kills,
        COUNT(*) FILTER (WHERE is_friendly = TRUE AND killer_id != victim_id) AS total_team_kills,
        SUM(distance) FILTER (WHERE NOT is_friendly AND killer_id != victim_id) AS total_kill_distance,
        killer_team AS team
    FROM round_kills k
    INNER JOIN active_rounds ar
        ON ar.round_id = k.round_id
        AND ar.server_id = k.server_id
    GROUP BY killer_id, killer_team, ar.round_id, ar.server_id
),

victim_stats AS (
    SELECT
        victim_id AS player_id,
        ar.round_id AS r_id,
        ar.server_id AS s_id,
        COUNT(*) AS total_deaths,
        COUNT(*) FILTER (WHERE killer_id = victim_id) AS total_suicides,
        COUNT(*) FILTER (WHERE killer_id = victim_id AND is_vehicle = TRUE) AS total_vehicle_deaths,
        victim_team AS team
    FROM round_kills k
    INNER JOIN active_rounds ar
        ON ar.round_id = k.round_id
        AND ar.server_id = k.server_id
    GROUP BY victim_id, victim_team, ar.round_id, ar.server_id
),

killer_weapon_count AS (
    SELECT
        killer_id AS player_id,
        killer_weapon_name,
        COUNT(*) AS cnt,
        ROW_NUMBER() OVER (PARTITION BY killer_id ORDER BY COUNT(*) DESC) AS rn
    FROM round_kills
    GROUP BY killer_id, killer_weapon_name
),

best_weapon AS (
    SELECT player_id, killer_weapon_name AS most_used_weapon
    FROM killer_weapon_count
    WHERE rn = 1
),

killer_hitzone_count AS (
    SELECT
        killer_id AS player_id,
        hit_zone,
        COUNT(*) AS cnt,
        ROW_NUMBER() OVER (PARTITION BY killer_id ORDER BY COUNT(*) DESC) AS rn
    FROM round_kills
    GROUP BY killer_id, hit_zone
),

best_hitzone AS (
    SELECT player_id, hit_zone AS most_hit_zone
    FROM killer_hitzone_count
    WHERE rn = 1
),

merged AS (
    SELECT DISTINCT
        COALESCE(k.player_id, v.player_id) AS player_id,
        COALESCE(k.r_id, v.r_id) AS round_id,
        COALESCE(k.s_id, v.s_id) AS server_id,
        COALESCE(k.team, v.team, '') AS team,
        COALESCE(k.total_kills, 0) AS total_kills,
        COALESCE(v.total_deaths, 0) AS total_deaths,
        COALESCE(v.total_suicides, 0) AS total_suicides,
        COALESCE(k.total_team_kills, 0) AS total_team_kills,
        COALESCE(k.total_headshots, 0) AS total_headshots,
        COALESCE(k.total_vehicle_kills, 0) AS total_vehicle_kills,
        COALESCE(v.total_vehicle_deaths, 0) AS total_vehicle_deaths,
        CASE WHEN k.total_kills > 0
            THEN k.total_kill_distance / k.total_kills
            ELSE 0
        END AS average_kill_distance,
        COALESCE(w.most_used_weapon, '') AS most_used_weapon,
        COALESCE(h.most_hit_zone, '') AS most_hit_zone,
        ROW_NUMBER() OVER (PARTITION BY COALESCE(k.player_id, v.player_id), COALESCE(k.r_id, v.r_id) ORDER BY COALESCE(k.total_kills, 0) DESC) AS rn
    FROM killer_stats k
    FULL OUTER JOIN victim_stats v ON k.player_id = v.player_id
    LEFT JOIN best_weapon w ON w.player_id = COALESCE(k.player_id, v.player_id)
    LEFT JOIN best_hitzone h ON h.player_id = COALESCE(k.player_id, v.player_id)
)

INSERT INTO rounds_stats (
    round_id, server_id, player_id, team,
    total_kills, total_deaths, total_suicides,
    average_kill_distance, total_team_kills,
    total_headshots, total_vehicle_kills, total_vehicle_deaths,
    most_used_weapon, most_hit_zone,
    created_at, updated_at
)
SELECT
    round_id, server_id, player_id, team,
    total_kills, total_deaths, total_suicides,
    average_kill_distance, total_team_kills,
    total_headshots, total_vehicle_kills, total_vehicle_deaths,
    most_used_weapon, most_hit_zone,
    NOW(),
    NOW()
FROM merged
WHERE rn = 1
ON CONFLICT (round_id, player_id) DO UPDATE SET
    team                = EXCLUDED.team,
    total_kills         = EXCLUDED.total_kills,
    total_deaths        = EXCLUDED.total_deaths,
    total_suicides      = EXCLUDED.total_suicides,
    average_kill_distance = EXCLUDED.average_kill_distance,
    total_team_kills    = EXCLUDED.total_team_kills,
    total_headshots     = EXCLUDED.total_headshots,
    total_vehicle_kills = EXCLUDED.total_vehicle_kills,
    total_vehicle_deaths = EXCLUDED.total_vehicle_deaths,
    most_used_weapon    = EXCLUDED.most_used_weapon,
    most_hit_zone       = EXCLUDED.most_hit_zone,
    updated_at          = EXCLUDED.updated_at;
`

	result := tx.WithContext(ctx).Exec(query)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("[repository:rounds_stats] failed to update in-progress round stats: %v", result.Error)
		return result.Error
	}

	if err := helpers.AdvisoryUnlock(ctx, tx, "update_rounds_stats_key"); err != nil {
		log.Printf("[repository:rounds_stats] failed to commit rounds_stats update: %v", err)
		return err
	}

	log.Printf("[repository:rounds_stats] in-progress round stats updated successfully, %d rows affected", result.RowsAffected)
	return nil
}

func (r *roundsStatsRepository) Update(ctx context.Context, stats *entities.RoundsStats, roundId uuid.UUID, playerId uuid.UUID) (*entities.RoundsStats, error) {
	tx := r.db.WithContext(ctx).Model(&entities.RoundsStats{}).Where("round_id = ? AND player_id = ?", roundId, playerId).Updates(stats)
	if tx.Error != nil {
		log.Printf("[repository:rounds_stats] failed to update round stats for round %v player %v: %v", roundId, playerId, tx.Error)
		return nil, fmt.Errorf("failed to update round stats: %w", tx.Error)
	}
	log.Printf("[repository:rounds_stats] stats updated for player %v in round %v", playerId, roundId)
	return stats, nil
}
