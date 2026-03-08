package gorm

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgresDB はDATABASE_URLを利用してPostgreSQL接続を初期化します。
func NewPostgresDB() (*gorm.DB, error) {
	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL が設定されていません")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("PostgreSQL接続に失敗しました: %w", err)
	}

	return db, nil
}
