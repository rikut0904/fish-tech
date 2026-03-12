package model

import "time"

// FishPlaceLinks は fish_place_links テーブルのGORMモデルです。
type FishPlaceLinks struct {
	FishID    string     `gorm:"column:fish_id;type:uuid;primaryKey"`
	PlaceID   string     `gorm:"column:place_id;type:text;primaryKey"`
	Score     *int       `gorm:"column:score"`
	Explain   string     `gorm:"column:explain"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

// TableName はテーブル名を返します。
func (FishPlaceLinks) TableName() string { return "fish_place_links" }
