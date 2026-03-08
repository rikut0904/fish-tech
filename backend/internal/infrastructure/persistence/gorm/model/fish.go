package model

import "time"

// Fish は fish テーブルのGORMモデルです。
type Fish struct {
	ID        string     `gorm:"column:id;type:uuid;primaryKey"`
	NameJa    string     `gorm:"column:name_ja;not null"`
	Name      string     `gorm:"column:name;not null"`
	Category  string     `gorm:"column:category"`
	Explain   string     `gorm:"column:explain"`
	ImageURL  string     `gorm:"column:image_url"`
	LinkURL   string     `gorm:"column:link_url"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

// TableName はテーブル名を返します。
func (Fish) TableName() string { return "fish" }
