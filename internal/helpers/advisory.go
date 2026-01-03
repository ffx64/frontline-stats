package helpers

import (
	"context"
	"log"

	"gorm.io/gorm"
)

func AdvisoryLock(ctx context.Context, db *gorm.DB, key string) (*gorm.DB, error) {
	tx := db.WithContext(ctx)

	query := `SELECT pg_try_advisory_lock(hashtext(?))`

	var hasLock bool
	if err := tx.Raw(query, key).Scan(&hasLock).Error; err != nil {
		log.Printf("[helpers:advisory_lock] erro ao tentar adquirir lock: %v", err)
		return nil, err
	}

	if !hasLock {
		log.Println("[helpers:advisory_lock] outra instância já está executando")
		return nil, nil
	}

	return tx, nil
}

func AdvisoryUnlock(ctx context.Context, tx *gorm.DB, key string) {
	query := `SELECT pg_advisory_unlock(hashtext(?))`
	if err := tx.Exec(query, key).Error; err != nil {
		log.Printf("[helpers:advisory_unlock] erro ao liberar lock: %v", err)
	}
}
