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
		log.Println("[postgres] usando conexão existente")
		return DB
	}

	appEnv := os.Getenv("APP_ENV") // "production", "dev", "test"

	if appEnv == "" {
		appEnv = "test"
		log.Println("[postgres] APP_ENV não definido, usando 'test'")
	}

	switch appEnv {
	case "test":
		// SQLite in-memory para testes
		var err error
		DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			log.Fatalf("[postgres] falha ao conectar SQLite para testes: %v", err)
		}
		log.Println("[postgres] banco de dados de teste (SQLite) conectado")
		return DB
	case "dev", "prod":
		// Postgres
		host := os.Getenv("POSTGRESQL_HOST")
		port := os.Getenv("POSTGRESQL_PORT")
		user := os.Getenv("POSTGRESQL_USERNAME")
		password := os.Getenv("POSTGRESQL_PASSWORD")
		dbname := os.Getenv("POSTGRESQL_DB")
		sslmode := os.Getenv("POSTGRESQL_SSLMODE")

		if dbname == "" {
			dbname = "gamestats"
			log.Printf("[postgres] POSTGRESQL_DB não definido, usando 'gamestats'")
		}

		if sslmode == "" {
			sslmode = "disable"
			log.Printf("[postgres] POSTGRESQL_SSLMODE não definido, usando 'disable'")
		}

		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode,
		)

		log.Printf("[postgres] tentando conectar ao Postgres em %s:%s/%s (%s:%s)", host, port, dbname, user, password)

		var sqlDB *sql.DB
		var err error

		for i := 1; i <= 10; i++ {
			sqlDB, err = sql.Open("postgres", dsn)
			if err != nil {
				log.Printf("[postgres] tentativa %d: falha ao abrir banco de dados: %v", i, err)
				time.Sleep(3 * time.Second)
				continue
			}
			if err = sqlDB.Ping(); err != nil {
				log.Printf("[postgres] tentativa %d: banco ainda não respondeu: %v", i, err)
				time.Sleep(3 * time.Second)
				continue
			}
			log.Printf("[postgres] conexão bem-sucedida na tentativa %d", i)
			break
		}

		if err != nil {
			log.Fatalf("[postgres] não consegui conectar ao banco depois de várias tentativas: %v", err)
		}

		DB, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatalf("[postgres] falha ao criar GORM DB: %v", err)
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		log.Printf("[postgres] banco de dados '%s' conectado com GORM", appEnv)
		return DB
	default:
		log.Fatalf("[postgres] APP_ENV inválido: %s", appEnv)
	}

	return nil
}
