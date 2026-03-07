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
  "imageUrl": "https://photos.google.com/..."
}
```

#### POST /api/admin/fishes

魚を登録する

**リクエスト**

```json
{
  "name": "ヒラメ",
  "category": "白身魚",
  "description": "淡白で上品な味わい",
  "imageUrl": "https://photos.google.com/...",
  "linkUrl": "https://example.com/fish/hirame"
}
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

```json
{
  "fishIdA": "fish-id-a",
  "fishIdB": "fish-id-b",
  "score": 4,
  "memo": "食感のバランスが良い"
}
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
| Content-Type | application/json |

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
