.PHONY: build, release

build:
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/dbbackup-linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/dbbackup-linux-arm64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/dbbackup-darwin-amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/dbbackup-darwin-arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/dbbackup-windows-amd64
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o bin/dbbackup-windows-arm64

release:
	echo $(version)
	@echo "Release Version: $(version)"
	@bash -c "echo -n $(version) > VERSION"
	git add .
	git commit -m $(version)
	git tag $(version)