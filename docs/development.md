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

このコマンドで `frontend`（3000）/ `admin`（3001）/ `backend`（8080）を同時起動できます。

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
| 管理画面（admin） | http://localhost:3001 |
| バックエンド API | http://localhost:8080 |
| Hello API | http://localhost:8080/api/hello |

## コーディング規約

- コメントは日本語で記載
- Go: 標準のフォーマッタ (`gofmt`) を使用
- TypeScript: ESLint + Prettier を使用
