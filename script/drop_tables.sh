#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
BACKEND_ENV_FILE="$ROOT_DIR/backend/.env"

if [ -z "${DATABASE_URL:-}" ] && [ -f "$BACKEND_ENV_FILE" ]; then
	set -a
	# backend/.env の設定値をそのまま環境変数として読み込む
	. "$BACKEND_ENV_FILE"
	set +a
fi

if [ "$#" -lt 1 ]; then
	echo "使い方: $0 <table_name> [table_name ...]" >&2
	exit 1
fi

if [ -z "${DATABASE_URL:-}" ]; then
	echo "エラー: DATABASE_URL が設定されていません。" >&2
	exit 1
fi

if ! command -v psql >/dev/null 2>&1; then
	echo "エラー: psql が見つかりません。postgresql-client をインストールしてください。" >&2
	exit 1
fi

TABLE_LIST=""
for table_name in "$@"; do
	case "$table_name" in
		*[!A-Za-z0-9_]*|'')
			echo "エラー: 不正なテーブル名です: $table_name" >&2
			exit 1
			;;
	esac
	if [ -n "$TABLE_LIST" ]; then
		TABLE_LIST="$TABLE_LIST, "
	fi
	TABLE_LIST="$TABLE_LIST\"$table_name\""
done

SQL="DROP TABLE IF EXISTS $TABLE_LIST CASCADE;"

echo "実行SQL: $SQL"
psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -c "$SQL"
