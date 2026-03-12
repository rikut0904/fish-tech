package model

import "time"

// PlaceCache は place_cache テーブルのGORMモデルです。
type PlaceCache struct {
	ID            string    `gorm:"column:id;type:text;primaryKey"`
	Name          string    `gorm:"column:name;not null"`
	Address       *string   `gorm:"column:address"`
	Lat           *float64  `gorm:"column:lat"`
	Lng           *float64  `gorm:"column:lng"`
	Coupon        *string   `gorm:"column:coupon"`
	Genre         *string   `gorm:"column:genre"`
	Card          *string   `gorm:"column:card"`
	Logo          *string   `gorm:"column:logo"`
	LargeAreaCode *string   `gorm:"column:large_area_code;index"`
	SmallAreaCode *string   `gorm:"column:small_area_code;index"`
	FetchedAt     time.Time `gorm:"column:fetched_at;not null"`
}

// TableName はテーブル名を返します。
func (PlaceCache) TableName() string { return "place_cache" }
