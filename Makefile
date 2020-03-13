PROJECT=chaos-commander
GOPATH ?= $(shell go env GOPATH)

# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

CURDIR := $(shell pwd)
path_to_add := $(addsuffix /bin,$(subst :,/bin:,$(GOPATH))):$(PWD)/tools/bin
export PATH := $(path_to_add):$(PATH)

GO        := GO111MODULE=off go
GOBUILD   := CGO_ENABLED=0 $(GO) build $(BUILD_FLAG) -tags codes
GOTEST    := CGO_ENABLED=0 $(GO) test -p 4
OVERALLS  := CGO_ENABLED=0 GO111MODULE=on overalls

ARCH      := "`uname -s`"
LINUX     := "Linux"
MAC       := "Darwin"
PACKAGE_LIST  := go list ./...| grep -vE "cmd" | grep -vE "test"
PACKAGES  := $$($(PACKAGE_LIST))
PACKAGE_DIRECTORIES := $(PACKAGE_LIST) | sed 's|github.com/pingcap/$(PROJECT)/||'
FILES     := $$(find $$($(PACKAGE_DIRECTORIES)) -name "*.go")

FILES     := $$(find . -name "*.go" | grep -vE "vendor")
GOFILTER := grep -vE 'vendor|render.Delims|bindata_assetfs|testutil|\.pb\.go'
GOCHECKER := $(GOFILTER) | awk '{ print } END { if (NR > 0) { exit 1 } }'
GOLINT := go list ./... | grep -vE 'vendor' | xargs -L1 -I {} golint {} 2>&1 | $(GOCHECKER)

FAILPOINT_ENABLE  := $$(find $$PWD/ -type d | grep -vE "(\.git|tools)" | xargs tools/bin/failpoint-ctl enable)
FAILPOINT_DISABLE := $$(find $$PWD/ -type d | grep -vE "(\.git|tools)" | xargs tools/bin/failpoint-ctl disable)

CHECK_LDFLAGS += $(LDFLAGS)

.PHONY: all syncer pcp

default: all

all: syncer pcp

build:
	$(GOBUILD)

RACE_FLAG =
ifeq ("$(WITH_RACE)", "1")
	RACE_FLAG = -race
	GOBUILD   = GOPATH=$(GOPATH) CGO_ENABLED=1 $(GO) build
endif

CHECK_FLAG =
ifeq ("$(WITH_CHECK)", "1")
	CHECK_FLAG = $(TEST_LDFLAGS)
endif

syncer:
	$(GOBUILD) $(RACE_FLAG) -ldflags '$(LDFLAGS) $(CHECK_FLAG)' -o bin/syncer cmd/syncer/*.go

pcp:
	$(GOBUILD) $(RACE_FLAG) -ldflags '$(LDFLAGS) $(CHECK_FLAG)' -o bin/pcp cmd/pcp/*.go

check:
	GO111MODULE=off go get golang.org/x/lint/golint

	@echo "vet"

	${GO} vet -all ./...

	@echo "golint"
	@ $(GOLINT)
	@echo "gofmt"
	@gofmt -s -l -w $(FILES) 2>&1 | $(GOCHECKER)
