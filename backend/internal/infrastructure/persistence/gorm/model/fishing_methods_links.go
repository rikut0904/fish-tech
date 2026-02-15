package model

import "time"

// FishingMethodsLinks は fishing_methods_links テーブルのGORMモデルです。
type FishingMethodsLinks struct {
FishID    string     "gorm:\"column:fish_id;type:uuid;primaryKey\""
	MethodID  string     "gorm:\"column:method_id;type:uuid;primaryKey\""
	IsPrimary bool       `gorm:"column:is_primary;not null;default:false"`
	Explain   string     `gorm:"column:explain"`
	CreatedAt time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

// TableName はテーブル名を返します。
func (FishingMethodsLinks) TableName() string { return "fishing_methods_links" }
