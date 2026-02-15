package model

import "time"

// FishingMethods は fishing_methods テーブルのGORMモデルです。
type FishingMethods struct {
	ID        string     `gorm:"column:id;primaryKey"`
	Name      string     `gorm:"column:name;uniqueIndex;not null"`
	Explain   string     `gorm:"column:explain"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

// TableName はテーブル名を返します。
func (FishingMethods) TableName() string { return "fishing_methods" }
