package admin

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/google/uuid"

	adminDomain "fish-tech/internal/domain/admin"
	"fish-tech/internal/shared/timeutil"
)

var (
	// ErrInvalidFishName は魚名が不正な場合のエラーです。
	ErrInvalidFishName = errors.New("魚名は必須です")
	// ErrInvalidFishCategory はカテゴリが不正な場合のエラーです。
	ErrInvalidFishCategory = errors.New("カテゴリは必須です")
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
	ListPairs(ctx context.Context, fishID string) ([]adminDomain.FishPair, error)
	CreatePair(ctx context.Context, pair adminDomain.FishPair) (adminDomain.FishPair, error)
	DeletePair(ctx context.Context, id string) error
	ExistsFish(ctx context.Context, id string) (bool, error)
	ExistsPair(ctx context.Context, fishIDa string, fishIDb string) (bool, error)
}

// MediaURLResolver は画像URL解決機能のインターフェースです。
type MediaURLResolver interface {
	ResolveMediaItemURL(ctx context.Context, mediaItemID string) (string, error)
}

// UseCase は管理画面向けユースケースのインターフェースです。
type UseCase interface {
	ListFishes(ctx context.Context) ([]adminDomain.Fish, error)
	CreateFish(ctx context.Context, name string, category string, description string, imageURL string, imageMediaID string, linkURL string) (adminDomain.Fish, error)
	DeleteFish(ctx context.Context, id string) error
	ListPairs(ctx context.Context, fishID string) ([]adminDomain.FishPair, error)
	CreatePair(ctx context.Context, fishIDa string, fishIDb string, score int, memo string) (adminDomain.FishPair, error)
	DeletePair(ctx context.Context, id string) error
}

type adminUseCase struct {
	repo     Repository
	resolver MediaURLResolver
}

// NewAdminUseCase は管理画面向けユースケースを作成します。
func NewAdminUseCase(repo Repository) UseCase {
	return NewAdminUseCaseWithResolver(repo, nil)
}

// NewAdminUseCaseWithResolver は画像URL解決機能付きの管理画面向けユースケースを作成します。
func NewAdminUseCaseWithResolver(repo Repository, resolver MediaURLResolver) UseCase {
	return &adminUseCase{repo: repo, resolver: resolver}
}

// ListFishes は魚一覧を返します。
func (u *adminUseCase) ListFishes(ctx context.Context) ([]adminDomain.Fish, error) {
	fishes, err := u.repo.ListFishes(ctx)
	if err != nil {
		return nil, err
	}

	if u.resolver == nil {
		return fishes, nil
	}

	for i := range fishes {
		if fishes[i].ImageMediaID == "" {
			continue
		}
		resolvedURL, resolveErr := u.resolver.ResolveMediaItemURL(ctx, fishes[i].ImageMediaID)
		if resolveErr != nil {
			continue
		}
		fishes[i].ImageURL = resolvedURL
	}

	return fishes, nil
}

// CreateFish は魚を追加します。
func (u *adminUseCase) CreateFish(ctx context.Context, name string, category string, description string, imageURL string, imageMediaID string, linkURL string) (adminDomain.Fish, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return adminDomain.Fish{}, ErrInvalidFishName
	}
	trimmedCategory := strings.TrimSpace(category)
	if trimmedCategory == "" {
		return adminDomain.Fish{}, ErrInvalidFishCategory
	}

	fish := adminDomain.Fish{
		ID:           uuid.NewString(),
		Name:         trimmedName,
		Category:     trimmedCategory,
		Description:  strings.TrimSpace(description),
		ImageURL:     strings.TrimSpace(imageURL),
		ImageMediaID: strings.TrimSpace(imageMediaID),
		LinkURL:      strings.TrimSpace(linkURL),
		CreatedAt:    timeutil.NowJST(),
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
func (u *adminUseCase) ListPairs(ctx context.Context, fishID string) ([]adminDomain.FishPair, error) {
	return u.repo.ListPairs(ctx, strings.TrimSpace(fishID))
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
		CreatedAt: timeutil.NowJST(),
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
