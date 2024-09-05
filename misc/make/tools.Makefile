# This makefile should be used to hold functions/variables

ARCH_ := $(ARCH)
ifeq ($(ARCH),x86_64)
	ARCH := amd64
else ifeq ($(ARCH),aarch64)
	ARCH := arm64 
endif

define github_url
    https://github.com/$(GITHUB)/releases/download/v$(VERSION)/$(ARCHIVE)
endef

# creates a directory bin.
bin:
	@ mkdir -p $@

# ~~~ Tools ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

# ~~ [ goose ] ~~~ https://github.com/pressly/goose ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# https://github.com/pressly/goose/releases/download/v3.22.0/goose_linux_x86_64
GOOSE := $(shell command -v goose || echo "bin/goose")
goose: bin/goose ## Install goose (database migration)

bin/goose: VERSION := 3.22.0
bin/goose: GITHUB  := pressly/goose
# bin/goose: ARCHIVE := goose_linux_$(ARCH)
bin/goose: ARCHIVE := goose_linux_$(ARCH_)
bin/goose: bin
	@ printf "Install goose... "
	@ curl -Ls $(call github_url) > $@ && chmod +x $@
	@ echo "done."

# ~~ [ air ] ~~~ https://github.com/cosmtrek/air ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# https://github.com/air-verse/air/releases/download/v1.52.3/air_1.52.3_linux_amd64

AIR := $(shell command -v air || echo "bin/air")
air: bin/air ## Installs air (go file watcher)

bin/air: VERSION := 1.52.3
bin/air: GITHUB  := air-verse/air
bin/air: ARCHIVE := air_$(VERSION)_$(OSTYPE)_$(ARCH)
bin/air: bin
	@ printf "Install air... "
	@ curl -Ls $(call github_url) > $@ && chmod +x $@
	@ echo "done."

# ~~ [ gotestsum ] ~~~ https://github.com/gotestyourself/gotestsum ~~~~~~~~~~~~~~~~~~~~~~~
# https://github.com/gotestyourself/gotestsum/releases/download/v1.12.0/gotestsum_1.12.0_linux_amd64.tar.gz

GOTESTSUM := $(shell command -v gotestsum || echo "bin/gotestsum")
gotestsum: bin/gotestsum ## Installs gotestsum (testing go code)

bin/gotestsum: VERSION := 1.12.0
bin/gotestsum: GITHUB  := gotestyourself/gotestsum
bin/gotestsum: ARCHIVE := gotestsum_$(VERSION)_$(OSTYPE)_$(ARCH).tar.gz
bin/gotestsum: bin
	@ printf "Install gotestsum... "
	@ curl -Ls $(shell echo $(call github_url) | tr A-Z a-z) | tar -zOxf - gotestsum > $@ && chmod +x $@
	@ echo "done."

# ~~ [ tparse ] ~~~ https://github.com/mfridman/tparse ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# https://github.com/mfridman/tparse/releases/download/v0.14.0/tparse_linux_x86_64

TPARSE := $(shell command -v tparse || echo "bin/tparse")
tparse: bin/tparse ## Installs tparse (testing go code)

bin/tparse: VERSION := 0.14.0
bin/tparse: GITHUB  := mfridman/tparse
bin/tparse: ARCHIVE := tparse_$(OSTYPE)_$(ARCH_)
bin/tparse: bin
	@ printf "Install tparse... "
	@ curl -Ls $(call github_url) > $@ && chmod +x $@
	@ echo "done."

# ~~ [ mockery ] ~~~ https://github.com/vektra/mockery ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# https://github.com/vektra/mockery/releases/download/v2.45.0/mockery_2.45.0_Linux_x86_64.tar.gz

MOCKERY := $(shell command -v mockery || echo "bin/mockery")
mockery: bin/mockery ## Installs mockery (mocks generation)

bin/mockery: VERSION := 2.45.0
bin/mockery: GITHUB  := vektra/mockery
bin/mockery: ARCHIVE := mockery_$(VERSION)_$(OSTYPE)_$(ARCH_).tar.gz
bin/mockery: bin
	@ printf "Install mockery... "
	@ curl -Ls $(call github_url) | tar -zOxf -  mockery > $@ && chmod +x $@
	@ echo "done."

# ~~ [ golangci-lint ] ~~~ https://github.com/golangci/golangci-lint ~~~~~~~~~~~~~~~~~~~~~
https://github.com/golangci/golangci-lint/releases/download/v1.60.3/golangci-lint-1.60.3-linux-amd64.tar.gz

GOLANGCI := $(shell command -v golangci-lint || echo "bin/golangci-lint")
golangci-lint: bin/golangci-lint ## Installs golangci-lint (linter)

bin/golangci-lint: VERSION := 1.60.3
bin/golangci-lint: GITHUB  := golangci/golangci-lint
bin/golangci-lint: ARCHIVE := golangci-lint-$(VERSION)-$(OSTYPE)-$(ARCH).tar.gz
bin/golangci-lint: bin
	@ printf "Install golangci-linter... "
	@ curl -Ls $(shell echo $(call github_url) | tr A-Z a-z) | tar -zOxf - $(shell printf golangci-lint-$(VERSION)-$(OSTYPE)-$(ARCH)/golangci-lint | tr A-Z a-z ) > $@ && chmod +x $@
	@ echo "done."