package handler

import (
	"net/http"

	"firebase-auth-go/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 魚の一覧を取得 (検索・フィルター・あいうえお順)
func GetFishList(c *gin.Context) {
	db := config.GetDB()
	var fishes []Fish

	// クエリパラメータ
	search := c.Query("search")
	month := c.Query("month")//旬
	species := c.Query("species") // 種類フィルター
	onlyLiked := c.Query("only_liked") == "true"
	uid, exists := c.Get("uid")

	query := db.WithContext(c.Request.Context()).Model(&model.Fish{})

	// お気に入り絞り込み (ログイン中)
	if onlyLiked && exists {
		query = query.Joins("JOIN user_fish_links ON user_fish_links.fish_id = fishes.id").
			Where("user_fish_links.user_id = ? AND user_fish_links.like = ?", uid.(string), true)
	}

	// 検索 (名前・学名)
	if search != "" {
		query = query.Where("name_ja LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 旬の月フィルター
	if month != "" {
		query = query.Joins("JOIN seasons ON seasons.fish_id = fishes.id").
			Where("seasons.month = ?", month)
	}

	// 種類フィルター
	if species != "" {
		query = query.Where("species = ?", species)
	}

	// あいうえお順で取得
	err := query.Distinct().Order("name_ja ASC").Find(&fishes).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "取得失敗"})
		return
	}

// お気に入りマップ
    likedMap := make(map[string]bool)
    if exists {
        var likedFishIDs []string
        db.Model(&model.UserFishLink{}).
            Where("user_id = ? AND `like` = ?", uid.(string), true).
            Pluck("fish_id", &likedFishIDs)
        for _, id := range likedFishIDs {
            likedMap[id] = true
        }
    }

    // レスポンス構築
    response := make([]model.FishListResponse, len(fishes))
    for i, f := range fishes {
        response[i] = model.FishListResponse{
            ID:       f.ID,
            NameJa:   f.NameJa,
            Name:     f.Name,
            ImageURL: f.ImageURL,
            IsLiked:  likedMap[f.ID],
        }
    }

    c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// 2. 魚の詳細情報を取得
func GetFishDetail(c *gin.Context) {
	id := c.Param("id")
	db := config.GetDB()

	var fish model.Fish
    err := db.WithContext(c.Request.Context()).
        Preload("Seasons").
        // 相性のいい魚：IsBad = false かつ Scoreが高い順
        Preload("CompatibleFish", func(db *gorm.DB) *gorm.DB {
            return db.Joins("JOIN fish_pair ON fish_pair.fish_b_id = fishes.id").
                Where("fish_pair.is_bad = ?", false).
                Order("fish_pair.score DESC").Limit(3)
        }).
        // 相性の悪い魚：IsBad = true
        Preload("IncompatibleFish", func(db *gorm.DB) *gorm.DB {
            return db.Joins("JOIN fish_pair ON fish_pair.fish_b_id = fishes.id").
                Where("fish_pair.is_bad = ?", true).
                Limit(3)
        }).
        First(&fish, "id = ?", id).Error

    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "魚が見つかりません"})
        return
    }

	// お気に入り状態の確認
	isLiked := false
	if uid, exists := c.Get("uid"); exists {
		var count int64
		db.Model(&model.UserFishLink{}).
			Where("user_id = ? AND fish_id = ? AND `like` = ?", uid.(string), id, true).
			Count(&count)
		isLiked = count > 0
	}

	// レスポンス
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"fish":     fish,    // 生息地、体長、種類、StoryURL、リレーションが含まれる
			"is_liked": isLiked,
		},
	})
}

// 3. お気に入りのトグル
func ToggleLike(c *gin.Context) {
	fishID := c.Param("id")
	uid, _ := c.Get("uid")
	db := config.GetDB()
	var like model.UserFishLink

	db.Where(model.UserFishLink{UserID: uid.(string), FishID: fishID}).FirstOrCreate(&like)
	
	if like.CreatedAt.IsZero() {
		like.Like = true
	} else {
		like.Like = !like.Like
	}
	db.Save(&like)

	c.JSON(http.StatusOK, gin.H{"success": true, "isLiked": like.Like})
}






//初期データ投入関数
func SeedData(c *gin.Context) {
	db := config.GetDB()

	// 1. 既存データのクリーンアップ（テスト用：必要に応じて）
	// db.Exec("DELETE FROM fish_pair")
	// db.Exec("DELETE FROM seasons")
	// db.Exec("DELETE FROM fishes")

	// 2. 魚のメインデータ作成
	// マグロ
	maguro := model.Fish{
		NameJa:   "クロマグロ",
		Name:     "Thunnus orientalis",
		Species:  "スズキ目サバ科",
		Habitat:  "太平洋、大西洋などの温帯・熱帯海域",
		Size:     "全長3m、体重400kg以上になることもある",
		Explain:  "「海のダイヤ」とも呼ばれる高級魚。時速80km以上で泳ぐことができます。",
		StoryURL: "https://example.com/story/maguro-vlog",
		ImageURL: "https://example.com/images/maguro.jpg",
	}

	// サンマ（相性がいい魚として）
	sanma := model.Fish{
		NameJa:   "サンマ",
		Name:     "Cololabis saira",
		Species:  "ダツ目サンマ科",
		Habitat:  "北太平洋の温帯海域",
		Size:     "全長約35cm程度",
		Explain:  "秋の味覚を代表する魚。銀色に輝く刀のような姿が特徴です。",
		StoryURL: "https://example.com/story/sanma-fishery",
		ImageURL: "https://example.com/images/sanma.jpg",
	}

	// 金魚（相性が悪い魚の例として）
	kingyo := model.Fish{
		NameJa:   "ワキン（金魚）",
		Name:     "Carassius auratus",
		Species:  "コイ目コイ科",
		Habitat:  "淡水（飼育下）",
		Size:     "15〜30cm程度",
		Explain:  "もっともポピュラーな金魚。マグロとは住む場所も性質も全く異なります。",
		StoryURL: "",
		ImageURL: "https://example.com/images/kingyo.jpg",
	}

	// データベースに保存（重複を避けるために FirstOrCreate を使用）
	db.FirstOrCreate(&maguro, model.Fish{NameJa: maguro.NameJa})
	db.FirstOrCreate(&sanma, model.Fish{NameJa: sanma.NameJa})
	db.FirstOrCreate(&kingyo, model.Fish{NameJa: kingyo.NameJa})

	// 3. 旬のデータ（マグロ：冬、サンマ：秋）
	maguroSeasons := []model.Season{
		{FishID: maguro.Base.ID, Month: 12},
		{FishID: maguro.Base.ID, Month: 1},
	}
	for _, s := range maguroSeasons {
		db.FirstOrCreate(&s, model.Season{FishID: s.FishID, Month: s.Month})
	}

	// 4. 相性データ（FishPair）
	pairs := []model.FishPair{
		{
			FishAID: maguro.Base.ID,
			FishBID: sanma.Base.ID,
			IsBad:   false, // 相性がいい
			Score:   90,
			Result:  "抜群",
			Explain: "同じ海域で獲れる旬の魚同士、食卓での相性も最高です。",
		},
		{
			FishAID: maguro.Base.ID,
			FishBID: kingyo.Base.ID,
			IsBad:   true, // 相性が悪い
			Score:   5,
			Result:  "最悪",
			Explain: "海水魚と淡水魚であり、マグロが食べてしまう恐れがあります。",
		},
	}
	for _, p := range pairs {
		// すでにペアが存在するか確認して保存
		db.Where(model.FishPair{FishAID: p.FishAID, FishBID: p.FishBID}).FirstOrCreate(&p)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "最新の構造体でデータを投入しました"})
}