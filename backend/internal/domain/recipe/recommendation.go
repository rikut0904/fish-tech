package recipe

// SeasonalFish は旬の魚情報です。
type SeasonalFish struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}

// RecipeRecommendation はレシピ提案情報です。
type RecipeRecommendation struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	ImageURL    *string `json:"imageUrl,omitempty"`
	RecipeURL   string  `json:"recipeUrl"`
	CookingTime *string `json:"cookingTime,omitempty"`
	Cost        *string `json:"cost,omitempty"`
	Score       *int    `json:"score,omitempty"`
	IsLikes     bool    `json:"isLikes"`
	Explain     string  `json:"explain"`
}
