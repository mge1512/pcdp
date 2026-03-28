.PHONY: build test install clean

build:
	CGO_ENABLED=0 go build -o pcd-lint .

test:
	go test -v ./...

install: build
	install -m 0755 pcd-lint /usr/local/bin/

clean:
	rm -f pcd-lint
	go clean -cache
