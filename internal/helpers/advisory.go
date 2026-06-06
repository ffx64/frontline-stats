package helpers

import (
	"context"
	"log"

	"gorm.io/gorm"
)

func AdvisoryLock(ctx context.Context, db *gorm.DB, key string) (*gorm.DB, error) {
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		log.Printf("[helpers:advisory_lock] failed to begin transaction: %v", tx.Error)
		return nil, tx.Error
	}

	var hasLock bool
	if err := tx.Raw("SELECT pg_try_advisory_xact_lock(hashtext(?))", key).Scan(&hasLock).Error; err != nil {
		tx.Rollback()
		log.Printf("[helpers:advisory_lock] failed to acquire lock: %v", err)
		return nil, err
	}

	if !hasLock {
		tx.Rollback()
		log.Println("[helpers:advisory_lock] another instance is already running")
		return nil, nil
	}

	return tx, nil
}

// AdvisoryUnlock commits the transaction, automatically releasing the advisory lock.
func AdvisoryUnlock(ctx context.Context, tx *gorm.DB, key string) error {
	if err := tx.Commit().Error; err != nil {
		log.Printf("[helpers:advisory_unlock] failed to commit transaction: %v", err)
		tx.Rollback()
		return err
	}
	return nil
}
