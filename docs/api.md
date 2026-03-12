# API仕様

## 概要

バックエンドAPIはGo + Echoフレームワークで構築されています。

## ベースURL

```
http://localhost:8080/api
```

## エンドポイント一覧

### Hello

| メソッド | パス | 説明 |
|----------|------|------|
| GET | `/api/hello` | 動作確認用エンドポイント |

### Public

| メソッド | パス | 説明 |
|----------|------|------|
| GET | `/api/fishes` | 魚一覧を取得 |
| GET | `/api/pairs` | 魚同士の相性一覧を取得 |
| GET | `/api/recipes` | レシピ一覧を検索 |
| GET | `/api/recipes/seasonal` | 旬の魚レシピを取得 |
| PATCH | `/api/recipes/favorite` | ユーザーのレシピお気に入り状態を更新 |

### Admin

| メソッド | パス | 説明 |
|----------|------|------|
| POST | `/api/admin/fishes/upload-image` | 画像をGoogle Photosへアップロード |
| POST | `/api/admin/fishes` | 魚を登録 |
| PATCH | `/api/admin/fishes/:id/seasons` | 魚の旬月を更新 |
| DELETE | `/api/admin/fishes/:id` | 魚を削除（関連する相性も削除） |
| POST | `/api/admin/pairs` | 魚同士の相性を登録 |
| DELETE | `/api/admin/pairs/:id` | 魚同士の相性を削除 |

`/api/admin/*` は管理画面オリジン（デフォルト: `http://localhost:3001`）からのアクセスのみ許可します。  
許可オリジンは環境変数 `ADMIN_ALLOWED_ORIGINS`（カンマ区切り）で変更できます。

#### GET /api/hello

動作確認用のエンドポイント

**レスポンス**

```json
{
  "message": "Hello fish-tech!"
}
```

**ステータスコード**

| コード | 説明 |
|--------|------|
| 200 | 成功 |

#### GET /api/fishes

魚一覧を取得する

**レスポンス**

```json
{
  "items": [
    {
      "id": "uuid",
      "name": "ヒラメ",
      "category": "白身魚",
      "description": "淡白で上品な味わい",
      "imageUrl": "https://photos.google.com/...",
      "linkUrl": "https://example.com/fish/hirame"
    }
  ]
}
```

#### GET /api/recipes

レシピ一覧を検索する

**クエリパラメータ**

| パラメータ | 型 | 必須 | 説明 |
|----------|---|---|---|
| `fishName` | string | 任意 | 魚名で検索 |
| `keyword` | string | 任意 | レシピ名や補足文で検索 |
| `userId` | string | 任意 | お気に入り状態判定に利用 |
| `favorite` | bool | 任意 | `true/false` でお気に入り状態を絞り込み |
| `count` | number | 任意 | 1ページ件数。未指定時 `10`、最大 `100` |
| `page` | number | 任意 | ページ番号。未指定時 `1` |

補足:

- `favorite` を使う場合は `userId` が必須
- 未対応クエリは `400` を返す

**レスポンス**

```json
{
  "page": 1,
  "perPage": 10,
  "count": 2,
  "total": 12,
  "items": [
    {
      "id": "12345",
      "title": "いわしの蒲焼き",
      "imageUrl": "https://recipe.r10s.jp/...",
      "recipeUrl": "https://recipe.rakuten.co.jp/recipe/12345/",
      "cookingTime": "約15分",
      "cost": "300円前後",
      "score": 5,
      "isLikes": false,
      "explain": "楽天レシピカテゴリ「いわし」の人気レシピです。"
    }
  ]
}
```

#### GET /api/recipes/seasonal

旬の魚レシピを取得する

**クエリパラメータ**

| パラメータ | 型 | 必須 | 説明 |
|----------|---|---|---|
| `fishName` | string | 任意 | 対象魚名。指定時に今月の旬でなければ `400` |
| `userId` | string | 任意 | お気に入り状態判定に利用 |
| `favorite` | bool | 任意 | `true/false` でお気に入り状態を絞り込み |
| `count` | number | 任意 | 1ページ件数。未指定時 `10`、最大 `100` |
| `page` | number | 任意 | ページ番号。未指定時 `1` |

補足:

- `favorite` を使う場合は `userId` が必須
- 未対応クエリは `400` を返す
- `fishName` を指定して今月の旬でない場合は `指定した魚は今月の旬ではありません` を返す

**レスポンス**

```json
{
  "month": 3,
  "selectedFishId": "uuid",
  "page": 1,
  "perPage": 10,
  "count": 1,
  "total": 4,
  "fishes": [
    {
      "id": "uuid",
      "name": "いわし",
      "imageUrl": "https://example.com/fish.jpg"
    }
  ],
  "recipes": [
    {
      "id": "12345",
      "title": "いわしの蒲焼き",
      "imageUrl": "https://recipe.r10s.jp/...",
      "recipeUrl": "https://recipe.rakuten.co.jp/recipe/12345/",
      "cookingTime": "約15分",
      "cost": "300円前後",
      "score": 5,
      "isLikes": false,
      "explain": "楽天レシピカテゴリ「いわし」の人気レシピです。"
    }
  ]
}
```

#### PATCH /api/recipes/favorite

ユーザーとレシピの関連を必要時に作成しつつ、お気に入り状態を更新する

**リクエスト**

- クエリパラメータまたは JSON body を指定
- 必須: `userId`, `recipeId`
- 任意: `isLikes` または `favorite`
- `user_recipe_links` に対象行がない場合は自動作成され、`isLikes` の初期値は `false`

```text
PATCH /api/recipes/favorite?userId=user-uuid&recipeId=12345&isLikes=true
```

```json
{
  "userId": "user-uuid",
  "recipeId": "12345",
  "isLikes": true
}
```

**レスポンス**

```json
{
  "userId": "user-uuid",
  "recipeId": "12345",
  "isLikes": true
}
```

#### POST /api/admin/fishes/upload-image

画像ファイルをGoogle Photosへアップロードする

**リクエスト**

- `multipart/form-data`
- フィールド名: `file`

**レスポンス**

```json
{
  "imageUrl": "https://lh3.googleusercontent.com/...",
  "imageMediaId": "ANX7n..."
}
```

#### POST /api/admin/fishes

魚を登録する

**リクエスト**

- クエリパラメータを指定
- 必須: `name`, `category`
- 任意: `description`, `imageUrl`, `imageMediaId`, `linkUrl`

```text
POST /api/admin/fishes?name=ヒラメ&category=白身魚&description=淡白で上品な味わい&imageUrl=https%3A%2F%2Flh3.googleusercontent.com%2F...&imageMediaId=ANX7n...&linkUrl=https%3A%2F%2Fexample.com%2Ffish%2Fhirame
```

**レスポンス**

```json
{
  "id": "uuid",
  "name": "ヒラメ",
  "category": "白身魚",
  "description": "淡白で上品な味わい",
  "imageUrl": "https://photos.google.com/...",
  "linkUrl": "https://example.com/fish/hirame",
  "months": []
}
```

#### PATCH /api/admin/fishes/:id/seasons

魚の旬月をまとめて更新する

**リクエスト**

- クエリ `months=3,4,5` または JSON body の `months` を指定
- 月は `1-12`
- 既存の旬月は全削除して置き換える

```text
PATCH /api/admin/fishes/<FISH_ID>/seasons?months=3,4,5
```

```json
{
  "months": [3, 4, 5]
}
```

**レスポンス**

- `204 No Content`

#### DELETE /api/admin/fishes/:id

魚を削除する

**レスポンス**

- `204 No Content`

#### GET /api/pairs

魚同士の相性一覧を取得する

**レスポンス**

```json
{
  "items": [
    {
      "id": "fish-id-a:fish-id-b",
      "fishIdA": "fish-id-a",
      "fishIdB": "fish-id-b",
      "score": 4,
      "memo": "食感のバランスが良い"
    }
  ]
}
```

#### POST /api/admin/pairs

魚同士の相性を登録する

**リクエスト**

- クエリパラメータを指定
- 必須: `fishIdA`, `fishIdB`, `score`
- 任意: `memo`

```text
POST /api/admin/pairs?fishIdA=fish-id-a&fishIdB=fish-id-b&score=4&memo=食感のバランスが良い
```

**レスポンス**

```json
{
  "id": "fish-id-a:fish-id-b",
  "fishIdA": "fish-id-a",
  "fishIdB": "fish-id-b",
  "score": 4,
  "memo": "食感のバランスが良い"
}
```

#### DELETE /api/admin/pairs/:id

魚同士の相性を削除する

**レスポンス**

- `204 No Content`

---

## 共通仕様

### リクエストヘッダー

| ヘッダー | 値 |
|----------|-----|
| Content-Type | `multipart/form-data`（`POST /api/admin/fishes/upload-image` のみ） |

### エラーレスポンス

```json
{
  "error": "エラーメッセージ"
}
```

### ステータスコード

| コード | 説明 |
|--------|------|
| 200 | 成功 |
| 400 | リクエスト不正 |
| 404 | リソースが見つからない |
| 500 | サーバーエラー |
