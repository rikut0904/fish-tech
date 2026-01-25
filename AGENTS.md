# AIエージェント向けガイドライン

このドキュメントは様々なAIエージェント（Claude Codeを除く）がプロジェクトで作業する際のルールを定義します。

## プロジェクト概要

- **プロジェクト名**: fish-tech
- **目的**: 金沢で特に有名な魚種以外の消費者認知度向上
- **主催**: 石川県・金沢市（漁業ハッカソン）

## 技術スタック

### フロントエンド

- Next.js 16.x（App Router）
- TypeScript 5.x
- TailwindCSS 4.x
- Node.js 24.9.0

### バックエンド

- Go 1.25.6
- Echo 4.x
- クリーンアーキテクチャ

## コーディング規約

### 共通

- コメントは**日本語**で記載すること
- ファイル名は小文字のスネークケースを使用

### Go

- パッケージ名は小文字
- エクスポートする関数・型にはコメントを付与

### TypeScript

- ESLint + Prettier を使用
- 型定義を明示的に記載
- `any` 型の使用は避ける

## ディレクトリ構成

### 全体構成

```
fish-tech/
├── frontend/                # フロントエンド
├── backend/                 # バックエンド
├── docs/                    # ドキュメント
├── docker-compose.yml
└── agent.md                 # 本ファイル
```

## フロントエンド新機能追加手順

1. `app/` 配下に新しいページディレクトリを作成（例: `app/fish/page.tsx`）
2. 必要に応じてコンポーネントを `app/components/` に作成
3. API連携が必要な場合は `fetch` または適切なライブラリを使用
4. スタイルは TailwindCSS のユーティリティクラスを使用

### ページ作成例

```tsx
// app/fish/page.tsx
export default function FishPage() {
  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold">魚図鑑</h1>
    </div>
  );
}
```

### API連携例

```tsx
// バックエンドAPIを呼び出す例
const response = await fetch('http://localhost:8080/api/hello');
const data = await response.json();
```

## バックエンド新機能追加手順

1. `internal/domain/` にエンティティを作成
2. `internal/usecase/` にユースケースを作成
3. `internal/interface/handler/` にハンドラーを作成
4. `internal/infrastructure/router/router.go` にルートを追加

## 禁止事項

### 共通

- 環境変数やシークレットをハードコードしない
- `node_modules/`, `.next/`, バイナリファイルを編集しない

### フロントエンド

- `pages/` ディレクトリを使用しない（App Router を使用）
- インラインスタイルの使用を避ける（TailwindCSS を使用）
- `use client` は必要な場合のみ使用

### バックエンド

- クリーンアーキテクチャの依存関係を逆転させない
  - `domain` → 他レイヤーへの依存禁止
  - `usecase` → `interface`, `infrastructure` への依存禁止

## 推奨事項

- 変更前に関連ファイルを読み込んで構造を理解する
- 既存のコードスタイルに合わせる
- 大きな変更は段階的に行う
- テストを追加・更新する

## ドキュメント更新

コード変更に伴い、以下のドキュメントも更新すること:

| 変更内容 | 更新対象 |
|----------|----------|
| 新機能追加 | `docs/features.md` |
| API追加・変更 | `docs/api.md` |
| 開発環境変更 | `docs/development.md` |

## 連絡先

不明点がある場合は、作業を進める前にユーザーに確認すること。
