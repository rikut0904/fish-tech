package model

import "time"

// UserRecipeLinks は user_recipe_links テーブルのGORMモデルです。
type UserRecipeLinks struct {
	UserID    string     `gorm:"column:user_id;primaryKey"`
	RecipeID  string     `gorm:"column:recipe_id;primaryKey"`
	Like      bool       `gorm:"column:like;not null;default:false"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

// TableName はテーブル名を返します。
func (UserRecipeLinks) TableName() string { return "user_recipe_links" }
