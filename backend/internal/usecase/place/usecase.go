package place

import (
	"context"
	"errors"
	"strings"
	"time"

	placeDomain "fish-tech/internal/domain/place"
)

const (
	defaultKeyword        = "魚"
	defaultCount          = 10
	defaultPage           = 1
	maxCount              = 100
	ishikawaLargeAreaCode = "Z063"
)

// ErrSmallAreaNotFound は指定された市コードが石川県のsmall_areaに存在しない場合のエラーです。
var ErrSmallAreaNotFound = errors.New("指定された市コードは石川県のsmall_areaに存在しません")
var ErrFavoriteUserIDRequired = errors.New("userIdは必須です")
var ErrFavoritePlaceIDRequired = errors.New("placeIdは必須です")
var ErrFavoriteUserNotFound = errors.New("指定されたuserIdは存在しません")
var ErrFavoritePlaceNotFound = errors.New("指定されたplaceIdは存在しません")

// Repository は店舗検索のデータ取得インターフェースです。
type Repository interface {
	SearchRecommendedPlaces(ctx context.Context, condition placeDomain.SearchCondition) ([]placeDomain.RecommendedPlace, error)
	IsValidSmallAreaCode(ctx context.Context, largeAreaCode string, smallAreaCode string) (bool, error)
	ExistsUser(ctx context.Context, userID string) (bool, error)
	ExistsPlace(ctx context.Context, placeID string) (bool, error)
	UpdatePlaceFavorite(ctx context.Context, userID string, placeID string, favorite bool, now time.Time) error
}

// UseCase はおすすめ店舗検索ユースケースのインターフェースです。
type UseCase interface {
	SearchRecommendedPlaces(ctx context.Context, condition placeDomain.SearchCondition) ([]placeDomain.RecommendedPlace, error)
	UpdatePlaceFavorite(ctx context.Context, userID string, placeID string, favorite bool) error
}

type placeUseCase struct {
	repo Repository
}

// NewPlaceUseCase はおすすめ店舗検索ユースケースを作成します。
func NewPlaceUseCase(repo Repository) UseCase {
	return &placeUseCase{repo: repo}
}

// SearchRecommendedPlaces はおすすめ店舗を検索します。
func (u *placeUseCase) SearchRecommendedPlaces(ctx context.Context, condition placeDomain.SearchCondition) ([]placeDomain.RecommendedPlace, error) {
	condition.FishName = strings.TrimSpace(condition.FishName)
	condition.Keyword = strings.TrimSpace(condition.Keyword)
	condition.CityCode = strings.ToUpper(strings.TrimSpace(condition.CityCode))
	condition.UserID = strings.TrimSpace(condition.UserID)
	condition.LargeArea = ishikawaLargeAreaCode
	condition.SmallAreaCode = condition.CityCode

	if condition.Keyword == "" {
		if condition.FishName != "" {
			condition.Keyword = condition.FishName + " 魚料理"
		} else {
			condition.Keyword = defaultKeyword
		}
	}
	if !strings.Contains(condition.Keyword, defaultKeyword) {
		condition.Keyword = defaultKeyword + " " + condition.Keyword
	}
	if condition.Count <= 0 {
		condition.Count = defaultCount
	}
	if condition.Count > maxCount {
		condition.Count = maxCount
	}
	if condition.Page <= 0 {
		condition.Page = defaultPage
	}
	if condition.SmallAreaCode != "" {
		isValid, err := u.repo.IsValidSmallAreaCode(ctx, condition.LargeArea, condition.SmallAreaCode)
		if err != nil {
			return nil, err
		}
		if !isValid {
			return nil, ErrSmallAreaNotFound
		}
	}

	return u.repo.SearchRecommendedPlaces(ctx, condition)
}

// UpdatePlaceFavorite は店舗のお気に入り状態を更新します。
func (u *placeUseCase) UpdatePlaceFavorite(ctx context.Context, userID string, placeID string, favorite bool) error {
	trimmedUserID := strings.TrimSpace(userID)
	if trimmedUserID == "" {
		return ErrFavoriteUserIDRequired
	}
	trimmedPlaceID := strings.TrimSpace(placeID)
	if trimmedPlaceID == "" {
		return ErrFavoritePlaceIDRequired
	}

	existsUser, err := u.repo.ExistsUser(ctx, trimmedUserID)
	if err != nil {
		return err
	}
	if !existsUser {
		return ErrFavoriteUserNotFound
	}

	existsPlace, err := u.repo.ExistsPlace(ctx, trimmedPlaceID)
	if err != nil {
		return err
	}
	if !existsPlace {
		return ErrFavoritePlaceNotFound
	}

	return u.repo.UpdatePlaceFavorite(ctx, trimmedUserID, trimmedPlaceID, favorite, time.Now())
}
