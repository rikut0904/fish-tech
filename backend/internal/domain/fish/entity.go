package fish

import "time"

// Fish は魚のドメインエンティティです。
type Fish struct {
	ID        string
	NameJa    string
	Name      string
	Explain   string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
