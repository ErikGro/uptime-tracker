package store

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Store struct {
	db *gorm.DB
}

func Open(dbPath string) (*Store, error) {
	gormLogger := logger.New(log.New(os.Stderr, "[gorm] ", log.LstdFlags), logger.Config{
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: true,
		SlowThreshold:             200 * time.Millisecond,
	})
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, fmt.Errorf("open sqlite at %q: %w", dbPath, err)
	}
	if err := db.AutoMigrate(&URL{}, &Check{}, &Setting{}); err != nil {
		return nil, fmt.Errorf("automigrate: %w", err)
	}
	s := &Store{db: db}
	if err := s.seedSettings(); err != nil {
		return nil, fmt.Errorf("seed settings: %w", err)
	}
	return s, nil
}

func (s *Store) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func isNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
