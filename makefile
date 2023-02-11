PROJECT_NAME := web-replay

ifdef OUTDIR
	OUTDIR := ${OUTDIR}
else
	OUTDIR := ./build
endif

ifeq (${PROD},true)
	LD_FLAGS := ${LD_FLAGS} -s -w
	OUTDIR := ${OUTDIR}/release
else
	LD_FLAGS := ${LD_FLAGS}
	OUTDIR := ${OUTDIR}/debug
endif

ifdef GOOS
	BUILD_TARGET := ${GOOS}
else
	BUILD_TARGET := $(shell go env GOOS)
endif

ifdef GOARCH
	BUILD_TARGET := ${BUILD_TARGET}-${GOARCH}
else
	BUILD_TARGET := ${BUILD_TARGET}-$(shell go env GOARCH)
endif

OUTDIR := ${OUTDIR}/${BUILD_TARGET}

ifeq (${GOOS},windows)
	BIN_NAME := $(PROJECT_NAME).exe
else ifeq ($(shell go env GOOS),windows)
	BIN_NAME := $(PROJECT_NAME).exe
else
	BIN_NAME := $(PROJECT_NAME)
endif

default: build

.PHONY: clean
clean:
	rm -r -f ./build

.PHONY: build
build:
	rm -rf ${OUTDIR}
	mkdir -p ${OUTDIR}
	mkdir -p ${OUTDIR}/bin
	
	cp -r commands/web-replay/patches ${OUTDIR}
	cp -r commands/web-replay-pack/web-replay-pack ${OUTDIR}/bin
	cp -r commands/web-replay-unpack/web-replay-unpack ${OUTDIR}/bin

	cd commands/web-replay && \
	env CGO_ENABLED=0 \
	go build -ldflags="$(LD_FLAGS)" -o "../../$(OUTDIR)/bin/$(BIN_NAME)" ./src/cmd/*.go

.PHONY: package
package:
	cd ${OUTDIR} && \
	tar -czvf "../${BUILD_TARGET}.tar.gz" *
 
.PHONY: build-all
build-all:
	env GOOS=linux GOARCH=arm make build
	env GOOS=linux GOARCH=arm64 make build
	env GOOS=linux GOARCH=amd64 make build
	env GOOS=darwin GOARCH=amd64 make build
	env GOOS=darwin GOARCH=arm64 make build
	env GOOS=windows GOARCH=amd64 make build
	env GOOS=windows GOARCH=arm64 make build

.PHONY: package-all
package-all:
	env GOOS=linux GOARCH=arm make package
	env GOOS=linux GOARCH=arm64 make package
	env GOOS=linux GOARCH=amd64 make package
	env GOOS=darwin GOARCH=amd64 make package
	env GOOS=darwin GOARCH=arm64 make package
	env GOOS=windows GOARCH=amd64 make package
	env GOOS=windows GOARCH=arm64 make package
