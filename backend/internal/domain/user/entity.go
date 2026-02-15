package user

import "time"

// User はユーザー情報です。
type User struct {
	UserID      string
	FirebaseUID string
	Name        string
	Mail        string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
