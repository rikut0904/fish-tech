package place

import "time"

// PlaceCache は店舗キャッシュ情報です。
type PlaceCache struct {
	ID        string
	Name      string
	Address   *string
	Genre     *string
	FetchedAt time.Time
}
