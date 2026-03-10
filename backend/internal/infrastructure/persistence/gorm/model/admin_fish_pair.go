package model

import "time"

// AdminFishPair は管理画面の魚相性情報を保存するGORMモデルです。
type AdminFishPair struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey"`
	FishIDa   string    `gorm:"column:fish_id_a;type:uuid;not null;index"`
	FishIDb   string    `gorm:"column:fish_id_b;type:uuid;not null;index"`
	Score     int       `gorm:"column:score;not null"`
	Memo      string    `gorm:"column:memo"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

// TableName はテーブル名を返します。
func (AdminFishPair) TableName() string { return "admin_fish_pairs" }
