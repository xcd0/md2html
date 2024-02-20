BIN        := md2html
DST        := ../build
GOARCH     := amd64
VERSION    := `gobump show | gojq -r .version`
REVISION   := `git rev-parse --short HEAD`

FLAGS_UNIX := -ldflags="-X main.version=$(VERSION) -X main.revision=$(REVISION) -s -w -extldflags='-static' -buildid=" -a -tags=netgo -trimpath
FLAGS_WIN  := -ldflags="-H windowsgui -X main.version=$(VERSION) -X main.revision=$(REVISION) -s -w -extldflags='-static' -buildid=" -a -tags=netgo -trimpath 

#.PHONY: build
.PHONY: all
build:
	bash update_version_info.sh
	mkdir -p $(DST)
ifeq ($(shell uname -o),Msys)
	CGO_ENABLED=0 go build -o $(DST)/$(BIN).exe
else
	CGO_ENABLED=0 go build -o $(DST)/$(BIN)
endif
	cp -rf $(DST)/$(BIN) ..

cp: build
	cp -rf $(DST)/$(BIN) `which $(BIN)`

run: build
	cd .. && ./$(BIN) readme.md

#all: build
all: release-build

release-build:
	bash update_version_info.sh
	goimports -w *.go
	gofmt -w *.go
	go generate
	CGO_ENABLED=0 GOARCH=$(GOARCH) GOOS=windows CGO_ENABLED=0 go build -o $(DST)/$(BIN)_windows.exe $(FLAGS_WIN)
	CGO_ENABLED=0 GOARCH=$(GOARCH) GOOS=darwin CGO_ENABLED=0 go build -o $(DST)/$(BIN)_macOS $(FLAGS_UNIX)
	CGO_ENABLED=0 GOARCH=$(GOARCH) GOOS=linux CGO_ENABLED=0 go build -o $(DST)/$(BIN)_linux $(FLAGS_UNIX)

release: release-build
	# for windows
	cd $(DST) && \
	mv $(BIN)_windows.exe $(BIN).exe && \
	upx --lzma $(BIN).exe && \
	zip md2html_binary_$(GOARCH)_windows.zip md2html.exe && \
	mv $(BIN).exe $(BIN)_windows.exe
	# for mac
	cd $(DST) && \
	mv $(BIN)_macOS $(BIN) && \
	upx --lzma $(BIN) && \
	zip md2html_binary_$(GOARCH)_macOS.zip md2html && \
	mv $(BIN) $(BIN)_macOS
	# for linux
	cd $(DST) && \
	mv $(BIN)_linux $(BIN) && \
	upx --lzma $(BIN) && \
	zip md2html_binary_$(GOARCH)_linux.zip md2html && \
	mv $(BIN) $(BIN)_linux
	# clean
	cd $(DST) && \
	rm $(DST)/$(BIN)_windows.exe && \
	rm $(DST)/$(BIN)_macOS && \
	rm $(DST)/$(BIN)_linux


clean:
	rm -rf build
	rm -f *.html
	rm -f *.mini.css

get:
	go get
	go install github.com/x-motemen/gobump/cmd/gobump@latest
	go install github.com/itchyny/gojq/cmd/gojq@latest
	go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest
	go install golang.org/x/tools/cmd/goimports@latest
	
	@echo このほかにリリースにupxを使用しています。
	@echo https://github.com/upx/upx/releases/latest から取得してください。
