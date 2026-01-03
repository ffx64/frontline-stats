package scheduler

import (
	"context"
	"log"
	"time"
)

func (s *Scheduler) updatePlayerStats() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	log.Println("[scheduler:stats] atualizando estatísticas dos jogadores...")
	if err := s.playersStatsRepo.UpdatePlayerStats(ctx); err != nil {
		log.Printf("[scheduler:stats] erro ao atualizar stats: %v", err)
		return
	}
	log.Println("[scheduler:stats] estatísticas atualizadas com sucesso")
}
