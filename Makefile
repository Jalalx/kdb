GIT_HASH=$(shell git rev-parse --short HEAD)
VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.GitHash=$(GIT_HASH)"

.PHONY: clean build run scripts

clean:
	rm -f bin/kdb
	rm -f bin/ask
	rm -f bin/learn
	rm -f bin/kdb-backup
	rm -f bin/text-transform

run:
	go run *.go

build:
	mkdir -p bin/
	cp -rf scripts/ask bin/ask
	cp -rf scripts/learn bin/learn
	cp -rf scripts/kdb-backup bin/kdb-backup
	cp -rf scripts/text-transform bin/text-transform
	go build $(LDFLAGS) -o bin/kdb

scripts:
	mkdir -p bin/
	cp -rf scripts/ask bin/ask
	cp -rf scripts/learn bin/learn
	cp -rf scripts/kdb-backup bin/kdb-backup
	cp -rf scripts/text-transform bin/text-transform

dist: build
	temp_dir=$(shell mktemp -d)
	mkdir -p "$temp_dir/macos-darwin-arm64"
	cp bin/* "$temp_dir/macos-darwin-arm64/"
	mkdir -p "dist"
	zip -r dist/macos-darwin-arm64.zip -j "$temp_dir/macos-darwin-arm64"
	rm -rf "$temp_dir"
