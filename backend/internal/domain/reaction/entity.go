package reaction

import "time"

// FishReaction は魚に対するユーザーのリアクションを表します。
type FishReaction struct {
	UserID    string
	FishID    string
	Like      bool
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// RecipeReaction は user_recipe_links の情報です。
type RecipeReaction struct {
	UserID    string
	RecipeID  string
	Like      bool
	CreatedAt time.Time
	UpdatedAt *time.Time
}
