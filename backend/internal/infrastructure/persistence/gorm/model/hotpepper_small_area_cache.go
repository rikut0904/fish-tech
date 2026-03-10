package model

import "time"

// HotpepperSmallAreaCache は hotpepper_small_area_cache テーブルのGORMモデルです。
type HotpepperSmallAreaCache struct {
	LargeAreaCode string    `gorm:"column:large_area_code;type:text;primaryKey"`
	Code          string    `gorm:"column:code;type:text;primaryKey"`
	Name          *string   `gorm:"column:name"`
	MiddleArea    *string   `gorm:"column:middle_area_code"`
	FetchedAt     time.Time `gorm:"column:fetched_at;not null"`
}

// TableName はテーブル名を返します。
func (HotpepperSmallAreaCache) TableName() string { return "hotpepper_small_area_cache" }
