package database

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"easyblog/config"
)

// connect to database by config
func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := openDB(cfg)
	if err != nil {
		return nil, err
	}
	err = configureDB(db, &cfg.PgConnCfg)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// close database
func CloseDB(db *gorm.DB) error {
	if db == nil {
		return errors.New("db pointer is nil")
	}
	sqlDb, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to fetch sql.DB: %w", err)
	}
	err = sqlDb.Close()
	if err != nil {
		return fmt.Errorf("failed to close sql.DB: %w", err)
	}
	return nil
}

// open database connection
func openDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		cfg.PgCfg.Host,
		cfg.PgCfg.Port,
		cfg.PgCfg.User,
		cfg.PgCfg.Password,
		cfg.PgCfg.Database,
		"UTC",
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres: %w", err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sql.DB: %w", err)
	}

	for i := 1; i <= cfg.PgConnCfg.ConnMaxRetries; i++ {
		if err = sqlDb.Ping(); err != nil {
			log.Printf("Ping failed, retry %d/%d: %v", i, cfg.PgConnCfg.ConnMaxRetries, err)
			time.Sleep(cfg.PgConnCfg.ConnRetryInterval)
		} else {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect postgres after %d retries: %w", cfg.PgConnCfg.ConnMaxRetries, err)
	}
	return db, nil
}

// setup database connection pool
func configureDB(db *gorm.DB, pgConnCfg *config.PgConnConfig) error {
	if db == nil {
		return errors.New("db pointer is nil")
	}
	sqlDb, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to fetch sql.DB: %w", err)
	}
	sqlDb.SetMaxOpenConns(pgConnCfg.MaxOpenConns)
	sqlDb.SetMaxIdleConns(pgConnCfg.MaxIdleConns)
	sqlDb.SetConnMaxLifetime(pgConnCfg.ConnMaxLifetime)
	sqlDb.SetConnMaxIdleTime(pgConnCfg.ConnMaxIdletime)
	return nil
}
