.PHONY: release clean install
dist/dstask: clean
	mkdir -p dstask
	go build -mod=vendor -o dist/dstask cmd/dstask.go

release:
	./do-release.sh

clean:
	rm -rf dist

install:
	cp dist/dstask /usr/local/bin
