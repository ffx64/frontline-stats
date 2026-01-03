package scheduler

import (
	"log"
	"time"

	"github.com/ffx64/gamestats-backend/internal/repositories"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron             *cron.Cron
	playersStatsRepo repositories.PlayersStatsRepository
	roundsStatsRepo  repositories.RoundsStatsRepository
}

func NewScheduler(playersStatsRepo repositories.PlayersStatsRepository, roundsStatsRepo repositories.RoundsStatsRepository) *Scheduler {
	return &Scheduler{
		cron:             cron.New(cron.WithSeconds()),
		playersStatsRepo: playersStatsRepo,
		roundsStatsRepo:  roundsStatsRepo,
	}
}

func (s *Scheduler) Start() {
	if _, err := s.cron.AddFunc("0 */15 * * * *", s.updatePlayerStats); err != nil {
		log.Fatalf("[scheduler] falha ao registrar job: %v", err)
	}

	if _, err := s.cron.AddFunc("0 */5 * * * *", s.updateRoundsStats); err != nil {
		log.Fatalf("[scheduler] falha ao registrar job: %v", err)
	}

	s.cron.Start()
	log.Println("[scheduler] scheduler iniciado com sucesso")
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	select {
	case <-ctx.Done():
		log.Println("[scheduler] scheduler parado com sucesso")
	case <-time.After(5 * time.Second):
		log.Println("[scheduler] timeout ao parar scheduler")
	}
}
