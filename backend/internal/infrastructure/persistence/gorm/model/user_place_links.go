package model

import "time"

// UserPlaceLinks は user_place_links テーブルのGORMモデルです。
type UserPlaceLinks struct {
	UserID    string     `gorm:"column:user_id;type:uuid;primaryKey"`
	PlaceID   string     `gorm:"column:place_id;type:text;primaryKey"`
	IsLikes   bool       `gorm:"column:is_likes;not null;default:false"`
	Explain   *string    `gorm:"column:explain"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

// TableName はテーブル名を返します。
func (UserPlaceLinks) TableName() string { return "user_place_links" }
