package recipe

import (
	"context"
	"errors"
	"log"
	"strings"
	"unicode"

	recipeDomain "fish-tech/internal/domain/recipe"
	"fish-tech/internal/shared/timeutil"
)

var (
	// ErrInvalidUserID はユーザーIDが不正な場合のエラーです。
	ErrInvalidUserID = errors.New("userIdは必須です")
	// ErrInvalidRecipeID はレシピIDが不正な場合のエラーです。
	ErrInvalidRecipeID = errors.New("recipeIdは必須です")
	// ErrUserNotFound はユーザーが存在しない場合のエラーです。
	ErrUserNotFound = errors.New("ユーザーが見つかりません")
	// ErrRecipeNotFound はレシピが存在しない場合のエラーです。
	ErrRecipeNotFound = errors.New("レシピが見つかりません")
	// ErrFavoriteUserIDRequired はfavorite指定時にuserIdが不足している場合のエラーです。
	ErrFavoriteUserIDRequired = errors.New("favorite指定時はuserIdが必要です")
	// ErrFishNotInSeason は指定した魚が今月の旬ではない場合のエラーです。
	ErrFishNotInSeason = errors.New("指定した魚は今月の旬ではありません")
	// ErrRakutenRecipesUnavailable は楽天レシピAPIからの取得に失敗した場合のエラーです。
	ErrRakutenRecipesUnavailable = errors.New("楽天レシピAPIからレシピを取得できませんでした")
)

const (
	defaultCount = 10
	defaultPage  = 1
	maxCount     = 100
)

// Repository はレシピ機能向けの永続化インターフェースです。
type Repository interface {
	ListSeasonalFishes(ctx context.Context, month int) ([]recipeDomain.SeasonalFish, error)
	FindFishByID(ctx context.Context, fishID string) (*recipeDomain.SeasonalFish, error)
	FindFishByName(ctx context.Context, fishName string) (*recipeDomain.SeasonalFish, error)
	FindSeasonalFishByName(ctx context.Context, fishName string, month int) (*recipeDomain.SeasonalFish, error)
	ListFreshRecipesByFishID(ctx context.Context, fishID string, userID string, favorite *bool, count int, page int) ([]recipeDomain.RecipeRecommendation, int, error)
	SearchRecipes(ctx context.Context, condition recipeDomain.SearchCondition) ([]recipeDomain.RecipeRecommendation, int, error)
	ReplaceFishRecipes(ctx context.Context, fishID string, recipes []recipeDomain.RecipeRecommendation) error
	ExistsUser(ctx context.Context, userID string) (bool, error)
	ExistsRecipe(ctx context.Context, recipeID string) (bool, error)
	UpdateRecipeFavorite(ctx context.Context, userID string, recipeID string, isLikes bool) error
}

// RakutenClient は楽天レシピAPIクライアントのインターフェースです。
type RakutenClient interface {
	Enabled() bool
	ListCategories(ctx context.Context) ([]RakutenCategory, error)
	GetCategoryRanking(ctx context.Context, categoryID string, categoryName string) ([]RakutenRecipe, error)
}

// RakutenCategory は楽天カテゴリ情報です。
type RakutenCategory struct {
	ID   string
	Name string
	Type string
}

// RakutenRecipe は楽天レシピ情報です。
type RakutenRecipe struct {
	ID              string
	Title           string
	ImageURL        *string
	RecipeURL       string
	CookingTime     *string
	Cost            *string
	Description     string
	Rank            int
	MatchedCategory string
}

// Response はレシピ画面向けレスポンスです。
type Response struct {
	Month          int                                 `json:"month"`
	SelectedFishID string                              `json:"selectedFishId"`
	Fishes         []recipeDomain.SeasonalFish         `json:"fishes"`
	Recipes        []recipeDomain.RecipeRecommendation `json:"recipes"`
	Page           int                                 `json:"page"`
	PerPage        int                                 `json:"perPage"`
	Count          int                                 `json:"count"`
	Total          int                                 `json:"total"`
}

// SearchResponse はレシピ検索レスポンスです。
type SearchResponse struct {
	Page    int                                 `json:"page"`
	PerPage int                                 `json:"perPage"`
	Count   int                                 `json:"count"`
	Total   int                                 `json:"total"`
	Items   []recipeDomain.RecipeRecommendation `json:"items"`
}

// UseCase はレシピ機能のユースケースです。
type UseCase interface {
	GetSeasonalRecipes(ctx context.Context, condition recipeDomain.SearchCondition) (Response, error)
	SearchRecipes(ctx context.Context, condition recipeDomain.SearchCondition) (SearchResponse, error)
	UpdateRecipeFavorite(ctx context.Context, userID string, recipeID string, isLikes bool) error
}

type recipeUseCase struct {
	repo    Repository
	rakuten RakutenClient
}

// NewRecipeUseCase はレシピユースケースを生成します。
func NewRecipeUseCase(repo Repository, rakuten RakutenClient) UseCase {
	return &recipeUseCase{repo: repo, rakuten: rakuten}
}

// GetSeasonalRecipes は旬の魚向けレシピ一覧を返します。
func (u *recipeUseCase) GetSeasonalRecipes(ctx context.Context, condition recipeDomain.SearchCondition) (Response, error) {
	month := int(timeutil.NowJST().Month())
	condition = normalizeCondition(condition)
	if condition.Favorite != nil && condition.UserID == "" {
		return Response{}, ErrFavoriteUserIDRequired
	}

	fishes, err := u.repo.ListSeasonalFishes(ctx, month)
	if err != nil {
		return Response{}, err
	}

	var selectedFish *recipeDomain.SeasonalFish
	if condition.FishName != "" {
		selectedFish, err = u.repo.FindSeasonalFishByName(ctx, condition.FishName, month)
		if err != nil {
			return Response{}, err
		}
		if selectedFish == nil {
			return Response{}, ErrFishNotInSeason
		}
	} else if len(fishes) > 0 {
		selectedFish = &fishes[0]
	}
	if selectedFish == nil {
		return Response{
			Month:   month,
			Fishes:  fishes,
			Recipes: []recipeDomain.RecipeRecommendation{},
			Page:    condition.Page,
			PerPage: condition.Count,
			Count:   0,
			Total:   0,
		}, nil
	}

	recipes, total, err := u.repo.ListFreshRecipesByFishID(ctx, selectedFish.ID, condition.UserID, condition.Favorite, condition.Count, condition.Page)
	if err != nil {
		return Response{}, err
	}
	if len(recipes) == 0 && u.rakuten != nil && u.rakuten.Enabled() {
		fetched, fetchErr := u.fetchRecipes(ctx, selectedFish.Name)
		if fetchErr != nil {
			log.Printf("recipe: 楽天APIから旬レシピ取得に失敗しました fish=%q err=%v", selectedFish.Name, fetchErr)
		} else if len(fetched) > 0 {
			if replaceErr := u.repo.ReplaceFishRecipes(ctx, selectedFish.ID, fetched); replaceErr != nil {
				return Response{}, replaceErr
			}
			recipes, total, err = u.repo.ListFreshRecipesByFishID(ctx, selectedFish.ID, condition.UserID, condition.Favorite, condition.Count, condition.Page)
			if err != nil {
				return Response{}, err
			}
		}
	}

	return Response{
		Month:          month,
		SelectedFishID: selectedFish.ID,
		Fishes:         fishes,
		Recipes:        recipes,
		Page:           condition.Page,
		PerPage:        condition.Count,
		Count:          len(recipes),
		Total:          total,
	}, nil
}

// SearchRecipes はレシピ一覧を検索します。
func (u *recipeUseCase) SearchRecipes(ctx context.Context, condition recipeDomain.SearchCondition) (SearchResponse, error) {
	condition = normalizeCondition(condition)
	if condition.Favorite != nil && condition.UserID == "" {
		return SearchResponse{}, ErrFavoriteUserIDRequired
	}

	items, total, err := u.repo.SearchRecipes(ctx, condition)
	if err != nil {
		return SearchResponse{}, err
	}

	if len(items) == 0 && condition.FishName != "" && u.rakuten != nil && u.rakuten.Enabled() {
		fish, findErr := u.repo.FindFishByName(ctx, condition.FishName)
		if findErr != nil {
			return SearchResponse{}, findErr
		}
		if fish != nil {
			fetched, fetchErr := u.fetchRecipes(ctx, fish.Name)
			if fetchErr != nil {
				log.Printf("recipe: 楽天APIから魚名指定レシピ取得に失敗しました fish=%q err=%v", fish.Name, fetchErr)
			} else if len(fetched) > 0 {
				if replaceErr := u.repo.ReplaceFishRecipes(ctx, fish.ID, fetched); replaceErr != nil {
					return SearchResponse{}, replaceErr
				}
				items, total, err = u.repo.SearchRecipes(ctx, condition)
				if err != nil {
					return SearchResponse{}, err
				}
			}
		}
	}
	if len(items) == 0 && condition.FishName == "" && u.rakuten != nil && u.rakuten.Enabled() {
		month := int(timeutil.NowJST().Month())
		fishes, fishErr := u.repo.ListSeasonalFishes(ctx, month)
		if fishErr != nil {
			return SearchResponse{}, fishErr
		}

		limit := len(fishes)
		if limit > 3 {
			limit = 3
		}
		var fetchedAny bool
		for i := 0; i < limit; i++ {
			fetched, fetchErr := u.fetchRecipes(ctx, fishes[i].Name)
			if fetchErr != nil {
				log.Printf("recipe: 楽天APIから汎用レシピ取得に失敗しました fish=%q err=%v", fishes[i].Name, fetchErr)
				continue
			}
			if len(fetched) == 0 {
				continue
			}
			if replaceErr := u.repo.ReplaceFishRecipes(ctx, fishes[i].ID, fetched); replaceErr != nil {
				continue
			}
			fetchedAny = true
		}
		if fetchedAny {
			items, total, err = u.repo.SearchRecipes(ctx, condition)
			if err != nil {
				return SearchResponse{}, err
			}
		}
	}

	return SearchResponse{
		Page:    condition.Page,
		PerPage: condition.Count,
		Count:   len(items),
		Total:   total,
		Items:   items,
	}, nil
}

// UpdateRecipeFavorite はユーザーのレシピお気に入り状態を更新します。
func (u *recipeUseCase) UpdateRecipeFavorite(ctx context.Context, userID string, recipeID string, isLikes bool) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return ErrInvalidUserID
	}

	recipeID = strings.TrimSpace(recipeID)
	if recipeID == "" {
		return ErrInvalidRecipeID
	}

	userExists, err := u.repo.ExistsUser(ctx, userID)
	if err != nil {
		return err
	}
	if !userExists {
		return ErrUserNotFound
	}

	recipeExists, err := u.repo.ExistsRecipe(ctx, recipeID)
	if err != nil {
		return err
	}
	if !recipeExists {
		return ErrRecipeNotFound
	}

	return u.repo.UpdateRecipeFavorite(ctx, userID, recipeID, isLikes)
}

func (u *recipeUseCase) fetchRecipes(ctx context.Context, fishName string) ([]recipeDomain.RecipeRecommendation, error) {
	categories, err := u.rakuten.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	category, ok := pickBestCategory(categories, fishName)
	if !ok {
		return []recipeDomain.RecipeRecommendation{}, nil
	}

	ranking, err := u.rakuten.GetCategoryRanking(ctx, category.ID, category.Name)
	if err != nil {
		return nil, err
	}

	result := make([]recipeDomain.RecipeRecommendation, 0, len(ranking))
	for _, item := range ranking {
		score := scoreFromRank(item.Rank)
		result = append(result, recipeDomain.RecipeRecommendation{
			ID:          item.ID,
			Title:       item.Title,
			ImageURL:    item.ImageURL,
			RecipeURL:   item.RecipeURL,
			CookingTime: item.CookingTime,
			Cost:        item.Cost,
			Score:       &score,
			Explain:     buildExplain(item.MatchedCategory, item.Description),
		})
	}

	return result, nil
}

func pickBestCategory(categories []RakutenCategory, fishName string) (RakutenCategory, bool) {
	fishAliases := buildFishAliases(fishName)
	bestScore := 0
	best := RakutenCategory{}

	for _, category := range categories {
		categoryAliases := buildFishAliases(category.Name)
		score := 0
		for _, fishAlias := range fishAliases {
			for _, categoryAlias := range categoryAliases {
				current := 0
				switch {
				case categoryAlias == fishAlias:
					current = 5
				case strings.Contains(categoryAlias, fishAlias):
					current = 4
				case strings.Contains(fishAlias, categoryAlias):
					current = 3
				case familyKeyword(fishAlias) != "" && familyKeyword(fishAlias) == categoryAlias:
					current = 2
				}
				if current > score {
					score = current
				}
			}
		}

		if score > bestScore || (score == bestScore && bestScore > 0 && len(category.Name) < len(best.Name)) {
			bestScore = score
			best = category
		}
	}

	return best, bestScore > 0
}

func normalize(value string) string {
	value = strings.TrimSpace(strings.ToLower(toHiragana(value)))
	replacer := strings.NewReplacer(" ", "", "　", "", "-", "", "（", "", "）", "", "(", "", ")", "")
	return replacer.Replace(value)
}

func scoreFromRank(rank int) int {
	score := 6 - rank
	if score < 1 {
		return 1
	}
	return score
}

func buildExplain(categoryName string, description string) string {
	base := "楽天レシピカテゴリ「" + categoryName + "」の人気レシピです。"
	description = strings.TrimSpace(description)
	if description == "" {
		return base
	}
	return base + " " + description
}

func buildFishAliases(value string) []string {
	normalized := normalize(value)
	if normalized == "" {
		return nil
	}

	aliases := []string{normalized}
	seen := map[string]struct{}{normalized: {}}
	add := func(candidate string) {
		candidate = normalize(candidate)
		if candidate == "" {
			return
		}
		if _, ok := seen[candidate]; ok {
			return
		}
		seen[candidate] = struct{}{}
		aliases = append(aliases, candidate)
	}

	if len([]rune(normalized)) >= 3 {
		switch {
		case strings.HasPrefix(normalized, "ま"):
			add(strings.TrimPrefix(normalized, "ま"))
		case strings.HasPrefix(normalized, "真"):
			add(strings.TrimPrefix(normalized, "真"))
		}
	}

	if strings.HasSuffix(normalized, "だい") {
		add(strings.TrimSuffix(normalized, "だい") + "たい")
		add("たい")
	}
	if strings.HasSuffix(normalized, "鯛") {
		add(strings.TrimSuffix(normalized, "鯛") + "たい")
		add("たい")
	}
	if strings.HasSuffix(normalized, "あじ") {
		add("あじ")
	}
	if strings.HasSuffix(normalized, "さば") {
		add("さば")
	}

	family := familyKeyword(normalized)
	if family != "" {
		add(family)
	}

	return aliases
}

func familyKeyword(value string) string {
	switch {
	case strings.HasSuffix(value, "だい"), strings.HasSuffix(value, "鯛"), strings.Contains(value, "たい"):
		return "たい"
	case strings.HasSuffix(value, "あじ"):
		return "あじ"
	case strings.HasSuffix(value, "さば"):
		return "さば"
	default:
		return ""
	}
}

func toHiragana(value string) string {
	return strings.Map(func(r rune) rune {
		if r >= 'ァ' && r <= 'ヶ' {
			return r - ('ァ' - 'ぁ')
		}
		if unicode.IsSpace(r) {
			return ' '
		}
		return r
	}, value)
}

func normalizeCondition(condition recipeDomain.SearchCondition) recipeDomain.SearchCondition {
	condition.FishName = strings.TrimSpace(condition.FishName)
	condition.Keyword = strings.TrimSpace(condition.Keyword)
	condition.UserID = strings.TrimSpace(condition.UserID)
	if condition.Count <= 0 {
		condition.Count = defaultCount
	}
	if condition.Count > maxCount {
		condition.Count = maxCount
	}
	if condition.Page <= 0 {
		condition.Page = defaultPage
	}

	return condition
}
