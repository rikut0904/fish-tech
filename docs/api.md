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
| GET | `/api/places/recommendations` | おすすめ店舗を取得（HotPepper + DBキャッシュ） |
| PATCH | `/api/places/favorite` | 店舗のお気に入り状態を更新 |

### Admin

| メソッド | パス | 説明 |
|----------|------|------|
| POST | `/api/admin/fishes/upload-image` | 画像をGoogle Photosへアップロード |
| POST | `/api/admin/fishes` | 魚を登録 |
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
  "linkUrl": "https://example.com/fish/hirame"
}
```

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

#### GET /api/places/recommendations

おすすめ店舗を取得する。  
`place_cache` をキャッシュとして利用し、`fetched_at` が1か月以内のデータはDBから返す。  
キャッシュが存在しない、または1か月より古い場合のみ HotPepper API を再取得して更新する。

**クエリパラメータ**

| パラメータ | 型 | 必須 | 説明 |
|----------|---|---|---|
| `fishName` | string | 任意 | 魚名（`keyword`未指定時は `魚名 + 魚料理` で検索） |
| `keyword` | string | 任意 | 検索キーワード（指定時も `魚` を含めて検索） |
| `cityCode` | string | 任意 | 石川県の HotPepper `small_area` コード |
| `userId` | string | 任意 | ユーザーID（`favorite=true` の場合は必須） |
| `favorite` | boolean | 任意 | `true` の場合はユーザーのお気に入り店舗のみ取得 |
| `count` | number | 任意 | 取得件数（1〜100、未指定は10） |
| `page` | number | 任意 | ページ番号（1以上、未指定は1） |

※ エリアは石川県で固定されます（HotPepper `large_area=Z063` を使用）。  
※ `cityCode` が石川県の `small_area` に存在しない場合は `404` を返します。
※ `favorite=true` で `userId` 未指定の場合は `400` を返します。  
※ キャッシュは `fetched_at` が30日以内のデータを使用し、該当キャッシュが無い（または古い）場合は HotPepper API を再取得します。
※ `cityCode` 判定用の `small_area` 一覧はDBにキャッシュし、1年経過時のみ再取得します。
※ `user_*_links` / `fish_user_links` / `user_place_links` は全組み合わせ事前作成を行わず、操作時にオンデマンド作成します。

**レスポンス**

```json
{
  "items": [
    {
      "name": "魚料理の店 さかな亭",
      "address": "石川県金沢市...",
      "lat": "36.561325",
      "lng": "136.656205",
      "coupon": "https://www.hotpepper.jp/strJ001234567/coupon/",
      "genre": "居酒屋",
      "card": "利用可",
      "logo": "https://imgfp.hotp.jp/IMGH/00/00/P000000000/P000000000_238.jpg"
    }
  ],
  "count": 1,
  "page": 1,
  "perPage": 10
}
```

---

#### PATCH /api/places/favorite

ユーザーの店舗お気に入り状態を1件更新する。  
`userId + placeId` の複合キーで更新し、存在しない場合は作成する。

**クエリパラメータ（対応）**

| パラメータ | 型 | 必須 | 説明 |
|----------|---|---|---|
| `userId` | string | 必須 | ユーザーID |
| `placeId` | string | 必須 | 店舗ID |
| `favorite` | boolean | 必須 | お気に入り状態 |

例: `PATCH /api/places/favorite?userId=...&placeId=J001142822&favorite=true`

**リクエストボディ（後方互換で対応）**

```json
{
  "userId": "uuid",
  "placeId": "J001142822",
  "favorite": true
}
```

**レスポンス**

- `204 No Content`

**エラー**

- `400`: `userId` または `placeId` が空
- `404`: 指定 `userId` / `placeId` が存在しない

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
