package user

import "time"

const (
	// RoleAdmin は管理者権限を表します。
	RoleAdmin = "admin"
	// RoleUser は一般ユーザー権限を表します。
	RoleUser = "user"
)

// User はユーザー情報です。
type User struct {
	UserID      string
	FirebaseUID string
	Name        string
	Mail        string
	Role        string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
