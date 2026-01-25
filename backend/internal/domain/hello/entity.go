package hello

// helloドメインのエンティティを表す
type Hello struct {
	Message string `json:"message"`
}

// 新しいHelloエンティティを作成する
func NewHello(message string) *Hello {
	return &Hello{
		Message: message,
	}
}
