VERSION=$(shell git describe --tags --always)

.PHONY: build_all
# build
build_all:
	rm -rf bin && mkdir bin bin/linux-amd64 bin/linux-arm64 bin/darwin-amd64 bin/darwin-arm64 \
	&& CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'main.Version=$(VERSION)'" -o ./bin/darwin-arm64/ ./... \
	&& CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.Version=$(VERSION)'" -o ./bin/darwin-amd64/ ./... \
	&& CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-X 'main.Version=$(VERSION)'" -o ./bin/linux-arm64/ ./... \
	&& CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.Version=$(VERSION)'" -o ./bin/linux-amd64/ ./...

.PHONY: build
# build
build:
	rm -rf bin && mkdir bin && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...