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

### バックエンド環境変数

```env
DATABASE_URL=postgresql://user:password@host:port/database
HOTPEPPER_API_KEY=xxxxxxxxxxxxxxxx
AUTO_MIGRATE=false
```

- `AUTO_MIGRATE` は起動時に DB マイグレーションを実行するかどうかの設定です。
- 既定値は `false` です。`true` を設定したときだけ、バックエンド起動時に `AutoMigrateAll` を実行します。

Google Photosへの画像アップロード機能を利用する場合は、下記の手順で取得した値を `backend/.env` に設定してください。

### Google Photos連携の設定手順

管理画面で魚画像をGoogle Photosへアップロードするには、Google Cloud で OAuth 2.0 の認証情報と Refresh Token が必要です。

**なぜ「デスクトップアプリ」？**  
fish-tech は Web アプリですが、Google Photos API は **バックエンド（Go）からのみ** 呼ばれます。管理画面のブラウザで「Googleでログイン」するのではなく、**一度だけ** 管理者が認可して取得した Refresh Token を `backend/.env` に置き、サーバーがそのトークンでアップロードする構成です。  
Refresh Token を取得する手順だけ「認可URLを開く → 表示されたコードをコピー → curl でトークン交換」としたいため、リダイレクトURL不要の **デスクトップアプリ** タイプを使うと設定が簡単です。**ウェブアプリケーション** タイプでも可能で、その場合はリダイレクトURI（例: `http://localhost:8088/callback`）を登録し、ローカルで一時サーバーを立てて `code` を受け取る形になります。

#### 1. Google Cloud でプロジェクトとAPIを有効化

1. [Google Cloud Console](https://console.cloud.google.com/) にアクセス
2. プロジェクトを選択（または新規作成）
3. **APIとサービス** → **ライブラリ** で「**Google Photos Library API**」を検索し、**有効にする**

#### 2. OAuth 同意画面の設定

1. **APIとサービス** → **OAuth 同意画面**
2. ユーザータイプは「**外部**」で作成（テスト時は自分のGoogleアカウントのみ追加可能）
3. **スコープの追加手順**
   - 同意画面の編集画面で「**スコープ**」のところまで進む
   - **「スコープを追加または削除」**（または "ADD OR REMOVE SCOPES"）をクリック
   - 開いたパネルで検索欄に「**Photos**」や「**photoslibrary**」と入力
   - 一覧から次をチェックして **「更新」** する  
     - `.../auth/photoslibrary.appendonly` … 写真の追加  
     - `.../auth/photoslibrary.readonly.appcreateddata` … アプリで追加した写真の閲覧（表示URL取得に必須）
   - 画面の「保存して次へ」でスコープを確定
4. **テストユーザー（重要）**  
   アプリが「テスト」状態のあいだは、**ここに追加したメールアドレスのアカウントだけ**が認可できます。認可で使う **自分の Google アカウントのメール**（例: `iwashi.kurukuru@gmail.com`）を必ず追加してください。追加していないと認可時に **403 access_denied** になります。  
   - OAuth 同意画面の編集で「**テストユーザー**」のステップへ進む → **「ADD USERS」／「ユーザーを追加」** → メールアドレスを入力して保存

#### 3. 認証情報（Client ID / Client Secret）の取得

1. **APIとサービス** → **認証情報** → **認証情報を作成** → **OAuth クライアント ID**
2. アプリケーションの種類は次のどちらかで作成します。
   - **デスクトップアプリ** … リダイレクトURL不要で、後述の「認可コードをコピーして curl」がそのまま使える（推奨）
   - **ウェブアプリケーション** … 使用する場合は「認可済みのリダイレクト URI」に `http://localhost:8088/callback` などを追加し、方法Bで code を受け取る
3. 名前を入力して作成すると、**クライアント ID** と **クライアント シークレット** が表示される  
   → これらを `GOOGLE_PHOTOS_CLIENT_ID` と `GOOGLE_PHOTOS_CLIENT_SECRET` に設定

#### 4. Refresh Token の取得

Refresh Token は「一度だけ」OAuth の認可フローをブラウザで行う必要があります。

**方法A: 認可コードをコピーして curl（デスクトップアプリ用・推奨）**

※ 認証情報を「デスクトップアプリ」で作成した場合に使えます。認可URLを開いたときに **401 invalid_client** が出る場合は、下記「トラブルシューティング: 401 invalid_client」を参照してください。

1. 以下のURLをブラウザで開く（`YOUR_CLIENT_ID` を実際の Client ID に置き換え）:

   ```
   https://accounts.google.com/o/oauth2/v2/auth?client_id=YOUR_CLIENT_ID&redirect_uri=urn:ietf:wg:oauth:2.0:oob&response_type=code&scope=https://www.googleapis.com/auth/photoslibrary.appendonly%20https://www.googleapis.com/auth/photoslibrary.readonly.appcreateddata
   ```

2. 表示された認可画面で「許可」をクリック
3. 画面に表示される **認可コード**（長い文字列）をコピー
4. ターミナルで以下を実行（`YOUR_CLIENT_ID` / `YOUR_CLIENT_SECRET` / `認可コード` を置き換え）:

   ```bash
   curl -X POST https://oauth2.googleapis.com/token \
     -d "client_id=YOUR_CLIENT_ID" \
     -d "client_secret=YOUR_CLIENT_SECRET" \
     -d "code=認可コード" \
     -d "grant_type=authorization_code" \
     -d "redirect_uri=urn:ietf:wg:oauth:2.0:oob"
   ```

5. レスポンスの JSON に含まれる `refresh_token` の値を `GOOGLE_PHOTOS_REFRESH_TOKEN` に設定

**方法B: ウェブアプリケーション + ローカルリダイレクトで取得**

- 認証情報を「**ウェブアプリケーション**」で作成した場合は、認可済みのリダイレクト URI に `http://localhost:8088/callback` などを登録し、そのパスで `code` クエリを受け取る簡易サーバーをローカルで起動します。ブラウザがリダイレクトされた先のURLから `code` を取得し、方法Aと同様に `redirect_uri=http://localhost:8088/callback` を指定して `curl` でトークン交換します。

**デスクトップアプリでできないときの対応（localhost リダイレクトで取得する）**

方法A（OOB・認可コードをコピー）で 401 などが出てうまくいかない場合は、**同じプロジェクトで「ウェブアプリケーション」の OAuth クライアントを1つ追加**し、localhost へリダイレクトするやり方で Refresh Token を取得します。取得後の `GOOGLE_PHOTOS_REFRESH_TOKEN` はそのまま使えるので、バックエンドの動きは変わりません。

1. **Google Cloud Console** → **APIとサービス** → **認証情報** → **認証情報を作成** → **OAuth クライアント ID**
2. アプリケーションの種類で **「ウェブアプリケーション」** を選択。名前は「Google Photos 用（ローカル取得）」など任意でよい。
3. **認可済みのリダイレクト URI** で **「URI を追加」** をクリックし、次を1件追加する:  
   `http://localhost:8088/callback`  
   保存して、表示された **クライアント ID** と **クライアント シークレット** を控える（ここでは「ウェブ」用の ID/Secret を使う）。
4. **ローカルで簡易サーバーを起動する**（別ターミナルで実行。ブラウザが `http://localhost:8088/callback?code=...` に飛ばされたときに表示して終了する用）:
   ```bash
   # 例: Node が入っている場合（1回だけ実行。Ctrl+C で止める）
   node -e "const http=require('http'); const s=http.createServer((q,r)=>{const u=new URL(q.url,'http://x'); const c=u.searchParams.get('code'); r.setHeader('Content-Type','text/html;charset=utf-8'); r.end(c ? '<p>認可コード（コピー）:</p><pre>'+c+'</pre><p>このウィンドウを閉じて、コピーしたコードで curl を実行してください。</p>' : 'code not found'); s.close();}); s.listen(8088, ()=>{ console.log('http://localhost:8088 で待機中。認可URLをブラウザで開いてください。');});"
   ```
   - 別の例（Python 3 がある場合）:
   ```bash
   python3 -c "
   from http.server import HTTPServer, BaseHTTPRequestHandler
   from urllib.parse import urlparse, parse_qs
   class H(BaseHTTPRequestHandler):
       def do_GET(self):
           q = parse_qs(urlparse(self.path).query)
           c = q.get('code', [''])[0]
           self.send_response(200); self.send_header('Content-type','text/html; charset=utf-8'); self.end_headers()
           self.wfile.write(('<p>認可コード（コピー）:</p><pre>%s</pre><p>このウィンドウを閉じて、コピーしたコードで curl を実行してください。</p>' % c).encode())
           raise SystemExit
   print('http://localhost:8088 で待機中。認可URLをブラウザで開いてください。')
   HTTPServer(('',8088), H).handle_request()
   "
   ```
5. **認可URLをブラウザで開く**（`YOUR_CLIENT_ID` を手順3で控えた**ウェブアプリケーション用**の Client ID に、`redirect_uri` はそのままに置き換え）:
   ```
   https://accounts.google.com/o/oauth2/v2/auth?client_id=YOUR_CLIENT_ID&redirect_uri=http://localhost:8088/callback&response_type=code&scope=https://www.googleapis.com/auth/photoslibrary.appendonly%20https://www.googleapis.com/auth/photoslibrary.readonly.appcreateddata
   ```
6. 認可して「許可」をクリックすると、ブラウザが `http://localhost:8088/callback?code=xxxxx` にリダイレクトされ、簡易サーバーが **認可コード** を表示する。そのコードをコピーする。
7. **トークン交換**（`YOUR_CLIENT_ID` / `YOUR_CLIENT_SECRET` は手順3のウェブ用の値、`認可コード` は手順6でコピーした値に置き換え）:
   ```bash
   curl -X POST https://oauth2.googleapis.com/token \
     -d "client_id=YOUR_CLIENT_ID" \
     -d "client_secret=YOUR_CLIENT_SECRET" \
     -d "code=認可コード" \
     -d "grant_type=authorization_code" \
     -d "redirect_uri=http://localhost:8088/callback"
   ```
8. レスポンスの JSON に含まれる `refresh_token` を `backend/.env` の `GOOGLE_PHOTOS_REFRESH_TOKEN` に設定する。  
   **Client ID / Client Secret** は、このあとバックエンドで使うので、**ウェブアプリケーション用**のものを `GOOGLE_PHOTOS_CLIENT_ID` と `GOOGLE_PHOTOS_CLIENT_SECRET` に設定する（デスクトップ用から差し替えてよい）。

これでデスクトップの OOB が使えない環境でも Refresh Token を取得できます。

#### 5. 任意: アルバムIDの設定

特定のGoogle Photosアルバムにだけアップロードしたい場合は、[Google Photos Library API のドキュメント](https://developers.google.com/photos/library/guides/list-albums) でアルバム一覧を取得し、対象アルバムの ID を `GOOGLE_PHOTOS_ALBUM_ID` に設定します。未設定の場合は「ライブラリ」に保存されます。

#### 6. backend/.env の例

```env
GOOGLE_PHOTOS_CLIENT_ID=xxxxx.apps.googleusercontent.com
GOOGLE_PHOTOS_CLIENT_SECRET=GOCSPX-xxxxx
GOOGLE_PHOTOS_REFRESH_TOKEN=1//0xxxxx
# 任意
GOOGLE_PHOTOS_ALBUM_ID=
```

設定後、管理画面（http://localhost:3001）の魚画像登録でアップロードが利用できます。設定が不足している場合は「Google Photos設定が不足しています」というエラーが返ります。

#### トラブルシューティング: 401 invalid_client

認可URLを開いたときやトークン交換時に **「401: invalid_client」「flowName=GeneralOAuthFlow」** が出る場合の対処です。

| 原因 | 確認・対処 |
|------|------------|
| **Client ID の誤り** | 認可URLに貼っている `client_id` が、Google Cloud Console の「認証情報」に表示されている **クライアント ID** と完全に一致しているか確認。前後のスペース・改行が入っていないか、別プロジェクトの ID になっていないかも確認。 |
| **認証情報の種類と redirect_uri の不一致** | **ウェブアプリケーション** で作成したクライアントでは `redirect_uri=urn:ietf:wg:oauth:2.0:oob` は使えません。認可URLや curl では必ず、コンソールに登録した **認可済みリダイレクト URI** をそのまま使ってください（例: `http://localhost:8088/callback`）。デスクトップでやりたい場合は、**デスクトップアプリ** で認証情報を作り直すか、下記「localhost リダイレクト」を使います。 |
| **OOB が使えない環境** | 「デスクトップアプリ」なのに認可URLで 401 になる場合、Google 側で OOB（`urn:ietf:wg:oauth:2.0:oob`）が制限されていることがあります。→ **上記「デスクトップアプリでできないときの対応」** に従い、同じプロジェクトで「ウェブアプリケーション」の OAuth クライアントを追加して localhost リダイレクトで Refresh Token を取得してください。 |
| **Client Secret の誤り（curl 時）** | `.env` や curl の `client_secret` に余計なスペース・改行・引用符が入っていないか確認。コンソールで「シークレットを表示」して再コピー。 |
| **.env の書き方** | 値は引用符で囲まなくてよいです。`GOOGLE_PHOTOS_CLIENT_ID=xxxxx` のように `=` の直後に値。行末にスペースやカンマを入れない。 |

上記を確認しても 401 が続く場合は、同じプロジェクトで **新しい OAuth クライアント ID** を「デスクトップアプリ」で作成し、新しく出た Client ID / Client Secret でやり直してみてください。

#### 403 access_denied が出るとき

認可URLを開いたあと、「**403: access_denied**」「Developer Information アプリ名: …」のようなモーダルやエラーが出る場合は、**OAuth 同意画面が「テスト」状態のため、テストユーザーに登録されていないアカウントは認可できない**ことが原因です。

**対処手順**

1. [Google Cloud Console](https://console.cloud.google.com/) を開く
2. **APIとサービス** → **OAuth 同意画面**
3. 画面の **「テストユーザー」** セクション（または編集して「テストユーザー」のステップ）を開く
4. **「ADD USERS」／「ユーザーを追加」** をクリック
5. 認可で使っている **Google アカウントのメールアドレス**（例: `iwashi.kurukuru@gmail.com`）を入力して追加・保存
6. もう一度、認可URLをブラウザで開き直して「許可」をクリックする

これで 403 が解消され、認可コードが表示されるはずです。まだ 403 の場合は、ブラウザでログインしている Google アカウントと、追加したテストユーザーのメールアドレスが一致しているか確認してください。

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
