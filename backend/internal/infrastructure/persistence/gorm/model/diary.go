package model

import "time"

// Diary は diary テーブルのGORMモデルです。
type Diary struct {
	ID        string     `gorm:"column:id;type:uuid;primaryKey"`
	UserID    string     `gorm:"column:user_id;type:uuid;not null;index"`
	FishID    string     `gorm:"column:fish_id;type:uuid;not null;index"`
	RecipeID  *string    `gorm:"column:recipe_id;type:varchar(255)"`
	PlaceID   *string    `gorm:"column:place_id;type:uuid"`
	Date      time.Time  `gorm:"column:date;not null"`
	Explain   *string    `gorm:"column:explain"`
	Score     *int       `gorm:"column:score"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

// TableName はテーブル名を返します。
func (Diary) TableName() string { return "diary" }
