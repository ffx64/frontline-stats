package scheduler

import (
	"log"
	"time"

	"github.com/ffx64/frontline-stats/internal/repositories"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron             *cron.Cron
	playersStatsRepo repositories.PlayersStatsRepository
	roundsStatsRepo  repositories.RoundsStatsRepository
	rdb              *redis.Client
}

func NewScheduler(playersStatsRepo repositories.PlayersStatsRepository, roundsStatsRepo repositories.RoundsStatsRepository, rdb *redis.Client) *Scheduler {
	return &Scheduler{
		cron:             cron.New(cron.WithSeconds()),
		playersStatsRepo: playersStatsRepo,
		roundsStatsRepo:  roundsStatsRepo,
		rdb:              rdb,
	}
}

func (s *Scheduler) Start() {
	if _, err := s.cron.AddFunc("0 */15 * * * *", s.updatePlayerStats); err != nil {
		log.Fatalf("[scheduler] failed to register job: %v", err)
	}

	if _, err := s.cron.AddFunc("0 */5 * * * *", s.updateRoundsStats); err != nil {
		log.Fatalf("[scheduler] failed to register job: %v", err)
	}

	s.cron.Start()
	log.Println("[scheduler] scheduler started")
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	select {
	case <-ctx.Done():
		log.Println("[scheduler] scheduler stopped")
	case <-time.After(5 * time.Second):
		log.Println("[scheduler] timeout stopping scheduler")
	}
}
