
USER_GH=eyedeekay
VERSION=$(ZERO_VERSION_B).001
ZERO_VERSION_B=`./get_latest_release.sh "i2p-zero/i2p-zero" | tr -d 'v'`
packagename=zerobundle
ZERO_VERSION=`./get_latest_release.sh "i2p-zero/i2p-zero"`
LINE_ITEM=`grep 'var zeroversion' gen.go`
LINE_ITEM_2=`grep 'var ZERO_VERSION' import/import.go`

echo:
	@echo "type make version to do release $(VERSION)"

all: echo fmt gen build

build:
	go build -tags netgo -o zero ./go-zero

gen: write test

update:
	sed -i "s|$(LINE_ITEM)|var zeroversion = \"$(ZERO_VERSION)\"|g" gen.go
	sed -i "s|$(LINE_ITEM_2)|var ZERO_VERSION = \"$(ZERO_VERSION)\"|g" import/import.go

write: update
	go run --tags=generate gen.go

itest:
	cd import/ && go test -v

fmt:
	gofmt -w -s *.go import/*.go go-zero/*.go  parts/*/unpacker.go

version:
	gothub release -s $(GITHUB_TOKEN) -u $(USER_GH) -r $(packagename) -t v$(VERSION) -d "version $(VERSION)"

del:
	gothub delete -s $(GITHUB_TOKEN) -u $(USER_GH) -r $(packagename) -t v$(VERSION)

tar:
	tar --exclude .git \
		--exclude .go \
		--exclude bin \
		--exclude examples \
		-cJvf ../$(packagename)_$(VERSION).orig.tar.xz .

