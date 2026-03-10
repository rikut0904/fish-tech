package model

import "time"

// User は user テーブルのGORMモデルです。
type User struct {
	UserID      string     "gorm:\"column:user_id;type:uuid;primaryKey\""
	FirebaseUID string     `gorm:"column:firebase_uid;not null;unique"`
	Name        string     `gorm:"column:name;not null"`
	Mail        string     `gorm:"column:mail;not null"`
	Role        string     `gorm:"column:role;type:varchar(20);not null;default:user"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

// TableName はテーブル名を返します。
func (User) TableName() string { return "user" }
