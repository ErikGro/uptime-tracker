package store

import (
	"database/sql"
	"time"
)

type Status string

const (
	StatusUnknown Status = "UNKNOWN"
	StatusUp      Status = "UP"
	StatusDown    Status = "DOWN"
)

type URL struct {
	ID                  uint   `gorm:"primaryKey"`
	Label               string `gorm:"not null"`
	URL                 string `gorm:"not null;uniqueIndex"`
	CurrentStatus       Status `gorm:"type:text;default:UNKNOWN;not null"`
	ConsecutiveFailures int    `gorm:"not null;default:0"`
	LastCheckedAt       sql.NullTime
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type Check struct {
	ID         uint      `gorm:"primaryKey"`
	URLID      uint      `gorm:"not null;index"`
	CheckedAt  time.Time `gorm:"not null;index"`
	StatusCode int
	LatencyMs  int64
	OK         bool
	Error      string
}

type Setting struct {
	Key   string `gorm:"primaryKey"`
	Value string `gorm:"not null"`
}
