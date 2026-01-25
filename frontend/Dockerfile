# ビルドステージ
FROM node:24.9.0-alpine AS builder

WORKDIR /app

# パッケージファイルをコピー
COPY package*.json ./
RUN npm ci

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN npm run build

# 実行ステージ
FROM node:24.9.0-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production

# ビルダーから必要なファイルをコピー
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/node_modules ./node_modules

EXPOSE 3000

CMD ["npm", "start"]
