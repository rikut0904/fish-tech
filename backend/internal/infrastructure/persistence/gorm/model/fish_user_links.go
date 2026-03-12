package model

import "time"

// FishUserLinks は fish_user_links テーブルのGORMモデルです。
type FishUserLinks struct {
	FishID    string     `gorm:"column:fish_id;type:uuid;primaryKey"`
	UserID    string     `gorm:"column:user_id;type:uuid;primaryKey"`
	IsLikes   bool       `gorm:"column:is_likes;not null;default:false"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

// TableName はテーブル名を返します。
func (FishUserLinks) TableName() string { return "fish_user_links" }
