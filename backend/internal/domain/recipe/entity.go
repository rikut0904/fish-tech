package recipe

import "time"

// RecipeCache はレシピキャッシュ情報です。
type RecipeCache struct {
	ID          string
	Title       string
	ImageURL    *string
	RecipeURL   string
	CookingTime *string
	Cost        *string
	FetchedAt   time.Time
}

// SearchCondition はレシピ検索条件です。
type SearchCondition struct {
	FishName string
	Keyword  string
	UserID   string
	Favorite *bool
	Count    int
	Page     int
}
