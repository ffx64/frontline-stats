package scheduler

import (
	"context"
	"log"
	"time"
)

func (s *Scheduler) updateRoundsStats() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	log.Println("[scheduler:stats] atualizando estatísticas das rodadas...")
	if err := s.roundsStatsRepo.UpdateRoundsStats(ctx); err != nil {
		log.Printf("[scheduler:stats] erro ao atualizar stats das rodadas: %v", err)
		return
	}
	log.Println("[scheduler:stats] estatísticas das rodadas atualizadas com sucesso")
}
