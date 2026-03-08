package model

import "time"

// AdminFish は管理画面の魚情報を保存するGORMモデルです。
type AdminFish struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey"`
	Name        string    `gorm:"column:name;not null"`
	Category    string    `gorm:"column:category"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
}

// TableName はテーブル名を返します。
func (AdminFish) TableName() string { return "admin_fishes" }
