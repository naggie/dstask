.PHONY: release clean install
dist/dstask: clean
	go build -mod=vendor -o dist/dstask cmd/dstask/main.go
	go build -mod=vendor -o dist/dstask-import cmd/dstask-import/main.go

release:
	./do-release.sh

clean:
	rm -rf dist

install:
	cp dist/dstask /usr/local/bin
	cp dist/dstask-import /usr/local/bin

test:
	go test ./...

update_deps:
	go get
	go mod vendor
	git add -f vendor
