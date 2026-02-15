# コントリビューションガイド

## 開発フロー

1. Issueを作成または既存のIssueを確認
2. `main`ブランチから作業ブランチを作成
3. 変更を実装
4. Pull Requestを作成

## ブランチ命名規則

### フロントエンド

- `feature/[機能名]` : 新機能追加
- `fix/[内容]`      : バグ修正
- `ref/[内容]`   : 軽微な修正やリファクタ
- `hotfix/[内容]`   : 緊急修正
- `chore/[内容]`    : 雑務・依存更新など

#### 命名例（フロントエンド）

| 用途         | ブランチ名例                       |
|--------------|------------------------------------|
| 新機能       | feature/login-form                 |
| バグ修正     | fix/header-logo                    |
| 軽微修正     | ref/typo-footer                    |
| 緊急修正     | hotfix/build-error                 |
| 依存更新     | chore/update-deps                  |

### バックエンド

- `feature/[機能名]`  : 新機能追加
- `fix/[内容]`       : バグ修正
- `ref/[内容]`  : リファクタリング
- `hotfix/[内容]`    : 緊急修正
- `chore/[内容]`     : 雑務・依存更新など

#### 命名例（バックエンド）

| 用途         | ブランチ名例                       |
|--------------|------------------------------------|
| 新機能       | feature/add-fish-api               |
| バグ修正     | fix/invalid-response               |
| リファクタ   | ref/domain-entity                  |
| 緊急修正     | hotfix/server-crash                |
| 依存更新     | chore/update-go-mod                |

## コミットメッセージ

```plaintext
[種別]/[変更内容]

fix/新機能追加
bug/バグ修正
docs/ドキュメント変更
ref/リファクタリング
test/テスト追加・修正
other/ビルド・設定変更
```

## Pull Request

- 関連するIssue番号を記載
- テンプレートに従って記述
- レビュー前にセルフチェック

## コードスタイル

- プロジェクトの既存コードに合わせる
