package dto

// UserResponse はユーザーレスポンスDTOです。
type UserResponse struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
	Mail   string `json:"mail"`
	Role   string `json:"role"`
}
