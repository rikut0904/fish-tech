package repository

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"gorm.io/gorm"

	adminDomain "fish-tech/internal/domain/admin"
	"fish-tech/internal/infrastructure/persistence/gorm/model"
	adminUsecase "fish-tech/internal/usecase/admin"
)

// AdminRepository は管理画面向けのリポジトリです。
type AdminRepository struct {
	db *gorm.DB
}

// NewAdminRepository は新しい管理画面リポジトリを作成します。
func NewAdminRepository(db *gorm.DB) (*AdminRepository, error) {
	return &AdminRepository{db: db}, nil
}

// ListFishes は魚一覧を返します。
func (r *AdminRepository) ListFishes(ctx context.Context) ([]adminDomain.Fish, error) {
	var rows []model.Fish
	if err := r.db.WithContext(ctx).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]adminDomain.Fish, 0, len(rows))
	for _, row := range rows {
		name := strings.TrimSpace(row.NameJa)
		if name == "" {
			name = row.Name
		}

		result = append(result, adminDomain.Fish{
			ID:           row.ID,
			Name:         name,
			Category:     row.Category,
			Description:  row.Explain,
			ImageURL:     row.ImageURL,
			ImageMediaID: row.ImageMediaID,
			LinkURL:      row.LinkURL,
			CreatedAt:    row.CreatedAt,
		})
	}

	return result, nil
}

// CreateFish は魚を作成します。
func (r *AdminRepository) CreateFish(ctx context.Context, fish adminDomain.Fish) (adminDomain.Fish, error) {
	updatedAt := fish.CreatedAt
	row := model.Fish{
		ID:           fish.ID,
		NameJa:       fish.Name,
		Name:         fish.Name,
		Category:     fish.Category,
		Explain:      fish.Description,
		ImageURL:     fish.ImageURL,
		ImageMediaID: fish.ImageMediaID,
		LinkURL:      fish.LinkURL,
		CreatedAt:    fish.CreatedAt,
		UpdatedAt:    &updatedAt,
	}

	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return adminDomain.Fish{}, err
	}

	return fish, nil
}

// DeleteFish は魚を削除し、関連相性も削除します。
func (r *AdminRepository) DeleteFish(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("fish_a_id = ? OR fish_b_id = ?", id, id).Delete(&model.FishPair{}).Error; err != nil {
			return err
		}

		result := tx.Where("id = ?", id).Delete(&model.Fish{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return adminUsecase.ErrFishNotFound
		}

		return nil
	})
}

// ListPairs は魚相性一覧を返します。
func (r *AdminRepository) ListPairs(ctx context.Context, fishID string) ([]adminDomain.FishPair, error) {
	var rows []model.FishPair
	query := r.db.WithContext(ctx).Order("created_at desc")
	if fishID != "" {
		query = query.Where("fish_a_id = ? OR fish_b_id = ?", fishID, fishID)
	}

	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]adminDomain.FishPair, 0, len(rows))
	for _, row := range rows {
		fishIDa, fishIDb := normalizePairIDs(row.FishAID, row.FishBID)
		result = append(result, adminDomain.FishPair{
			ID:      buildPairID(fishIDa, fishIDb),
			FishIDa: fishIDa,
			FishIDb: fishIDb,
			Score:   row.Score,
			Memo:    row.Explain,
		})
	}

	return result, nil
}

// CreatePair は魚相性を作成します。
func (r *AdminRepository) CreatePair(ctx context.Context, pair adminDomain.FishPair) (adminDomain.FishPair, error) {
	fishIDa, fishIDb := normalizePairIDs(pair.FishIDa, pair.FishIDb)

	updatedAt := pair.CreatedAt
	row := model.FishPair{
		FishAID:   fishIDa,
		FishBID:   fishIDb,
		Result:    derivePairResult(pair.Score),
		Score:     pair.Score,
		Explain:   pair.Memo,
		CreatedAt: pair.CreatedAt,
		UpdatedAt: &updatedAt,
	}

	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return adminDomain.FishPair{}, err
	}

	pair.FishIDa = fishIDa
	pair.FishIDb = fishIDb
	pair.ID = buildPairID(fishIDa, fishIDb)
	return pair, nil
}

// DeletePair は魚相性を削除します。
func (r *AdminRepository) DeletePair(ctx context.Context, id string) error {
	fishIDa, fishIDb, err := parsePairID(id)
	if err != nil {
		return adminUsecase.ErrPairNotFound
	}
	fishIDa, fishIDb = normalizePairIDs(fishIDa, fishIDb)

	result := r.db.WithContext(ctx).Where("fish_a_id = ? AND fish_b_id = ?", fishIDa, fishIDb).Delete(&model.FishPair{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return adminUsecase.ErrPairNotFound
	}

	return nil
}

// ExistsFish は魚の存在確認を行います。
func (r *AdminRepository) ExistsFish(ctx context.Context, id string) (bool, error) {
	var row model.Fish
	err := r.db.WithContext(ctx).Select("id").Where("id = ?", id).Take(&row).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}

// ExistsPair は魚相性の存在確認を行います。
func (r *AdminRepository) ExistsPair(ctx context.Context, fishIDa string, fishIDb string) (bool, error) {
	var row model.FishPair
	err := r.db.WithContext(ctx).
		Select("fish_a_id").
		Where("(fish_a_id = ? AND fish_b_id = ?) OR (fish_a_id = ? AND fish_b_id = ?)", fishIDa, fishIDb, fishIDb, fishIDa).
		Take(&row).
		Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}

func buildPairID(fishIDa string, fishIDb string) string {
	fishIDa, fishIDb = normalizePairIDs(fishIDa, fishIDb)
	return fishIDa + ":" + fishIDb
}

func parsePairID(id string) (string, string, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid pair id")
	}

	fishIDa := strings.TrimSpace(parts[0])
	fishIDb := strings.TrimSpace(parts[1])
	if fishIDa == "" || fishIDb == "" {
		return "", "", fmt.Errorf("invalid pair id")
	}

	return fishIDa, fishIDb, nil
}

func normalizePairIDs(fishIDa string, fishIDb string) (string, string) {
	ids := []string{fishIDa, fishIDb}
	sort.Strings(ids)
	return ids[0], ids[1]
}

func derivePairResult(score int) string {
	switch {
	case score >= 4:
		return "good"
	case score <= 2:
		return "bad"
	default:
		return "neutral"
	}
}
