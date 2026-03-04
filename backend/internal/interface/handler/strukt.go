package handler

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 共通のフィールドを定義
type Base struct {
	ID        string         `gorm:"primaryKey;type:uuid" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// UUIDを自動生成します
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return
}

// ユーザー情報
type User struct {
	Base
	FirebaseUID string `gorm:"uniqueIndex;not null" json:"firebaseUid"`
	Name        string `gorm:"not null" json:"name"`
	Mail        string `gorm:"not null" json:"mail"`
}

// 魚の基本情報
type Fish struct {
	Base
	NameJa  string `gorm:"not null" json:"nameJa"`  // 和名
	Name    string `gorm:"not null" json:"name"`    // 学名
	Species  string `gorm:"not null" json:"species"`//種類
	Habitat  string `gorm:"type:text" json:"habitat"`//生息地
	Size     string `json:"size"`//体長
	StoryURL string `json:"storyUrl"`//ストーリーURL
	Explain string `gorm:"type:text" json:"explain"` // 説明
	ImageURL string `json:"imageUrl"`               // 写真URL（追加）

	// リレーション
	Seasons         []Season         `gorm:"foreignKey:FishID" json:"seasons"`
	FishingMethods  []FishingMethod  `gorm:"many2many:fishing_methods_links;" json:"fishingMethods"`
	Recipes         []RecipeCache    `gorm:"many2many:fish_recipe_links;" json:"recipes"`
	CompatibleFish  []Fish           `gorm:"many2many:fish_pair;joinForeignKey:FishAID;joinReferences:FishBID" json:"compatibleFish"`
	//相性の悪い魚 (Scoreが低いものを想定)
    IncompatibleFish []Fish 		 `gorm:"many2many:fish_pair;joinForeignKey:FishAID;joinReferences:FishBID" json:"incompatibleFish"`
}

// 旬の月
type Season struct {
	FishID    string    `gorm:"primaryKey" json:"fishId"`
	Month     int       `gorm:"primaryKey" json:"month"` // 1-12
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 漁法
type FishingMethod struct {
	Base
	Name    string `gorm:"uniqueIndex;not null" json:"name"`
	Explain string `gorm:"type:text" json:"explain"`
}

// レシピ
type RecipeCache struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	ImageURL    string    `json:"imageUrl"`
	RecipeURL   string    `gorm:"not null" json:"recipeUrl"`
	CookingTime string    `json:"cookingTime"`
	Cost        string    `json:"cost"`
	FetchedAt   time.Time `json:"fetchedAt"`
}

// ユーザーのお気に入り魚
type UserFishLink struct {
	UserID    string    `gorm:"primaryKey" json:"userId"`
	FishID    string    `gorm:"primaryKey" json:"fishId"`
	Like      bool      `gorm:"default:true" json:"like"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 相性のいい魚のペア
type FishPair struct {
	FishAID   string `gorm:"primaryKey" json:"fishAId"`
	FishBID   string `gorm:"primaryKey" json:"fishBId"`
	// true = 相性が悪い / false = 相性がいい
    IsBad     bool   `gorm:"default:false" json:"isBad"`
	Result    string `json:"result"`
	Explain   string `gorm:"type:text" json:"explain"`
	Score     int    `json:"score"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// APIレスポンス用構造体

// 一覧表示用の軽量な構造体
type FishListResponse struct {
	ID       string `json:"id"`
	NameJa   string `json:"nameJa"`
	Name     string `json:"name"`    // 学名
	ImageURL string `json:"imageUrl"`
	IsLiked  bool   `json:"isLiked"` 
}

// 詳細表示用の構造体
type FishDetailResponse struct {
	Fish
	IsLiked bool `json:"isLiked"`
}