package model

import "time"

// UserFishLinks は user_fish_links テーブルのGORMモデルです。
type UserFishLinks struct {
	UserID    string     `gorm:"column:user_id;primaryKey"`
	FishID    string     `gorm:"column:fish_id;primaryKey"`
	Like      bool       `gorm:"column:like;not null;default:false"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

// TableName はテーブル名を返します。
func (UserFishLinks) TableName() string { return "user_fish_links" }
