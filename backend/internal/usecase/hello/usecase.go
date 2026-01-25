package hello

import (
	"fish-tech/internal/domain/hello"
)

// helloユースケースのインターフェースを定義する
type UseCase interface {
	GetHello() *hello.Hello
}

type helloUseCase struct{}

// 新しいhelloユースケースを作成する
func NewHelloUseCase() UseCase {
	return &helloUseCase{}
}

// helloメッセージを返す
func (u *helloUseCase) GetHello() *hello.Hello {
	return hello.NewHello("Hello fish-tech!")
}
