package model

import "time"

// RecipeCache は recipe_cache テーブルのGORMモデルです。
type RecipeCache struct {
	ID          string    `gorm:"column:id;type:text;primaryKey"`
	Title       string    `gorm:"column:title;not null"`
	ImageURL    *string   `gorm:"column:image_url"`
	RecipeURL   string    `gorm:"column:recipe_url;not null"`
	CookingTime *string   `gorm:"column:cooking_time"`
	Cost        *string   `gorm:"column:cost"`
	FetchedAt   time.Time `gorm:"column:fetched_at;not null"`
}

// TableName はテーブル名を返します。
func (RecipeCache) TableName() string { return "recipe_cache" }
