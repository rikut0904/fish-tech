package model

import "time"

// FishPair は fish_pair テーブルのGORMモデルです。
type FishPair struct {
	FishAID   string     `gorm:"column:fish_a_id;type:uuid;primaryKey;uniqueIndex:uq_fish_pair_ids"`
	FishBID   string     `gorm:"column:fish_b_id;type:uuid;primaryKey;uniqueIndex:uq_fish_pair_ids"`
	Result    string     `gorm:"column:result;not null"`
	Explain   string     `gorm:"column:explain"`
	Score     int        `gorm:"column:score;not null"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

// TableName はテーブル名を返します。
func (FishPair) TableName() string { return "fish_pair" }
