package model

import "time"

// Season は season テーブルのGORMモデルです。
type Season struct {
	FishID    string     "gorm:\"column:fish_id;type:uuid;primaryKey\""
	Month     int        `gorm:"column:month;primaryKey"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

// TableName はテーブル名を返します。
func (Season) TableName() string { return "season" }
