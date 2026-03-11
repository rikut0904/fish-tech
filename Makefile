GOCACHE := /tmp/go-build-cache
GOMODCACHE := /tmp/go-mod-cache

.PHONY: swag fmt lint test run up upd build down logs

swag:
	cp frontend/public/swagger.json /tmp/fish-tech-swagger.prev.json 2>/dev/null || true
	cd backend && GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g ./cmd/api/main.go -o ../frontend/public --outputTypes json --parseInternal
	node ./script/update_swagger_version.js ./backend/cmd/api/main.go ./frontend/public/swagger.json /tmp/fish-tech-swagger.prev.json
	rm -f /tmp/fish-tech-swagger.prev.json

fmt: swag
	cd backend && GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go fix ./...
	cd backend && GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go run golang.org/x/tools/cmd/goimports@v0.31.0 -w ./cmd ./internal
	cd backend && gofmt -w ./cmd ./internal
	cd frontend && npm run lint -- --fix
	cd admin && npm run lint -- --fix

lint:
	cd frontend && npm run lint
	cd admin && npm run lint

test: fmt
	cd backend && GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go test ./...

up: 
	docker compose up --build

upd:
	docker compose up --build -d

build: swag
	docker compose build

down:
	docker compose down

logs:
	docker compose logs -f
