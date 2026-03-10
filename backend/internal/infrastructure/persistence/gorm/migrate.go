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

func migrateFishUserLinks(db *gorm.DB) error {
	if !db.Migrator().HasTable(&model.FishUserLinks{}) {
		return nil
	}

	if err := db.Exec(`
		INSERT INTO fish_user_links (fish_id, user_id, is_likes, created_at)
		SELECT f.id, u.user_id, FALSE, NOW()
		FROM fish AS f
		CROSS JOIN "user" AS u
		ON CONFLICT (fish_id, user_id) DO NOTHING
	`).Error; err != nil {
		return fmt.Errorf("fish_user_links の初期バックフィルに失敗しました: %w", err)
	}

	if err := db.Exec(`
		CREATE OR REPLACE FUNCTION sync_fish_user_links_new_fish()
		RETURNS trigger AS $$
		BEGIN
			INSERT INTO fish_user_links (fish_id, user_id, is_likes, created_at)
			SELECT NEW.id, u.user_id, FALSE, NOW()
			FROM "user" AS u
			ON CONFLICT (fish_id, user_id) DO NOTHING;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`).Error; err != nil {
		return fmt.Errorf("fish_user_links 新規魚トリガー関数の作成に失敗しました: %w", err)
	}

	if err := db.Exec(`DROP TRIGGER IF EXISTS trg_sync_fish_user_links_new_fish ON fish;`).Error; err != nil {
		return fmt.Errorf("fish_user_links 新規魚トリガーの削除に失敗しました: %w", err)
	}
	if err := db.Exec(`
		CREATE TRIGGER trg_sync_fish_user_links_new_fish
		AFTER INSERT ON fish
		FOR EACH ROW
		EXECUTE FUNCTION sync_fish_user_links_new_fish();
	`).Error; err != nil {
		return fmt.Errorf("fish_user_links 新規魚トリガーの作成に失敗しました: %w", err)
	}

	if err := db.Exec(`
		CREATE OR REPLACE FUNCTION sync_fish_user_links_new_user()
		RETURNS trigger AS $$
		BEGIN
			INSERT INTO fish_user_links (fish_id, user_id, is_likes, created_at)
			SELECT f.id, NEW.user_id, FALSE, NOW()
			FROM fish AS f
			ON CONFLICT (fish_id, user_id) DO NOTHING;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`).Error; err != nil {
		return fmt.Errorf("fish_user_links 新規ユーザートリガー関数の作成に失敗しました: %w", err)
	}

	if err := db.Exec(`DROP TRIGGER IF EXISTS trg_sync_fish_user_links_new_user ON "user";`).Error; err != nil {
		return fmt.Errorf("fish_user_links 新規ユーザートリガーの削除に失敗しました: %w", err)
	}
	if err := db.Exec(`
		CREATE TRIGGER trg_sync_fish_user_links_new_user
		AFTER INSERT ON "user"
		FOR EACH ROW
		EXECUTE FUNCTION sync_fish_user_links_new_user();
	`).Error; err != nil {
		return fmt.Errorf("fish_user_links 新規ユーザートリガーの作成に失敗しました: %w", err)
	}

	if err := db.Exec(`
		CREATE OR REPLACE FUNCTION cleanup_fish_user_links_on_fish_delete()
		RETURNS trigger AS $$
		BEGIN
			DELETE FROM fish_user_links WHERE fish_id = OLD.id;
			RETURN OLD;
		END;
		$$ LANGUAGE plpgsql;
	`).Error; err != nil {
		return fmt.Errorf("fish_user_links 魚削除クリーンアップ関数の作成に失敗しました: %w", err)
	}
	if err := db.Exec(`DROP TRIGGER IF EXISTS trg_cleanup_fish_user_links_on_fish_delete ON fish;`).Error; err != nil {
		return fmt.Errorf("fish_user_links 魚削除クリーンアップトリガーの削除に失敗しました: %w", err)
	}
	if err := db.Exec(`
		CREATE TRIGGER trg_cleanup_fish_user_links_on_fish_delete
		AFTER DELETE ON fish
		FOR EACH ROW
		EXECUTE FUNCTION cleanup_fish_user_links_on_fish_delete();
	`).Error; err != nil {
		return fmt.Errorf("fish_user_links 魚削除クリーンアップトリガーの作成に失敗しました: %w", err)
	}

	if err := db.Exec(`
		CREATE OR REPLACE FUNCTION cleanup_fish_user_links_on_user_delete()
		RETURNS trigger AS $$
		BEGIN
			DELETE FROM fish_user_links WHERE user_id = OLD.user_id;
			RETURN OLD;
		END;
		$$ LANGUAGE plpgsql;
	`).Error; err != nil {
		return fmt.Errorf("fish_user_links ユーザー削除クリーンアップ関数の作成に失敗しました: %w", err)
	}
	if err := db.Exec(`DROP TRIGGER IF EXISTS trg_cleanup_fish_user_links_on_user_delete ON "user";`).Error; err != nil {
		return fmt.Errorf("fish_user_links ユーザー削除クリーンアップトリガーの削除に失敗しました: %w", err)
	}
	if err := db.Exec(`
		CREATE TRIGGER trg_cleanup_fish_user_links_on_user_delete
		AFTER DELETE ON "user"
		FOR EACH ROW
		EXECUTE FUNCTION cleanup_fish_user_links_on_user_delete();
	`).Error; err != nil {
		return fmt.Errorf("fish_user_links ユーザー削除クリーンアップトリガーの作成に失敗しました: %w", err)
	}

	return nil
}

func migrateUserPlaceLinks(db *gorm.DB) error {
	if !db.Migrator().HasTable(&model.UserPlaceLinks{}) {
		return nil
	}

	if err := db.Exec(`
		INSERT INTO user_place_links (user_id, place_id, is_likes, created_at)
		SELECT u.user_id, p.id, FALSE, NOW()
		FROM "user" AS u
		CROSS JOIN place_cache AS p
		ON CONFLICT (user_id, place_id) DO NOTHING
	`).Error; err != nil {
		return fmt.Errorf("user_place_links の初期バックフィルに失敗しました: %w", err)
	}

	if err := db.Exec(`
		CREATE OR REPLACE FUNCTION sync_user_place_links_new_place()
		RETURNS trigger AS $$
		BEGIN
			INSERT INTO user_place_links (user_id, place_id, is_likes, created_at)
			SELECT u.user_id, NEW.id, FALSE, NOW()
			FROM "user" AS u
			ON CONFLICT (user_id, place_id) DO NOTHING;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`).Error; err != nil {
		return fmt.Errorf("user_place_links 新規店舗トリガー関数の作成に失敗しました: %w", err)
	}

	if err := db.Exec(`DROP TRIGGER IF EXISTS trg_sync_user_place_links_new_place ON place_cache;`).Error; err != nil {
		return fmt.Errorf("user_place_links 新規店舗トリガーの削除に失敗しました: %w", err)
	}
	if err := db.Exec(`
		CREATE TRIGGER trg_sync_user_place_links_new_place
		AFTER INSERT ON place_cache
		FOR EACH ROW
		EXECUTE FUNCTION sync_user_place_links_new_place();
	`).Error; err != nil {
		return fmt.Errorf("user_place_links 新規店舗トリガーの作成に失敗しました: %w", err)
	}

	if err := db.Exec(`
		CREATE OR REPLACE FUNCTION sync_user_place_links_new_user()
		RETURNS trigger AS $$
		BEGIN
			INSERT INTO user_place_links (user_id, place_id, is_likes, created_at)
			SELECT NEW.user_id, p.id, FALSE, NOW()
			FROM place_cache AS p
			ON CONFLICT (user_id, place_id) DO NOTHING;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`).Error; err != nil {
		return fmt.Errorf("user_place_links 新規ユーザートリガー関数の作成に失敗しました: %w", err)
	}

	if err := db.Exec(`DROP TRIGGER IF EXISTS trg_sync_user_place_links_new_user ON "user";`).Error; err != nil {
		return fmt.Errorf("user_place_links 新規ユーザートリガーの削除に失敗しました: %w", err)
	}
	if err := db.Exec(`
		CREATE TRIGGER trg_sync_user_place_links_new_user
		AFTER INSERT ON "user"
		FOR EACH ROW
		EXECUTE FUNCTION sync_user_place_links_new_user();
	`).Error; err != nil {
		return fmt.Errorf("user_place_links 新規ユーザートリガーの作成に失敗しました: %w", err)
	}

	if err := db.Exec(`
		CREATE OR REPLACE FUNCTION cleanup_user_place_links_on_place_delete()
		RETURNS trigger AS $$
		BEGIN
			DELETE FROM user_place_links WHERE place_id = OLD.id;
			RETURN OLD;
		END;
		$$ LANGUAGE plpgsql;
	`).Error; err != nil {
		return fmt.Errorf("user_place_links 店舗削除クリーンアップ関数の作成に失敗しました: %w", err)
	}
	if err := db.Exec(`DROP TRIGGER IF EXISTS trg_cleanup_user_place_links_on_place_delete ON place_cache;`).Error; err != nil {
		return fmt.Errorf("user_place_links 店舗削除クリーンアップトリガーの削除に失敗しました: %w", err)
	}
	if err := db.Exec(`
		CREATE TRIGGER trg_cleanup_user_place_links_on_place_delete
		AFTER DELETE ON place_cache
		FOR EACH ROW
		EXECUTE FUNCTION cleanup_user_place_links_on_place_delete();
	`).Error; err != nil {
		return fmt.Errorf("user_place_links 店舗削除クリーンアップトリガーの作成に失敗しました: %w", err)
	}

	if err := db.Exec(`
		CREATE OR REPLACE FUNCTION cleanup_user_place_links_on_user_delete()
		RETURNS trigger AS $$
		BEGIN
			DELETE FROM user_place_links WHERE user_id = OLD.user_id;
			RETURN OLD;
		END;
		$$ LANGUAGE plpgsql;
	`).Error; err != nil {
		return fmt.Errorf("user_place_links ユーザー削除クリーンアップ関数の作成に失敗しました: %w", err)
	}
	if err := db.Exec(`DROP TRIGGER IF EXISTS trg_cleanup_user_place_links_on_user_delete ON "user";`).Error; err != nil {
		return fmt.Errorf("user_place_links ユーザー削除クリーンアップトリガーの削除に失敗しました: %w", err)
	}
	if err := db.Exec(`
		CREATE TRIGGER trg_cleanup_user_place_links_on_user_delete
		AFTER DELETE ON "user"
		FOR EACH ROW
		EXECUTE FUNCTION cleanup_user_place_links_on_user_delete();
	`).Error; err != nil {
		return fmt.Errorf("user_place_links ユーザー削除クリーンアップトリガーの作成に失敗しました: %w", err)
	}

	return nil
}

func migrateUserFishLinks(db *gorm.DB) error {
	if !db.Migrator().HasTable(&model.UserFishLinks{}) {
		return nil
	}

	if err := db.Exec(`
		INSERT INTO user_fish_links (user_id, fish_id, "like", created_at)
		SELECT u.user_id, f.id, FALSE, NOW()
		FROM "user" AS u
		CROSS JOIN fish AS f
		ON CONFLICT (user_id, fish_id) DO NOTHING
	`).Error; err != nil {
		return fmt.Errorf("user_fish_links の初期バックフィルに失敗しました: %w", err)
	}

	if err := db.Exec(`
		CREATE OR REPLACE FUNCTION sync_user_fish_links_new_fish()
		RETURNS trigger AS $$
		BEGIN
			INSERT INTO user_fish_links (user_id, fish_id, "like", created_at)
			SELECT u.user_id, NEW.id, FALSE, NOW()
			FROM "user" AS u
			ON CONFLICT (user_id, fish_id) DO NOTHING;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`).Error; err != nil {
		return fmt.Errorf("user_fish_links 新規魚トリガー関数の作成に失敗しました: %w", err)
	}
	if err := db.Exec(`DROP TRIGGER IF EXISTS trg_sync_user_fish_links_new_fish ON fish;`).Error; err != nil {
		return fmt.Errorf("user_fish_links 新規魚トリガーの削除に失敗しました: %w", err)
	}
	if err := db.Exec(`
		CREATE TRIGGER trg_sync_user_fish_links_new_fish
		AFTER INSERT ON fish
		FOR EACH ROW
		EXECUTE FUNCTION sync_user_fish_links_new_fish();
	`).Error; err != nil {
		return fmt.Errorf("user_fish_links 新規魚トリガーの作成に失敗しました: %w", err)
	}

	if err := db.Exec(`
		CREATE OR REPLACE FUNCTION sync_user_fish_links_new_user()
		RETURNS trigger AS $$
		BEGIN
			INSERT INTO user_fish_links (user_id, fish_id, "like", created_at)
			SELECT NEW.user_id, f.id, FALSE, NOW()
			FROM fish AS f
			ON CONFLICT (user_id, fish_id) DO NOTHING;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`).Error; err != nil {
		return fmt.Errorf("user_fish_links 新規ユーザートリガー関数の作成に失敗しました: %w", err)
	}
	if err := db.Exec(`DROP TRIGGER IF EXISTS trg_sync_user_fish_links_new_user ON "user";`).Error; err != nil {
		return fmt.Errorf("user_fish_links 新規ユーザートリガーの削除に失敗しました: %w", err)
	}
	if err := db.Exec(`
		CREATE TRIGGER trg_sync_user_fish_links_new_user
		AFTER INSERT ON "user"
		FOR EACH ROW
		EXECUTE FUNCTION sync_user_fish_links_new_user();
	`).Error; err != nil {
		return fmt.Errorf("user_fish_links 新規ユーザートリガーの作成に失敗しました: %w", err)
	}

	if err := db.Exec(`
		CREATE OR REPLACE FUNCTION cleanup_user_fish_links_on_fish_delete()
		RETURNS trigger AS $$
		BEGIN
			DELETE FROM user_fish_links WHERE fish_id = OLD.id;
			RETURN OLD;
		END;
		$$ LANGUAGE plpgsql;
	`).Error; err != nil {
		return fmt.Errorf("user_fish_links 魚削除クリーンアップ関数の作成に失敗しました: %w", err)
	}
	if err := db.Exec(`DROP TRIGGER IF EXISTS trg_cleanup_user_fish_links_on_fish_delete ON fish;`).Error; err != nil {
		return fmt.Errorf("user_fish_links 魚削除クリーンアップトリガーの削除に失敗しました: %w", err)
	}
	if err := db.Exec(`
		CREATE TRIGGER trg_cleanup_user_fish_links_on_fish_delete
		AFTER DELETE ON fish
		FOR EACH ROW
		EXECUTE FUNCTION cleanup_user_fish_links_on_fish_delete();
	`).Error; err != nil {
		return fmt.Errorf("user_fish_links 魚削除クリーンアップトリガーの作成に失敗しました: %w", err)
	}

	if err := db.Exec(`
		CREATE OR REPLACE FUNCTION cleanup_user_fish_links_on_user_delete()
		RETURNS trigger AS $$
		BEGIN
			DELETE FROM user_fish_links WHERE user_id = OLD.user_id;
			RETURN OLD;
		END;
		$$ LANGUAGE plpgsql;
	`).Error; err != nil {
		return fmt.Errorf("user_fish_links ユーザー削除クリーンアップ関数の作成に失敗しました: %w", err)
	}
	if err := db.Exec(`DROP TRIGGER IF EXISTS trg_cleanup_user_fish_links_on_user_delete ON "user";`).Error; err != nil {
		return fmt.Errorf("user_fish_links ユーザー削除クリーンアップトリガーの削除に失敗しました: %w", err)
	}
	if err := db.Exec(`
		CREATE TRIGGER trg_cleanup_user_fish_links_on_user_delete
		AFTER DELETE ON "user"
		FOR EACH ROW
		EXECUTE FUNCTION cleanup_user_fish_links_on_user_delete();
	`).Error; err != nil {
		return fmt.Errorf("user_fish_links ユーザー削除クリーンアップトリガーの作成に失敗しました: %w", err)
	}

	return nil
}
