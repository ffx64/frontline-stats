package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/ffx64/gamestats-backend/internal/cache"
)

func (s *Scheduler) updatePlayerStats() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	log.Println("[scheduler:stats] updating player stats...")
	if err := s.playersStatsRepo.UpdatePlayerStats(ctx); err != nil {
		log.Printf("[scheduler:stats] failed to update player stats: %v", err)
		return
	}

	cache.Delete(ctx, s.rdb,
		cache.KeyLeaderboardKills,
		cache.KeyLeaderboardHeadshots,
		cache.KeyLeaderboardVehicles,
	)
	log.Println("[scheduler:stats] player stats updated successfully")
}
