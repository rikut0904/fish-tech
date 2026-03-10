package gorm

import (
	"fmt"

	"gorm.io/gorm"

	"fish-tech/internal/infrastructure/persistence/gorm/model"
)

// AutoMigrateAll は全テーブルを作成・更新します。
func AutoMigrateAll(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Fish{},
		&model.Season{},
		&model.FishStory{},
		&model.FishingMethods{},
		&model.FishingMethodsLinks{},
		&model.RecipeCache{},
		&model.FishRecipeLinks{},
		&model.PlaceCache{},
		&model.FishPlaceLinks{},
		&model.Diary{},
		&model.FishPair{},
		&model.UserFishLinks{},
		&model.FishUserLinks{},
		&model.UserPlaceLinks{},
		&model.UserRecipeLinks{},
		&model.HotpepperSmallAreaCache{},
		&model.AdminFish{},
		&model.AdminFishPair{},
	); err != nil {
		return err
	}

	if err := migratePlaceCacheLogo(db); err != nil {
		return err
	}
	if err := migratePlaceCacheTel(db); err != nil {
		return err
	}
	if err := migratePlaceCacheAreaAndIndexes(db); err != nil {
		return err
	}
	if err := migrateHotpepperSmallAreaCacheIndexes(db); err != nil {
		return err
	}
	if err := migrateUserFishLinks(db); err != nil {
		return err
	}
	if err := migrateFishUserLinks(db); err != nil {
		return err
	}
	if err := migrateUserPlaceLinks(db); err != nil {
		return err
	}

	return nil
}

func migratePlaceCacheAreaAndIndexes(db *gorm.DB) error {
	if !db.Migrator().HasTable(&model.PlaceCache{}) {
		return nil
	}

	if !db.Migrator().HasColumn(&model.PlaceCache{}, "large_area_code") {
		if err := db.Migrator().AddColumn(&model.PlaceCache{}, "large_area_code"); err != nil {
			return fmt.Errorf("place_cache.large_area_code の追加に失敗しました: %w", err)
		}
	}

	if err := db.Exec(`
		UPDATE place_cache
		SET large_area_code = 'Z063'
		WHERE large_area_code IS NULL
		  AND address ILIKE '%石川県%'
	`).Error; err != nil {
		return fmt.Errorf("place_cache.large_area_code の初期補完に失敗しました: %w", err)
	}

	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_place_cache_large_area_fetched_at ON place_cache (large_area_code, fetched_at DESC)`).Error; err != nil {
		return fmt.Errorf("place_cache.large_area_code,fetched_at インデックス作成に失敗しました: %w", err)
	}
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_place_cache_small_area_fetched_at ON place_cache (small_area_code, fetched_at DESC)`).Error; err != nil {
		return fmt.Errorf("place_cache.small_area_code,fetched_at インデックス作成に失敗しました: %w", err)
	}

	return nil
}

func migratePlaceCacheLogo(db *gorm.DB) error {
	if !db.Migrator().HasTable(&model.PlaceCache{}) {
		return nil
	}

	if !db.Migrator().HasColumn(&model.PlaceCache{}, "logo") {
		if err := db.Migrator().AddColumn(&model.PlaceCache{}, "logo"); err != nil {
			return fmt.Errorf("place_cache.logo の追加に失敗しました: %w", err)
		}
	}

	if db.Migrator().HasColumn(&model.PlaceCache{}, "photo_pc") || db.Migrator().HasColumn(&model.PlaceCache{}, "photo_mobile") {
		if err := db.Exec(`
			UPDATE place_cache
			SET logo = COALESCE(NULLIF(TRIM(logo), ''), NULLIF(TRIM(photo_pc), ''), NULLIF(TRIM(photo_mobile), ''))
			WHERE logo IS NULL OR TRIM(logo) = ''
		`).Error; err != nil {
			return fmt.Errorf("place_cache.logo へのデータ移行に失敗しました: %w", err)
		}
	}

	if db.Migrator().HasColumn(&model.PlaceCache{}, "photo_pc") {
		if err := db.Migrator().DropColumn(&model.PlaceCache{}, "photo_pc"); err != nil {
			return fmt.Errorf("place_cache.photo_pc の削除に失敗しました: %w", err)
		}
	}

	if db.Migrator().HasColumn(&model.PlaceCache{}, "photo_mobile") {
		if err := db.Migrator().DropColumn(&model.PlaceCache{}, "photo_mobile"); err != nil {
			return fmt.Errorf("place_cache.photo_mobile の削除に失敗しました: %w", err)
		}
	}

	return nil
}

func migratePlaceCacheTel(db *gorm.DB) error {
	if !db.Migrator().HasTable(&model.PlaceCache{}) {
		return nil
	}

	if db.Migrator().HasColumn(&model.PlaceCache{}, "tel") {
		if err := db.Migrator().DropColumn(&model.PlaceCache{}, "tel"); err != nil {
			return fmt.Errorf("place_cache.tel の削除に失敗しました: %w", err)
		}
	}

	return nil
}

func migrateHotpepperSmallAreaCacheIndexes(db *gorm.DB) error {
	if !db.Migrator().HasTable(&model.HotpepperSmallAreaCache{}) {
		return nil
	}

	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_hsac_large_area_fetched_at ON hotpepper_small_area_cache (large_area_code, fetched_at DESC)`).Error; err != nil {
		return fmt.Errorf("hotpepper_small_area_cache インデックス作成に失敗しました: %w", err)
	}

	return nil
}

func migrateFishUserLinks(db *gorm.DB) error {
	if !db.Migrator().HasTable(&model.FishUserLinks{}) {
		return nil
	}

	if err := dropLegacyFishUserLinksAutomation(db); err != nil {
		return err
	}
	if err := ensureForeignKeyConstraint(db, `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_fish_user_links_fish') THEN
				ALTER TABLE fish_user_links
				ADD CONSTRAINT fk_fish_user_links_fish
				FOREIGN KEY (fish_id) REFERENCES fish(id) ON DELETE CASCADE;
			END IF;
		END $$;
	`); err != nil {
		return fmt.Errorf("fish_user_links.fish_id FK制約の追加に失敗しました: %w", err)
	}
	if err := ensureForeignKeyConstraint(db, `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_fish_user_links_user') THEN
				ALTER TABLE fish_user_links
				ADD CONSTRAINT fk_fish_user_links_user
				FOREIGN KEY (user_id) REFERENCES "user"(user_id) ON DELETE CASCADE;
			END IF;
		END $$;
	`); err != nil {
		return fmt.Errorf("fish_user_links.user_id FK制約の追加に失敗しました: %w", err)
	}

	return nil
}

func migrateUserPlaceLinks(db *gorm.DB) error {
	if !db.Migrator().HasTable(&model.UserPlaceLinks{}) {
		return nil
	}

	if err := dropLegacyUserPlaceLinksAutomation(db); err != nil {
		return err
	}
	if err := ensureForeignKeyConstraint(db, `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_user_place_links_user') THEN
				ALTER TABLE user_place_links
				ADD CONSTRAINT fk_user_place_links_user
				FOREIGN KEY (user_id) REFERENCES "user"(user_id) ON DELETE CASCADE;
			END IF;
		END $$;
	`); err != nil {
		return fmt.Errorf("user_place_links.user_id FK制約の追加に失敗しました: %w", err)
	}
	if err := ensureForeignKeyConstraint(db, `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_user_place_links_place') THEN
				ALTER TABLE user_place_links
				ADD CONSTRAINT fk_user_place_links_place
				FOREIGN KEY (place_id) REFERENCES place_cache(id) ON DELETE CASCADE;
			END IF;
		END $$;
	`); err != nil {
		return fmt.Errorf("user_place_links.place_id FK制約の追加に失敗しました: %w", err)
	}

	return nil
}

func migrateUserFishLinks(db *gorm.DB) error {
	if !db.Migrator().HasTable(&model.UserFishLinks{}) {
		return nil
	}

	if err := dropLegacyUserFishLinksAutomation(db); err != nil {
		return err
	}
	if err := ensureForeignKeyConstraint(db, `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_user_fish_links_user') THEN
				ALTER TABLE user_fish_links
				ADD CONSTRAINT fk_user_fish_links_user
				FOREIGN KEY (user_id) REFERENCES "user"(user_id) ON DELETE CASCADE;
			END IF;
		END $$;
	`); err != nil {
		return fmt.Errorf("user_fish_links.user_id FK制約の追加に失敗しました: %w", err)
	}
	if err := ensureForeignKeyConstraint(db, `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_user_fish_links_fish') THEN
				ALTER TABLE user_fish_links
				ADD CONSTRAINT fk_user_fish_links_fish
				FOREIGN KEY (fish_id) REFERENCES fish(id) ON DELETE CASCADE;
			END IF;
		END $$;
	`); err != nil {
		return fmt.Errorf("user_fish_links.fish_id FK制約の追加に失敗しました: %w", err)
	}

	return nil
}

func dropLegacyFishUserLinksAutomation(db *gorm.DB) error {
	queries := []string{
		`DROP TRIGGER IF EXISTS trg_sync_fish_user_links_new_fish ON fish;`,
		`DROP TRIGGER IF EXISTS trg_sync_fish_user_links_new_user ON "user";`,
		`DROP TRIGGER IF EXISTS trg_cleanup_fish_user_links_on_fish_delete ON fish;`,
		`DROP TRIGGER IF EXISTS trg_cleanup_fish_user_links_on_user_delete ON "user";`,
		`DROP FUNCTION IF EXISTS sync_fish_user_links_new_fish();`,
		`DROP FUNCTION IF EXISTS sync_fish_user_links_new_user();`,
		`DROP FUNCTION IF EXISTS cleanup_fish_user_links_on_fish_delete();`,
		`DROP FUNCTION IF EXISTS cleanup_fish_user_links_on_user_delete();`,
	}
	for _, query := range queries {
		if err := db.Exec(query).Error; err != nil {
			return fmt.Errorf("fish_user_links 旧トリガー/関数削除に失敗しました: %w", err)
		}
	}
	return nil
}

func dropLegacyUserPlaceLinksAutomation(db *gorm.DB) error {
	queries := []string{
		`DROP TRIGGER IF EXISTS trg_sync_user_place_links_new_place ON place_cache;`,
		`DROP TRIGGER IF EXISTS trg_sync_user_place_links_new_user ON "user";`,
		`DROP TRIGGER IF EXISTS trg_cleanup_user_place_links_on_place_delete ON place_cache;`,
		`DROP TRIGGER IF EXISTS trg_cleanup_user_place_links_on_user_delete ON "user";`,
		`DROP FUNCTION IF EXISTS sync_user_place_links_new_place();`,
		`DROP FUNCTION IF EXISTS sync_user_place_links_new_user();`,
		`DROP FUNCTION IF EXISTS cleanup_user_place_links_on_place_delete();`,
		`DROP FUNCTION IF EXISTS cleanup_user_place_links_on_user_delete();`,
	}
	for _, query := range queries {
		if err := db.Exec(query).Error; err != nil {
			return fmt.Errorf("user_place_links 旧トリガー/関数削除に失敗しました: %w", err)
		}
	}
	return nil
}

func dropLegacyUserFishLinksAutomation(db *gorm.DB) error {
	queries := []string{
		`DROP TRIGGER IF EXISTS trg_sync_user_fish_links_new_fish ON fish;`,
		`DROP TRIGGER IF EXISTS trg_sync_user_fish_links_new_user ON "user";`,
		`DROP TRIGGER IF EXISTS trg_cleanup_user_fish_links_on_fish_delete ON fish;`,
		`DROP TRIGGER IF EXISTS trg_cleanup_user_fish_links_on_user_delete ON "user";`,
		`DROP FUNCTION IF EXISTS sync_user_fish_links_new_fish();`,
		`DROP FUNCTION IF EXISTS sync_user_fish_links_new_user();`,
		`DROP FUNCTION IF EXISTS cleanup_user_fish_links_on_fish_delete();`,
		`DROP FUNCTION IF EXISTS cleanup_user_fish_links_on_user_delete();`,
	}
	for _, query := range queries {
		if err := db.Exec(query).Error; err != nil {
			return fmt.Errorf("user_fish_links 旧トリガー/関数削除に失敗しました: %w", err)
		}
	}
	return nil
}

func ensureForeignKeyConstraint(db *gorm.DB, query string) error {
	return db.Exec(query).Error
}
