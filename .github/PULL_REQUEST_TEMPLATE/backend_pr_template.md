# Pull Request（バックエンド用）

人に見せるためのPRテンプレートです。必要に応じて内容を編集してください。

## 概要

このPRで対応した内容を簡潔に記載してください。

## 実装内容

### 新規作成

- 対象パッケージ・ファイルを記載

### 修正

- 対象パッケージ・ファイルを記載

### 削除

- 対象パッケージ・ファイルを記載

## 動作確認

- [ ] APIの動作やエラーの確認ポイントを記載
- [ ] 必要に応じてリクエスト例やレスポンス例を記載

## 確認方法

```bash
# バックエンドのディレクトリに移動
cd backend/

# 必要なモジュールのインストール
# 例: go mod tidy

# サーバー起動
# 例: go run cmd/api/main.go

cd backend

docker compose up --build
# 上手くいかないときは docker compose down してから再度 docker compose up --build を試してください。
# docker compose up --build --no-cache を試すのも有効です。

go mod download

go run ./cmd/api/main.go

```

## 関連Issue

- 関連するIssue番号があれば記載してください（例: #123）

## 備考

その他、注意点や補足事項があれば記載してください。
