package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/lib/pq"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	if DB != nil {
		log.Println("[postgres] using existing connection")
		return DB
	}

	appEnv := os.Getenv("APP_ENV") // "production", "dev", "test"

	if appEnv == "" {
		appEnv = "test"
		log.Println("[postgres] APP_ENV not set, defaulting to 'test'")
	}

	switch appEnv {
	case "test":
		var err error
		DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			log.Fatalf("[postgres] failed to connect SQLite for tests: %v", err)
		}
		log.Println("[postgres] test database (SQLite) connected")
		return DB
	case "dev", "prod":
		host := os.Getenv("POSTGRESQL_HOST")
		port := os.Getenv("POSTGRESQL_PORT")
		user := os.Getenv("POSTGRESQL_USERNAME")
		password := os.Getenv("POSTGRESQL_PASSWORD")
		dbname := os.Getenv("POSTGRESQL_DB")
		sslmode := os.Getenv("POSTGRESQL_SSLMODE")

		if dbname == "" {
			dbname = "gamestats"
			log.Printf("[postgres] POSTGRESQL_DB not set, defaulting to 'gamestats'")
		}

		if sslmode == "" {
			sslmode = "disable"
			log.Printf("[postgres] POSTGRESQL_SSLMODE not set, defaulting to 'disable'")
		}

		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode,
		)

		log.Printf("[postgres] connecting to Postgres at %s:%s/%s", host, port, dbname)

		var sqlDB *sql.DB
		var err error

		for i := 1; i <= 10; i++ {
			sqlDB, err = sql.Open("postgres", dsn)
			if err != nil {
				log.Printf("[postgres] attempt %d: failed to open database: %v", i, err)
				time.Sleep(3 * time.Second)
				continue
			}
			if err = sqlDB.Ping(); err != nil {
				log.Printf("[postgres] attempt %d: database not responding: %v", i, err)
				time.Sleep(3 * time.Second)
				continue
			}
			log.Printf("[postgres] connection successful on attempt %d", i)
			break
		}

		if err != nil {
			log.Fatalf("[postgres] failed to connect to database after multiple attempts: %v", err)
		}

		DB, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatalf("[postgres] failed to create GORM DB: %v", err)
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		log.Printf("[postgres] database '%s' connected with GORM", appEnv)
		return DB
	default:
		log.Fatalf("[postgres] invalid APP_ENV: %s", appEnv)
	}

	return nil
}
