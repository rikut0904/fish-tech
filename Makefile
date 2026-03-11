GOCACHE := /tmp/go-build-cache
GOMODCACHE := /tmp/go-mod-cache

.PHONY: swag fmt test run build

swag:
	cd backend && GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g ./cmd/api/main.go -o ../frontend/public --outputTypes json --parseInternal

fmt: swag
	cd backend && gofmt -w ./cmd ./internal
	cd frontend && npm run lint -- --fix
	cd admin && npm run lint -- --fix

test: 
	cd backend && GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go test ./...
	cd frontend && npm run lint
	cd admin && npm run lint

up: 
	docker compose up --build

upd:
	docker compose up --build -d

build:
	docker compose build

down:
	docker compose down

logs:
	docker compose logs -f
