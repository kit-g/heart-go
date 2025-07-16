build-ApiFunction:
	GOOS=linux GOARCH=arm64 go build \
		-ldflags="-s -w \
			-X 'heart/internal/routerx.Commit=$(shell git rev-parse HEAD)' \
			-X 'heart/internal/routerx.DeployedAt=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)'" \
		-o $(ARTIFACTS_DIR)/bootstrap ./cmd/api/main.go

build-BackgroundFunction:
	GOOS=linux GOARCH=arm64 go build \
		-ldflags="-s -w \
			-X 'heart/internal/routerx.Commit=$(shell git rev-parse HEAD)' \
			-X 'heart/internal/routerx.DeployedAt=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)'" \
		-o $(ARTIFACTS_DIR)/bootstrap ./cmd/background/main.go


