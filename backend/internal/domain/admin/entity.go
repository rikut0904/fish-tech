package admin

import "time"

// Fish は管理画面で扱う魚エンティティです。
type Fish struct {
	ID          string
	Name        string
	Category    string
	Description string
	ImageURL    string
	LinkURL     string
	CreatedAt   time.Time
}

// FishPair は管理画面で扱う魚同士の相性エンティティです。
type FishPair struct {
	ID        string
	FishIDa   string
	FishIDb   string
	Score     int
	Memo      string
	CreatedAt time.Time
}
