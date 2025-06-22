package models

import (
	"time"
)

type DbPlayerKillAgg struct {
	PlayerId    string    `gorm:"primaryKey;index"`
	UnixHour    time.Time `gorm:"primaryKey"`
	KillCount   int
	FFKillCount int
	DeathCount  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
