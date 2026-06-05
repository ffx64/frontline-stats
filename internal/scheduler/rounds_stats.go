package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/ffx64/frontline-stats/internal/cache"
)

func (s *Scheduler) updateRoundsStats() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	log.Println("[scheduler:stats] updating round stats...")
	if err := s.roundsStatsRepo.UpdateRoundsStats(ctx); err != nil {
		log.Printf("[scheduler:stats] failed to update round stats: %v", err)
		return
	}

	cache.DeletePattern(ctx, s.rdb, "round:scoreboard:*")
	log.Println("[scheduler:stats] round stats updated successfully")
}
