GIT_REPO   := github.com/naggie/dstask

DIST_FILE  := dstask
DIST_DIR   := dist
SRC_FILE   := $(addsuffix .go, $(DIST_FILE))
SRC_DIR    := cmd
ARM5_FILE  := $(DIST_FILE)-linux-arm5
AMD64_FILE := $(DIST_FILE)-linux-amd64
DRWN_FILE  := $(DIST_FILE)-darwin-amd64

RELEASE_FILE := RELEASE.md

# VERSION HAS TO BE SET BY HAND
VERSION    := ""
GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILD_DATE := $(shell date)

LDFLAGS  = -s -w
LDFLAGS += -X \"$(GIT_REPO).GIT_COMMIT=$(GIT_COMMIT)\"
LDFLAGS += -X \"$(GIT_REPO).VERSION=$(VERSION)\"
LDFLAGS += -X \"$(GIT_REPO).BUILD_DATE=$(BUILD_DATE)\"

export CGO_ENABLED=0

install: $(DIST_DIR)/$(DIST_FILE)
	install $< /usr/local/bin/

$(DIST_DIR)/$(DIST_FILE): $(SRC_DIR)/$(SRC_FILE) $(DIST_DIR)
	go build -mod=vendor -o $@ $<

$(DIST_DIR):
	mkdir -p $@
	
release: \
	$(DIST_DIR)/$(ARM5_FILE) \
	$(DIST_DIR)/$(AMD64_FILE) \
	$(DIST_DIR)/$(DRWN_FILE) \
	$(RELEASE_FILE)
	hub release create \
	--draft \
	-a $(DIST_DIR)/$(ARM5_FILE)\\#"dstask linux-arm5" \
	-a $(DIST_DIR)/$(AMD64_FILE)\\#"dstask linux-amd64" \
	-a $(DIST_DIR)/$(DRWN_FILE)\\#"dstask darwin-amd64" \
	-F $(RELEASE_FILE) \
	$(VERSION)

$(DIST_DIR)/$(ARM5_FILE): export GOOS=linux
$(DIST_DIR)/$(ARM5_FILE): export GOARCH=arm
$(DIST_DIR)/$(ARM5_FILE): export GOARM=5
$(DIST_DIR)/$(ARM5_FILE): $(SRCDIR)/$(SRCFILE) $(DISTDIR)
	go build -mod=vendor -ldflags="$(LDFLAGS)" -o $@ $<

$(DIST_DIR)/$(AMD64_FILE): export GOOS=linux
$(DIST_DIR)/$(AMD64_FILE): export GOARCH=amd64
$(DIST_DIR)/$(AMD64_FILE): $(SRCDIR)/$(SRCFILE) $(DISTDIR)
	go build -mod=vendor -ldflags="$(LDFLAGS)" -o $@ $<

$(DIST_DIR)/$(DRWN_FILE): export GOOS=darwin
$(DIST_DIR)/$(DRWN_FILE): export GOARCH=amd64
$(DIST_DIR)/$(DRWN_FILE): $(SRCDIR)/$(SRCFILE) $(DISTDIR)
	go build -mod=vendor -ldflags="$(LDFLAGS)" -o $@ $<

$(RELEASE_FILE):
	# file doesn't exist or is for old version, replace
	( !-f $@ || head -n 1 $@ | grep -vq $(VERSION)) && printf "$(VERSION)\n\n\n" > $@
	vim "+ normal G $" $@

.PHONY: clean

clean:
	rm -rf $(DISTDIR)
