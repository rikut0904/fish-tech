package model

import "time"

// FishStory は fish_story テーブルのGORMモデルです。
type FishStory struct {
	ID        string     "gorm:\"column:id;type:uuid;primaryKey\""
	FishID    string     "gorm:\"column:fish_id;type:uuid;not null;index\""
	Story     []byte     `gorm:"column:story;type:jsonb;not null"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

// TableName はテーブル名を返します。
func (FishStory) TableName() string { return "fish_story" }
