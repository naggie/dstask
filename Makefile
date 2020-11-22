.PHONY: release clean install
dist/dstask: clean
	go build -mod=vendor -o dist/dstask cmd/dstask/main.go
	go build -mod=vendor -o dist/dstask-sync cmd/dstask-sync/main.go

release:
	./do-release.sh

clean:
	rm -rf dist

install:
	cp dist/dstask /usr/local/bin
	cp dist/dstask-sync /usr/local/bin

test:
	go test ./...

update_deps:
	go get -i
	go mod vendor
