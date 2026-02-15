# アーキテクチャ設計について

バックエンドはクリーンアーキテクチャを採用しています。

## レイヤー構成

| レイヤー | ディレクトリ | 役割 |
|----------|-------------|------|
| ドメイン層 | `internal/domain/` | エンティティ（ビジネスオブジェクト） |
| ユースケース層 | `internal/usecase/` | ビジネスロジック |
| インターフェース層 | `internal/interface/` | ハンドラー（コントローラー） |
| インフラストラクチャ層 | `internal/infrastructure/` | ルーター、DB接続等 |

## 依存関係の方向

```plaintext
infrastructure → interface → usecase → domain
```

外側のレイヤーは内側のレイヤーに依存しますが、内側から外側への依存は禁止とする。

## 新機能の追加方法

1. `internal/domain/` にエンティティを追加
2. `internal/usecase/` にユースケースを追加
3. `internal/interface/handler/` にハンドラーを追加
4. `internal/infrastructure/router/router.go` にルートを追加
