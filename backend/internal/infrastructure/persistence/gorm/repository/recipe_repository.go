package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	recipeDomain "fish-tech/internal/domain/recipe"
	"fish-tech/internal/infrastructure/persistence/gorm/model"
	"fish-tech/internal/shared/timeutil"
)

const recipeCacheValidDuration = 24 * time.Hour

// RecipeRepository はレシピ機能向けのリポジトリです。
type RecipeRepository struct {
	db *gorm.DB
}

type recipeSearchRow struct {
	ID          string
	Title       string
	ImageURL    *string
	RecipeURL   string
	CookingTime *string
	Cost        *string
	Score       *int
	Explain     string
	IsLikes     bool
}

// NewRecipeRepository はレシピ機能向けリポジトリを作成します。
func NewRecipeRepository(db *gorm.DB) (*RecipeRepository, error) {
	if err := db.AutoMigrate(&model.RecipeCache{}, &model.FishRecipeLinks{}, &model.UserRecipeLinks{}); err != nil {
		return nil, err
	}

	return &RecipeRepository{db: db}, nil
}

// ListSeasonalFishes は指定月の旬魚一覧を返します。
func (r *RecipeRepository) ListSeasonalFishes(ctx context.Context, month int) ([]recipeDomain.SeasonalFish, error) {
	var rows []struct {
		ID       string
		Name     string
		ImageURL string
	}

	if err := r.db.WithContext(ctx).
		Table("fish AS f").
		Select("f.id AS id, COALESCE(NULLIF(TRIM(f.name_ja), ''), f.name) AS name, f.image_url AS image_url").
		Joins("JOIN season AS s ON s.fish_id = f.id").
		Where("s.month = ?", month).
		Order("f.created_at DESC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		if err := r.db.WithContext(ctx).
			Table("fish AS f").
			Select("f.id AS id, COALESCE(NULLIF(TRIM(f.name_ja), ''), f.name) AS name, f.image_url AS image_url").
			Order("f.created_at DESC").
			Limit(12).
			Scan(&rows).Error; err != nil {
			return nil, err
		}
	}

	result := make([]recipeDomain.SeasonalFish, 0, len(rows))
	for _, row := range rows {
		result = append(result, recipeDomain.SeasonalFish{
			ID:       row.ID,
			Name:     row.Name,
			ImageURL: strings.TrimSpace(row.ImageURL),
		})
	}

	return result, nil
}

// FindFishByID は魚を取得します。
func (r *RecipeRepository) FindFishByID(ctx context.Context, fishID string) (*recipeDomain.SeasonalFish, error) {
	var row struct {
		ID       string
		Name     string
		ImageURL string
	}

	err := r.db.WithContext(ctx).
		Table("fish AS f").
		Select("f.id AS id, COALESCE(NULLIF(TRIM(f.name_ja), ''), f.name) AS name, f.image_url AS image_url").
		Where("f.id = ?", strings.TrimSpace(fishID)).
		Take(&row).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &recipeDomain.SeasonalFish{
		ID:       row.ID,
		Name:     row.Name,
		ImageURL: strings.TrimSpace(row.ImageURL),
	}, nil
}

// FindFishByName は魚名から魚を取得します。
func (r *RecipeRepository) FindFishByName(ctx context.Context, fishName string) (*recipeDomain.SeasonalFish, error) {
	return r.findFishByName(ctx, fishName, 0, false)
}

// FindSeasonalFishByName は指定月の旬魚から魚を取得します。
func (r *RecipeRepository) FindSeasonalFishByName(ctx context.Context, fishName string, month int) (*recipeDomain.SeasonalFish, error) {
	return r.findFishByName(ctx, fishName, month, true)
}

// ListFreshRecipesByFishID は指定魚の新しいレシピキャッシュを返します。
func (r *RecipeRepository) ListFreshRecipesByFishID(ctx context.Context, fishID string, userID string, favorite *bool, count int, page int) ([]recipeDomain.RecipeRecommendation, int, error) {
	condition := recipeDomain.SearchCondition{
		UserID:   strings.TrimSpace(userID),
		Favorite: favorite,
		Count:    count,
		Page:     page,
	}
	return r.searchRecipes(ctx, condition, strings.TrimSpace(fishID), 0, false)
}

// SearchRecipes は魚名やキーワードからレシピを検索します。
func (r *RecipeRepository) SearchRecipes(ctx context.Context, condition recipeDomain.SearchCondition) ([]recipeDomain.RecipeRecommendation, int, error) {
	return r.searchRecipes(ctx, condition, "", 0, false)
}

// ReplaceFishRecipes は魚に紐づくレシピキャッシュを更新します。
func (r *RecipeRepository) ReplaceFishRecipes(ctx context.Context, fishID string, recipes []recipeDomain.RecipeRecommendation) error {
	now := timeutil.NowJST()

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, recipe := range recipes {
			cache := model.RecipeCache{
				ID:          recipe.ID,
				Title:       recipe.Title,
				ImageURL:    recipe.ImageURL,
				RecipeURL:   recipe.RecipeURL,
				CookingTime: recipe.CookingTime,
				Cost:        recipe.Cost,
				FetchedAt:   now,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"title", "image_url", "recipe_url", "cooking_time", "cost", "fetched_at"}),
			}).Create(&cache).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("fish_id = ?", strings.TrimSpace(fishID)).Delete(&model.FishRecipeLinks{}).Error; err != nil {
			return err
		}

		for _, recipe := range recipes {
			link := model.FishRecipeLinks{
				FishID:    strings.TrimSpace(fishID),
				RecipeID:  recipe.ID,
				Score:     recipe.Score,
				Explain:   recipe.Explain,
				CreatedAt: now,
				UpdatedAt: &now,
			}
			if err := tx.Create(&link).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *RecipeRepository) findFishByName(ctx context.Context, fishName string, month int, onlySeasonal bool) (*recipeDomain.SeasonalFish, error) {
	trimmedFishName := strings.TrimSpace(fishName)
	if trimmedFishName == "" {
		return nil, nil
	}

	pattern := "%" + trimmedFishName + "%"
	query := r.db.WithContext(ctx).
		Table("fish AS f").
		Select("f.id AS id, COALESCE(NULLIF(TRIM(f.name_ja), ''), f.name) AS name, f.image_url AS image_url")
	if onlySeasonal {
		query = query.Joins("JOIN season AS s ON s.fish_id = f.id").Where("s.month = ?", month)
	}

	var row struct {
		ID       string
		Name     string
		ImageURL string
	}
	err := query.
		Where("f.name_ja = ? OR f.name = ? OR f.name_ja ILIKE ? OR f.name ILIKE ?", trimmedFishName, trimmedFishName, pattern, pattern).
		Order(fmt.Sprintf("CASE WHEN COALESCE(NULLIF(TRIM(f.name_ja), ''), f.name) = %s THEN 0 WHEN f.name = %s THEN 0 ELSE 1 END, f.created_at DESC", quoteLiteral(trimmedFishName), quoteLiteral(trimmedFishName))).
		Take(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &recipeDomain.SeasonalFish{
		ID:       row.ID,
		Name:     row.Name,
		ImageURL: strings.TrimSpace(row.ImageURL),
	}, nil
}

func (r *RecipeRepository) searchRecipes(ctx context.Context, condition recipeDomain.SearchCondition, fishID string, month int, onlySeasonal bool) ([]recipeDomain.RecipeRecommendation, int, error) {
	freshFrom := timeutil.NowJST().Add(-recipeCacheValidDuration)
	base := r.db.WithContext(ctx).
		Table("fish_recipe_links AS frl").
		Joins("JOIN recipe_cache AS rc ON rc.id = frl.recipe_id").
		Joins("JOIN fish AS f ON f.id = frl.fish_id").
		Where("rc.fetched_at >= ?", freshFrom)

	if strings.TrimSpace(condition.UserID) != "" {
		base = base.Joins("LEFT JOIN user_recipe_links AS url ON url.recipe_id = rc.id AND url.user_id = ?", strings.TrimSpace(condition.UserID))
	} else {
		base = base.Joins("LEFT JOIN user_recipe_links AS url ON 1 = 0")
	}
	if onlySeasonal {
		base = base.Joins("JOIN season AS s ON s.fish_id = f.id").Where("s.month = ?", month)
	}
	if strings.TrimSpace(fishID) != "" {
		base = base.Where("frl.fish_id = ?", strings.TrimSpace(fishID))
	}
	if strings.TrimSpace(condition.FishName) != "" {
		pattern := "%" + strings.TrimSpace(condition.FishName) + "%"
		base = base.Where("(COALESCE(NULLIF(TRIM(f.name_ja), ''), f.name) ILIKE ? OR f.name ILIKE ?)", pattern, pattern)
	}
	if strings.TrimSpace(condition.Keyword) != "" {
		pattern := "%" + strings.TrimSpace(condition.Keyword) + "%"
		base = base.Where("(rc.title ILIKE ? OR frl.explain ILIKE ?)", pattern, pattern)
	}
	if condition.Favorite != nil {
		base = base.Where("COALESCE(url.is_likes, FALSE) = ?", *condition.Favorite)
	}

	var total int64
	if err := base.Session(&gorm.Session{}).Distinct("rc.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	rows := make([]recipeSearchRow, 0, condition.Count)
	offset := (condition.Page - 1) * condition.Count
	err := base.
		Select(`
			rc.id,
			rc.title,
			rc.image_url,
			rc.recipe_url,
			rc.cooking_time,
			rc.cost,
			MAX(frl.score) AS score,
			MIN(frl.explain) AS explain,
			COALESCE(BOOL_OR(url.is_likes), FALSE) AS is_likes
		`).
		Group("rc.id, rc.title, rc.image_url, rc.recipe_url, rc.cooking_time, rc.cost").
		Order("COALESCE(MAX(frl.score), 0) DESC, MAX(rc.fetched_at) DESC, rc.id ASC").
		Limit(condition.Count).
		Offset(offset).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	return toRecipeRecommendations(rows), int(total), nil
}

func toRecipeRecommendations(rows []recipeSearchRow) []recipeDomain.RecipeRecommendation {
	result := make([]recipeDomain.RecipeRecommendation, 0, len(rows))
	for _, row := range rows {
		result = append(result, recipeDomain.RecipeRecommendation{
			ID:          row.ID,
			Title:       row.Title,
			ImageURL:    row.ImageURL,
			RecipeURL:   row.RecipeURL,
			CookingTime: row.CookingTime,
			Cost:        row.Cost,
			Score:       row.Score,
			IsLikes:     row.IsLikes,
			Explain:     row.Explain,
		})
	}

	return result
}

func quoteLiteral(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

// ExistsUser はユーザーの存在確認を行います。
func (r *RecipeRepository) ExistsUser(ctx context.Context, userID string) (bool, error) {
	var row model.User
	err := r.db.WithContext(ctx).Select("user_id").Where("user_id = ?", strings.TrimSpace(userID)).Take(&row).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}

// ExistsRecipe はレシピの存在確認を行います。
func (r *RecipeRepository) ExistsRecipe(ctx context.Context, recipeID string) (bool, error) {
	var row model.RecipeCache
	err := r.db.WithContext(ctx).Select("id").Where("id = ?", strings.TrimSpace(recipeID)).Take(&row).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}

// UpdateRecipeFavorite はユーザーとレシピの関連を必要時に作成しつつ更新します。
func (r *RecipeRepository) UpdateRecipeFavorite(ctx context.Context, userID string, recipeID string, isLikes bool) error {
	now := timeutil.NowJST()
	link := model.UserRecipeLinks{
		UserID:    strings.TrimSpace(userID),
		RecipeID:  strings.TrimSpace(recipeID),
		IsLikes:   isLikes,
		CreatedAt: now,
		UpdatedAt: &now,
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "recipe_id"}},
		DoUpdates: clause.Assignments(map[string]any{"is_likes": isLikes, "updated_at": now}),
	}).Create(&link).Error
}
