# ビルドステージ
FROM node:24.9.0-alpine AS builder

WORKDIR /app
RUN npm install -g npm@latest

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
RUN npm install -g npm@latest

ENV NODE_ENV=production
ENV PORT=3001

# ビルダーから必要なファイルをコピー
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/node_modules ./node_modules

EXPOSE 3001

CMD ["npm", "start"]
