package model

import "time"

// Fish は fish テーブルのGORMモデルです。
type Fish struct {
	ID        string     `gorm:"column:id;primaryKey"`
	NameJa    string     `gorm:"column:name_ja;not null"`
	Name      string     `gorm:"column:name;not null"`
	Explain   string     `gorm:"column:explain"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

// TableName はテーブル名を返します。
func (Fish) TableName() string { return "fish" }
