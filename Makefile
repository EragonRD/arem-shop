.PHONY: run build test fmt docker-up docker-down docker-logs

# ── Développement local (sans Docker) ─────────────────────────
run:
	go run ./cmd/api

build:
	go build ./...

test:
	go test ./...

fmt:
	gofmt -w ./cmd ./internal

# ── Docker ────────────────────────────────────────────────────
docker-up:
	bash scripts/docker-start.sh

docker-up-detach:
	bash scripts/docker-start.sh -d

docker-down:
	bash scripts/docker-stop.sh

docker-clean:
	bash scripts/docker-stop.sh --volumes

docker-logs:
	bash scripts/docker-logs.sh

docker-test:
	bash scripts/docker-test.sh
