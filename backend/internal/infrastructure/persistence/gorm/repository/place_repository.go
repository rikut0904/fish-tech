package repository

import (
	"context"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	placeDomain "fish-tech/internal/domain/place"
	"fish-tech/internal/infrastructure/hotpepper"
	"fish-tech/internal/infrastructure/persistence/gorm/model"
)

const cacheValidDuration = time.Hour * 24 * 30
const asyncUpsertBatchSize = 100
const smallAreaCacheValidDuration = time.Hour * 24 * 365

type placeCacheRow struct {
	ID            string
	Name          string
	Address       *string
	Lat           *float64
	Lng           *float64
	Coupon        *string
	Genre         *string
	Card          *string
	Logo          *string
	SmallAreaCode *string
}

// PlaceRepository は店舗検索の永続化リポジトリです。
type PlaceRepository struct {
	db              *gorm.DB
	hotpepperClient *hotpepper.Client
}

// NewPlaceRepository は新しい店舗検索リポジトリを作成します。
func NewPlaceRepository(db *gorm.DB, hotpepperClient *hotpepper.Client) (*PlaceRepository, error) {
	return &PlaceRepository{db: db, hotpepperClient: hotpepperClient}, nil
}

// SearchRecommendedPlaces はキャッシュを優先し、必要時にHotPepper APIを再取得します。
func (r *PlaceRepository) SearchRecommendedPlaces(ctx context.Context, condition placeDomain.SearchCondition) ([]placeDomain.RecommendedPlace, error) {
	now := time.Now()
	freshFrom := now.Add(-cacheValidDuration)

	if condition.FishName != "" {
		fishID, found, err := r.findFishIDByName(ctx, condition.FishName)
		if err != nil {
			return nil, err
		}
		if found {
			cached, err := r.findFreshPlacesByFishID(ctx, fishID, condition.LargeArea, condition.CityCode, condition.UserID, condition.Favorite, freshFrom, condition.Count, condition.Page)
			if err != nil {
				return nil, err
			}
			if condition.Favorite {
				if len(cached) > 0 {
					return cached, nil
				}
			} else if len(cached) >= condition.Count {
				return cached, nil
			}

			fetched, err := r.hotpepperClient.SearchRecommendedPlaces(ctx, condition)
			if err != nil {
				return nil, err
			}
			r.persistFetchedPlacesAsync(fetched, fishID, now)
			if condition.Favorite {
				return cached, nil
			}
			return fetched, nil
		}
	}

	cached, err := r.findFreshPlacesByKeyword(ctx, condition.Keyword, condition.LargeArea, condition.CityCode, condition.UserID, condition.Favorite, freshFrom, condition.Count, condition.Page)
	if err != nil {
		return nil, err
	}
	if condition.Favorite {
		if len(cached) > 0 {
			return cached, nil
		}
	} else if len(cached) >= condition.Count {
		return cached, nil
	}

	fetched, err := r.hotpepperClient.SearchRecommendedPlaces(ctx, condition)
	if err != nil {
		return nil, err
	}
	r.persistFetchedPlacesAsync(fetched, "", now)
	if condition.Favorite {
		return cached, nil
	}

	return fetched, nil
}

func (r *PlaceRepository) findFishIDByName(ctx context.Context, fishName string) (string, bool, error) {
	var fish model.Fish
	err := r.db.WithContext(ctx).
		Where("name_ja = ? OR name = ?", fishName, fishName).
		First(&fish).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", false, nil
		}
		return "", false, err
	}

	return fish.ID, true, nil
}

// IsValidSmallAreaCode は石川県のsmall_areaコードとして有効かを返します。
func (r *PlaceRepository) IsValidSmallAreaCode(ctx context.Context, largeAreaCode string, smallAreaCode string) (bool, error) {
	trimmedLargeAreaCode := strings.TrimSpace(largeAreaCode)
	trimmedSmallAreaCode := strings.TrimSpace(smallAreaCode)
	if trimmedLargeAreaCode == "" || trimmedSmallAreaCode == "" {
		return false, nil
	}

	freshFrom := time.Now().Add(-smallAreaCacheValidDuration)
	freshExists, err := r.hasFreshSmallAreaCache(ctx, trimmedLargeAreaCode, freshFrom)
	if err != nil {
		return false, err
	}
	if !freshExists {
		fetched, fetchErr := r.hotpepperClient.FetchSmallAreasByLargeArea(ctx, trimmedLargeAreaCode)
		if fetchErr != nil {
			return false, fetchErr
		}
		if upsertErr := r.replaceSmallAreaCache(ctx, trimmedLargeAreaCode, fetched, time.Now()); upsertErr != nil {
			return false, upsertErr
		}
		for _, area := range fetched {
			if strings.EqualFold(strings.TrimSpace(area.Code), trimmedSmallAreaCode) {
				return true, nil
			}
		}
		return false, nil
	}

	var count int64
	if err := r.db.WithContext(ctx).
		Table("hotpepper_small_area_cache").
		Where("large_area_code = ? AND code = ? AND fetched_at >= ?", trimmedLargeAreaCode, trimmedSmallAreaCode, freshFrom).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *PlaceRepository) hasFreshSmallAreaCache(ctx context.Context, largeAreaCode string, freshFrom time.Time) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Table("hotpepper_small_area_cache").
		Where("large_area_code = ? AND fetched_at >= ?", largeAreaCode, freshFrom).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PlaceRepository) replaceSmallAreaCache(ctx context.Context, largeAreaCode string, areas []hotpepper.SmallArea, fetchedAt time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("hotpepper_small_area_cache").
			Where("large_area_code = ?", largeAreaCode).
			Delete(&model.HotpepperSmallAreaCache{}).Error; err != nil {
			return err
		}

		if len(areas) == 0 {
			return nil
		}

		rows := make([]model.HotpepperSmallAreaCache, 0, len(areas))
		for _, area := range areas {
			code := strings.TrimSpace(area.Code)
			if code == "" {
				continue
			}
			name := strings.TrimSpace(area.Name)
			middleArea := strings.TrimSpace(area.MiddleArea)
			row := model.HotpepperSmallAreaCache{
				LargeAreaCode: largeAreaCode,
				Code:          code,
				FetchedAt:     fetchedAt,
			}
			if name != "" {
				row.Name = &name
			}
			if middleArea != "" {
				row.MiddleArea = &middleArea
			}
			rows = append(rows, row)
		}
		if len(rows) == 0 {
			return nil
		}

		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "large_area_code"}, {Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "middle_area_code", "fetched_at"}),
		}).CreateInBatches(&rows, asyncUpsertBatchSize).Error
	})
}

// ExistsUser はuserの存在確認を行います。
func (r *PlaceRepository) ExistsUser(ctx context.Context, userID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Table("\"user\"").Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsPlace はplace_cacheの存在確認を行います。
func (r *PlaceRepository) ExistsPlace(ctx context.Context, placeID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Table("place_cache").Where("id = ?", placeID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdatePlaceFavorite はユーザーの店舗お気に入り状態を更新します。
func (r *PlaceRepository) UpdatePlaceFavorite(ctx context.Context, userID string, placeID string, favorite bool, now time.Time) error {
	row := model.UserPlaceLinks{
		UserID:    userID,
		PlaceID:   placeID,
		IsLikes:   favorite,
		CreatedAt: now,
		UpdatedAt: &now,
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "place_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"is_likes", "updated_at"}),
	}).Create(&row).Error
}

func (r *PlaceRepository) findFreshPlacesByFishID(ctx context.Context, fishID string, largeAreaCode string, cityCode string, userID string, favorite bool, freshFrom time.Time, count int, page int) ([]placeDomain.RecommendedPlace, error) {
	rows := make([]placeCacheRow, 0, count)
	offset := (page - 1) * count
	query := r.db.WithContext(ctx).
		Table("fish_place_links AS fpl").
		Select("pc.id, pc.name, pc.address, pc.lat, pc.lng, pc.coupon, pc.genre, pc.card, pc.logo, pc.small_area_code").
		Joins("JOIN place_cache AS pc ON pc.id = fpl.place_id").
		Where("fpl.fish_id = ? AND pc.fetched_at >= ?", fishID, freshFrom).
		Where("pc.large_area_code = ?", largeAreaCode)
	if favorite {
		query = query.Joins("JOIN user_place_links AS upl ON upl.place_id = pc.id").
			Where("upl.user_id = ? AND upl.is_likes = TRUE", userID)
	}

	trimmedCityCode := strings.TrimSpace(cityCode)
	if trimmedCityCode != "" {
		query = query.Where("pc.small_area_code = ?", trimmedCityCode)
	}

	err := query.Order("COALESCE(fpl.score, 0) DESC, pc.fetched_at DESC").
		Limit(count).
		Offset(offset).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return toRecommendedPlaces(rows), nil
}

func (r *PlaceRepository) findFreshPlacesByKeyword(ctx context.Context, keyword string, largeAreaCode string, cityCode string, userID string, favorite bool, freshFrom time.Time, count int, page int) ([]placeDomain.RecommendedPlace, error) {
	rows := make([]placeCacheRow, 0, count)
	offset := (page - 1) * count
	query := r.db.WithContext(ctx).
		Table("place_cache").
		Select("id, name, address, lat, lng, coupon, genre, card, logo, small_area_code").
		Where("fetched_at >= ?", freshFrom).
		Where("large_area_code = ?", largeAreaCode)
	if favorite {
		query = query.Joins("JOIN user_place_links AS upl ON upl.place_id = place_cache.id").
			Where("upl.user_id = ? AND upl.is_likes = TRUE", userID)
	}

	trimmedKeyword := strings.TrimSpace(keyword)
	if trimmedKeyword != "" && trimmedKeyword != "魚" {
		likePattern := "%" + trimmedKeyword + "%"
		query = query.Where("name ILIKE ? OR address ILIKE ? OR genre ILIKE ?", likePattern, likePattern, likePattern)
	}
	trimmedCityCode := strings.TrimSpace(cityCode)
	if trimmedCityCode != "" {
		query = query.Where("small_area_code = ?", trimmedCityCode)
	}

	err := query.Order("fetched_at DESC").
		Limit(count).
		Offset(offset).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return toRecommendedPlaces(rows), nil
}

func (r *PlaceRepository) upsertPlaces(ctx context.Context, places []placeDomain.RecommendedPlace, fetchedAt time.Time) error {
	rows := make([]model.PlaceCache, 0, len(places))
	for _, place := range places {
		var address *string
		if strings.TrimSpace(place.Address) != "" {
			value := strings.TrimSpace(place.Address)
			address = &value
		}
		var lat *float64
		if parsedLat, err := strconv.ParseFloat(strings.TrimSpace(place.Lat), 64); err == nil {
			lat = &parsedLat
		}
		var lng *float64
		if parsedLng, err := strconv.ParseFloat(strings.TrimSpace(place.Lng), 64); err == nil {
			lng = &parsedLng
		}
		var coupon *string
		if strings.TrimSpace(place.Coupon) != "" {
			value := strings.TrimSpace(place.Coupon)
			coupon = &value
		}

		var genre *string
		if strings.TrimSpace(place.Genre) != "" {
			value := strings.TrimSpace(place.Genre)
			genre = &value
		}
		var card *string
		if strings.TrimSpace(place.Card) != "" {
			value := strings.TrimSpace(place.Card)
			card = &value
		}
		var logo *string
		if strings.TrimSpace(place.Logo) != "" {
			value := strings.TrimSpace(place.Logo)
			logo = &value
		}
		var largeAreaCode *string
		if strings.TrimSpace(place.LargeAreaCode) != "" {
			value := strings.TrimSpace(place.LargeAreaCode)
			largeAreaCode = &value
		}
		var smallAreaCode *string
		if strings.TrimSpace(place.SmallAreaCode) != "" {
			value := strings.TrimSpace(place.SmallAreaCode)
			smallAreaCode = &value
		}

		row := model.PlaceCache{
			ID:            place.ID,
			Name:          place.Name,
			Address:       address,
			Lat:           lat,
			Lng:           lng,
			Coupon:        coupon,
			Genre:         genre,
			Card:          card,
			Logo:          logo,
			LargeAreaCode: largeAreaCode,
			SmallAreaCode: smallAreaCode,
			FetchedAt:     fetchedAt,
		}
		rows = append(rows, row)
	}

	if len(rows) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "address", "lat", "lng", "coupon", "genre", "card", "logo", "large_area_code", "small_area_code", "fetched_at"}),
	}).CreateInBatches(&rows, asyncUpsertBatchSize).Error; err != nil {
		return err
	}

	return nil
}

func (r *PlaceRepository) upsertFishPlaceLinks(ctx context.Context, fishID string, places []placeDomain.RecommendedPlace, now time.Time) error {
	links := make([]model.FishPlaceLinks, 0, len(places))
	for _, place := range places {
		link := model.FishPlaceLinks{
			FishID:    fishID,
			PlaceID:   place.ID,
			CreatedAt: now,
			UpdatedAt: &now,
		}
		links = append(links, link)
	}

	if len(links) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "fish_id"}, {Name: "place_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
	}).CreateInBatches(&links, asyncUpsertBatchSize).Error; err != nil {
		return err
	}

	return nil
}

func (r *PlaceRepository) persistFetchedPlacesAsync(places []placeDomain.RecommendedPlace, fishID string, now time.Time) {
	if len(places) == 0 {
		return
	}

	placesCopy := make([]placeDomain.RecommendedPlace, len(places))
	copy(placesCopy, places)

	go func() {
		bgCtx := context.Background()
		if err := r.upsertPlaces(bgCtx, placesCopy, now); err != nil {
			return
		}
		if strings.TrimSpace(fishID) != "" {
			_ = r.upsertFishPlaceLinks(bgCtx, fishID, placesCopy, now)
		}
	}()
}

func toRecommendedPlaces(rows []placeCacheRow) []placeDomain.RecommendedPlace {
	places := make([]placeDomain.RecommendedPlace, 0, len(rows))
	for _, row := range rows {
		address := ""
		if row.Address != nil {
			address = *row.Address
		}
		lat := ""
		if row.Lat != nil {
			lat = strconv.FormatFloat(*row.Lat, 'f', -1, 64)
		}
		lng := ""
		if row.Lng != nil {
			lng = strconv.FormatFloat(*row.Lng, 'f', -1, 64)
		}
		coupon := ""
		if row.Coupon != nil {
			coupon = *row.Coupon
		}
		genre := ""
		if row.Genre != nil {
			genre = *row.Genre
		}
		card := ""
		if row.Card != nil {
			card = *row.Card
		}
		logo := ""
		if row.Logo != nil {
			logo = *row.Logo
		}
		smallAreaCode := ""
		if row.SmallAreaCode != nil {
			smallAreaCode = *row.SmallAreaCode
		}

		places = append(places, placeDomain.RecommendedPlace{
			ID:            row.ID,
			Name:          row.Name,
			Address:       address,
			Lat:           lat,
			Lng:           lng,
			Coupon:        coupon,
			Genre:         genre,
			Card:          card,
			Logo:          logo,
			SmallAreaCode: smallAreaCode,
		})
	}

	return places
}
