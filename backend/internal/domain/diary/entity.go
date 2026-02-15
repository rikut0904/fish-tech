package diary

import "time"

// Diary は日記情報です。
type Diary struct {
	ID        string
	UserID    string
	FishID    string
	RecipeID  *string
	PlaceID   *string
	Date      time.Time
	Explain   *string
	Score     *int
	CreatedAt time.Time
	UpdatedAt *time.Time
}
