.PHONY: build

build:
	@mkdir -p bin
	#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/dbbackup-linux-amd64 ./cmd/cli
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/dbbackup-linux-arm64 ./cmd/cli
	#CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/dbbackup-darwin-amd64 ./cmd/cli
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/dbbackup-darwin-arm64 ./cmd/cli
	#CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/dbbackup-windows-amd64 ./cmd/cli
	@#CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o bin/dbbackup-windows-arm64