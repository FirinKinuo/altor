VERSION := $(file < VERSION)
PROJECT := $(shell cat go.mod | grep module | sed 's|module\s||g')

clean:
	rm -rf _build/ release/

build:
	go mod download
	CGO_ENABLED=0 go build -tags release -o $(PROJECT) main.go

build-all:
	mkdir -p _build
	GOOS=darwin  GOARCH=amd64   CGO_ENABLED=0 go build -tags release -o _build/$(PROJECT)-$(VERSION)-darwin-amd64 cmd/main/main.go
	GOOS=linux   GOARCH=amd64   CGO_ENABLED=0 go build -tags release -o _build/$(PROJECT)-$(VERSION)-linux-amd64 cmd/main/main.go
	GOOS=linux   GOARCH=arm     CGO_ENABLED=0 go build -tags release -o _build/$(PROJECT)-$(VERSION)-linux-arm cmd/main/main.go
	GOOS=linux   GOARCH=arm64   CGO_ENABLED=0 go build -tags release -o _build/$(PROJECT)-$(VERSION)-linux-arm64 cmd/main/main.go
	GOOS=windows GOARCH=amd64   CGO_ENABLED=0 go build -tags release -o _build/$(PROJECT)-$(VERSION)-windows-amd64.exe cmd/main/main.go
	cd _build; sha256sum * > sha256sums.txt

install:
	make build
	mv ./$(PROJECT) /usr/bin/$(PROJECT)

release:
	make clean
	make build-all
	mkdir release
	cp _build/* release
	cd release; sha256sum --quiet --check sha256sums.txt && \
	gh release create v$(VERSION) -d -t v$(VERSION) *
