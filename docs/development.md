# 開発ガイド

## 開発環境

| 項目 | バージョン |
|------|-----------|
| OS | WSL2 (Ubuntu 24.04.3) |
| Docker | 29.1.2 |
| Node.js | 24.9.0 |
| npm | 11.7.0 |
| Go | 1.25.6 |

## プロジェクト構成

```
fish-tech/
├── frontend/                # フロントエンド (Next.js)
│   ├── app/                 # App Router
│   ├── public/              # 静的ファイル
│   ├── Dockerfile
│   └── package.json
├── backend/                 # バックエンド (Go + Echo)
│   ├── cmd/
│   │   └── api/             # エントリーポイント
│   │       └── main.go
│   └── internal/
│       ├── domain/          # ドメイン層（エンティティ）
│       ├── usecase/         # ユースケース層（ビジネスロジック）
│       ├── interface/       # インターフェース層（ハンドラー）
│       └── infrastructure/  # インフラストラクチャ層（ルーター等）
├── docs/                    # ドキュメント
├── docker-compose.yml
└── README.md
```

## 起動方法

### Docker Composeで起動（推奨）

```bash
docker compose up --build
```

### 個別起動

#### フロントエンド

```bash
cd frontend
npm install
npm run dev
```

#### バックエンド

```bash
cd backend
go mod download
go run ./cmd/api/main.go
```

## アクセス先

| サービス | URL |
|----------|-----|
| フロントエンド | http://localhost:3000 |
| バックエンド API | http://localhost:8080 |
| Hello API | http://localhost:8080/api/hello |

## クリーンアーキテクチャについて

バックエンドはクリーンアーキテクチャを採用しています。

### レイヤー構成

| レイヤー | ディレクトリ | 役割 |
|----------|-------------|------|
| ドメイン層 | `internal/domain/` | エンティティ（ビジネスオブジェクト） |
| ユースケース層 | `internal/usecase/` | ビジネスロジック |
| インターフェース層 | `internal/interface/` | ハンドラー（コントローラー） |
| インフラストラクチャ層 | `internal/infrastructure/` | ルーター、DB接続等 |

### 依存関係の方向

```
infrastructure → interface → usecase → domain
```

外側のレイヤーは内側のレイヤーに依存しますが、内側から外側への依存は禁止とする。

## 新機能の追加方法

1. `internal/domain/` にエンティティを追加
2. `internal/usecase/` にユースケースを追加
3. `internal/interface/handler/` にハンドラーを追加
4. `internal/infrastructure/router/router.go` にルートを追加

## コーディング規約

- コメントは日本語で記載
- Go: 標準のフォーマッタ (`gofmt`) を使用
- TypeScript: ESLint + Prettier を使用

