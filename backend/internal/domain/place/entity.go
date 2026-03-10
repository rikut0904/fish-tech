package place

import "time"

// PlaceCache は店舗キャッシュ情報です。
type PlaceCache struct {
	ID            string
	Name          string
	Address       *string
	Lat           *float64
	Lng           *float64
	Coupon        *string
	Genre         *string
	Card          *string
	Logo          *string
	LargeAreaCode *string
	FetchedAt     time.Time
}

// SearchCondition はおすすめ店舗検索条件です。
type SearchCondition struct {
	FishName      string
	Keyword       string
	CityCode      string
	UserID        string
	Favorite      bool
	Count         int
	Page          int
	LargeArea     string
	SmallAreaCode string
}

// RecommendedPlace はおすすめ店舗情報です。
type RecommendedPlace struct {
	ID            string `json:"-"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	Lat           string `json:"lat"`
	Lng           string `json:"lng"`
	Coupon        string `json:"coupon"`
	Genre         string `json:"genre"`
	Card          string `json:"card"`
	Logo          string `json:"logo"`
	SmallAreaCode string `json:"-"`
	LargeAreaCode string `json:"-"`
}
