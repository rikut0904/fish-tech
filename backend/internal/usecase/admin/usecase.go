package admin

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	adminDomain "fish-tech/internal/domain/admin"
)

var (
	// ErrInvalidFishName は魚名が不正な場合のエラーです。
	ErrInvalidFishName = errors.New("魚名は必須です")
	// ErrInvalidPair は相性データが不正な場合のエラーです。
	ErrInvalidPair = errors.New("相性データが不正です")
	// ErrPairAlreadyExists は同一ペアが既に存在する場合のエラーです。
	ErrPairAlreadyExists = errors.New("同じ魚ペアは既に登録されています")
	// ErrFishNotFound は指定魚が存在しない場合のエラーです。
	ErrFishNotFound = errors.New("魚が見つかりません")
	// ErrPairNotFound は指定相性データが存在しない場合のエラーです。
	ErrPairNotFound = errors.New("相性データが見つかりません")
)

// Repository は管理画面向けデータアクセスのインターフェースです。
type Repository interface {
	ListFishes(ctx context.Context) ([]adminDomain.Fish, error)
	CreateFish(ctx context.Context, fish adminDomain.Fish) (adminDomain.Fish, error)
	DeleteFish(ctx context.Context, id string) error
	ListPairs(ctx context.Context) ([]adminDomain.FishPair, error)
	CreatePair(ctx context.Context, pair adminDomain.FishPair) (adminDomain.FishPair, error)
	DeletePair(ctx context.Context, id string) error
	ExistsFish(ctx context.Context, id string) (bool, error)
	ExistsPair(ctx context.Context, fishIDa string, fishIDb string) (bool, error)
}

// UseCase は管理画面向けユースケースのインターフェースです。
type UseCase interface {
	ListFishes(ctx context.Context) ([]adminDomain.Fish, error)
	CreateFish(ctx context.Context, name string, category string, description string, imageURL string, linkURL string) (adminDomain.Fish, error)
	DeleteFish(ctx context.Context, id string) error
	ListPairs(ctx context.Context) ([]adminDomain.FishPair, error)
	CreatePair(ctx context.Context, fishIDa string, fishIDb string, score int, memo string) (adminDomain.FishPair, error)
	DeletePair(ctx context.Context, id string) error
}

type adminUseCase struct {
	repo Repository
}

// NewAdminUseCase は管理画面向けユースケースを作成します。
func NewAdminUseCase(repo Repository) UseCase {
	return &adminUseCase{repo: repo}
}

// ListFishes は魚一覧を返します。
func (u *adminUseCase) ListFishes(ctx context.Context) ([]adminDomain.Fish, error) {
	return u.repo.ListFishes(ctx)
}

// CreateFish は魚を追加します。
func (u *adminUseCase) CreateFish(ctx context.Context, name string, category string, description string, imageURL string, linkURL string) (adminDomain.Fish, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return adminDomain.Fish{}, ErrInvalidFishName
	}

	fish := adminDomain.Fish{
		ID:          uuid.NewString(),
		Name:        trimmedName,
		Category:    strings.TrimSpace(category),
		Description: strings.TrimSpace(description),
		ImageURL:    strings.TrimSpace(imageURL),
		LinkURL:     strings.TrimSpace(linkURL),
		CreatedAt:   time.Now(),
	}

	return u.repo.CreateFish(ctx, fish)
}

// DeleteFish は魚を削除し、関連する相性データも削除します。
func (u *adminUseCase) DeleteFish(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrFishNotFound
	}

	return u.repo.DeleteFish(ctx, id)
}

// ListPairs は魚相性一覧を返します。
func (u *adminUseCase) ListPairs(ctx context.Context) ([]adminDomain.FishPair, error) {
	return u.repo.ListPairs(ctx)
}

// CreatePair は魚相性を追加します。
func (u *adminUseCase) CreatePair(ctx context.Context, fishIDa string, fishIDb string, score int, memo string) (adminDomain.FishPair, error) {
	fishIDa = strings.TrimSpace(fishIDa)
	fishIDb = strings.TrimSpace(fishIDb)
	if fishIDa == "" || fishIDb == "" || fishIDa == fishIDb || score < 1 || score > 5 {
		return adminDomain.FishPair{}, ErrInvalidPair
	}

	normalizedA, normalizedB := normalizePairIDs(fishIDa, fishIDb)

	existsA, err := u.repo.ExistsFish(ctx, normalizedA)
	if err != nil {
		return adminDomain.FishPair{}, err
	}
	if !existsA {
		return adminDomain.FishPair{}, ErrFishNotFound
	}

	existsB, err := u.repo.ExistsFish(ctx, normalizedB)
	if err != nil {
		return adminDomain.FishPair{}, err
	}
	if !existsB {
		return adminDomain.FishPair{}, ErrFishNotFound
	}

	existsPair, err := u.repo.ExistsPair(ctx, normalizedA, normalizedB)
	if err != nil {
		return adminDomain.FishPair{}, err
	}
	if existsPair {
		return adminDomain.FishPair{}, ErrPairAlreadyExists
	}

	pair := adminDomain.FishPair{
		ID:        uuid.NewString(),
		FishIDa:   normalizedA,
		FishIDb:   normalizedB,
		Score:     score,
		Memo:      strings.TrimSpace(memo),
		CreatedAt: time.Now(),
	}

	return u.repo.CreatePair(ctx, pair)
}

// DeletePair は魚相性を削除します。
func (u *adminUseCase) DeletePair(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrPairNotFound
	}

	return u.repo.DeletePair(ctx, id)
}

func normalizePairIDs(a string, b string) (string, string) {
	ids := []string{a, b}
	sort.Strings(ids)
	return ids[0], ids[1]
}
