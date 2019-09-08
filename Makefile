.PHONY: release clean unstall
dist/dstask:
	mkdir -p dstask
	go build -mod=vendor -o dist/dstask cmd/dstask.go

release:
	./do-release.sh

clean:
	rm -rf dist

install: dist/dstask
	cp dist/dstask /usr/local/bin
