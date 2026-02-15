package model

import "time"

// FishRecipeLinks は fish_recipe_links テーブルのGORMモデルです。
type FishRecipeLinks struct {
	FishID    string     `gorm:"column:fish_id;primaryKey"`
	RecipeID  string     `gorm:"column:recipe_id;primaryKey"`
	Score     *int       `gorm:"column:score"`
	Explain   string     `gorm:"column:explain"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

// TableName はテーブル名を返します。
func (FishRecipeLinks) TableName() string { return "fish_recipe_links" }
