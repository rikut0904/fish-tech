package hello

import (
	"context"
	"time"
	"your_project/internal/domain/fish"
)

type RecipeUsecase struct {
	repo         FishRepository        // DB操作用
	recipeClient *api.RakutenRecipeClient // API操作用
}

func (u *RecipeUsecase) GetRecipesByFish(ctx context.Context, fishID string) ([]fish.Recipe, error) {
	// 1. まずはDBのキャッシュ(recipe_cache)を確認
	// ER図の fish_recipe_links と recipe_cache をJOINして取得
	cachedRecipes, err := u.repo.GetCachedRecipes(fishID)
	
	// キャッシュが有効（例えば1日以内）であればそれを返す
	if err == nil && len(cachedRecipes) > 0 {
		return cachedRecipes, nil
	}

	// 2. キャッシュがない場合、fishテーブルから rakuten_category_id を取得
	f, err := u.repo.GetFishByID(fishID)
	if err != nil || f.RakutenCategoryID == "" {
		return nil, err
	}

	// 3. 楽天APIから最新レシピを取得
	newRecipes, err := u.recipeClient.FetchRankingByCategoryId(f.RakutenCategoryID)
	if err != nil {
		return nil, err
	}

	// 4. 取得したデータを recipe_cache に保存（キャッシュ更新）
	// 同時に fish_recipe_links で魚と紐付け
	go u.repo.SaveRecipesToCache(fishID, newRecipes)

	return newRecipes, nil
}