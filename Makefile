.PHONY: all

BUILD_FILE=main.go

TARGET_NAME=feature
VERSION=1.0.0
RELEASE=2.el7
BUILT=`date`

GOOS=linux
GOARCH=amd64
ARCH=x86_64
GO111MODULE=on

ENV=prd

all: pull format clean build compress package
prd: pull clean build package

pull:
	git pull
	git submodule foreach git pull

format:
	gofmt -w .

build:
	$(eval GIT_COMMIT=$(shell git rev-parse --short HEAD))
	GOOS=$(GOOS) GOARCH=$(GOARCH) GO111MODULE=$(GO111MODULE) go build -ldflags "-X main.Version=$(VERSION) \
	-X main.CommitId=$(GIT_COMMIT) -X 'main.Built=$(BUILT) -s -w'" -v -o $(TARGET_NAME) $(BUILD_FILE)

clean:
	rm -rf ./objs $(TARGET_NAME)

package:
	dos2unix ./scripts/*.sh
	chmod +x ./scripts/*.sh
	./scripts/package.sh package $(TARGET_NAME) $(ENV)

release:
	dos2unix ./scripts/*.sh
	chmod +x ./scripts/*.sh
	./scripts/package.sh release $(TARGET_NAME)

compress:
	upx $(TARGET_NAME)

rpmbuild:
	./scripts/package.sh rpmbuild $(TARGET_NAME) $(VERSION)-$(RELEASE) $(ARCH)
