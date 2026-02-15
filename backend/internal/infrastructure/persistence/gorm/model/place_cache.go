package model

import "time"

// PlaceCache は place_cache テーブルのGORMモデルです。
type PlaceCache struct {
ID        string    "gorm:\"column:id;type:uuid;primaryKey\""
	Name      string    `gorm:"column:name;not null"`
	Address   *string   `gorm:"column:address"`
	Genre     *string   `gorm:"column:genre"`
	FetchedAt time.Time `gorm:"column:fetched_at;not null"`
}

// TableName はテーブル名を返します。
func (PlaceCache) TableName() string { return "place_cache" }
