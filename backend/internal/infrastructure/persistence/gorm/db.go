package gorm

import (
	"fmt"
	"os"
	"strings"
	"time"

	"fish-tech/internal/shared/timeutil"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgresDB はDATABASE_URLを利用してPostgreSQL接続を初期化します。
func NewPostgresDB() (*gorm.DB, error) {
	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL が設定されていません")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return timeutil.NowJST()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("PostgreSQL接続に失敗しました: %w", err)
	}
	if err := db.Exec("SET TIME ZONE 'Asia/Tokyo'").Error; err != nil {
		fmt.Fprintf(os.Stderr, "warn: DBセッションのタイムゾーン設定に失敗しました: %v\n", err)
	}

	return db, nil
}
