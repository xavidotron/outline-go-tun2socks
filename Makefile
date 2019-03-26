GOCMD=go
GOMOBILE=gomobile
GOBIND=$(GOMOBILE) bind
GOBUILD=$(GOCMD) build
BUILDDIR=$(shell pwd)/build
IMPORT_PATH=github.com/Jigsaw-Code/outline-go-tun2socks
LDFLAGS='-s -w'
TUN2SOCKS_VERSION=v1.13.0
TUN2SOCKS_SRC_PATH=$(GOPATH)/src/github.com/eycorsican/go-tun2socks
TUN2SOCKS_MOD_PATH=$(GOPATH)/pkg/mod/github.com/eycorsican/go-tun2socks\@$(TUN2SOCKS_VERSION)

ANDROID_BUILDDIR=$(BUILDDIR)/android
ANDROID_ARTIFACT=$(ANDROID_BUILDDIR)/tun2socks.aar
IOS_BUILDDIR=$(BUILDDIR)/ios
IOS_ARTIFACT=$(IOS_BUILDDIR)/Tun2socks.framework
MACOS_BUILDDIR=$(BUILDDIR)/macos
MACOS_ARTIFACT=$(MACOS_BUILDDIR)/Tun2socks.framework
WINDOWS_BUILDDIR=$(BUILDDIR)/windows

ANDROID_BUILD_CMD="cd $(BUILDDIR) && GO111MODULE=off $(GOBIND) -a -ldflags $(LDFLAGS) -target=android -tags android -o $(ANDROID_ARTIFACT) $(IMPORT_PATH)/android"
IOS_BUILD_CMD="cd $(BUILDDIR) &&  GO111MODULE=off $(GOBIND) -a -ldflags $(LDFLAGS) -bundleid org.outline.tun2socks -target=ios/arm,ios/arm64 -tags ios -o $(IOS_ARTIFACT) $(IMPORT_PATH)/apple"
MACOS_BUILD_CMD="cd $(BUILDDIR) &&  GO111MODULE=off $(GOBIND) -a -ldflags $(LDFLAGS) -bundleid org.outline.tun2socks -target=ios/amd64 -tags ios -o $(MACOS_ARTIFACT) $(IMPORT_PATH)/apple"
WINDOWS_BUILD_CMD="cd $(TUN2SOCKS_SRC_PATH) && go get -d ./... && BUILD_TAGS='dnscache dnsfallback socks' make windows && cp build/tun2socks-windows*.exe $(WINDOWS_BUILDDIR)/"

define build
	$(call modularize)
	mkdir -p $(1)
	cd $(TUN2SOCKS_MOD_PATH) && make copy
	eval $(2)
	cd $(TUN2SOCKS_MOD_PATH) && make clean
	$(call undo_modularize)
endef

# Workaround to modularize go-tun2socks and gomobile.
define modularize
	# We need to call `make copy` in go-tun2socks, but the downloaded
	# module does not grant us write permissions.
	# TODO: add module support in go-tun2socks upstream.
	chmod -R u+w $(TUN2SOCKS_MOD_PATH)
	# gomobile does not yet support modules.
	# Symlink the current module and the go-tun2socks module in $GOPATH.
	# go-tun2socks should not be in $GOPATH for this to work.
	# TODO: remove this once gomobile enables modules in Go 1.13.
	ln -s $(shell pwd) $(GOPATH)/src/$(IMPORT_PATH)
	ln -s $(TUN2SOCKS_MOD_PATH) $(TUN2SOCKS_SRC_PATH)
endef

define undo_modularize
	rm $(GOPATH)/src/$(IMPORT_PATH) || true
	rm $(TUN2SOCKS_SRC_PATH) || true
endef

.PHONY: android ios macos windows clean

all: android ios macos windows

android:
	$(call build,$(ANDROID_BUILDDIR),$(ANDROID_BUILD_CMD))

ios:
	$(call build,$(IOS_BUILDDIR),$(IOS_BUILD_CMD))

macos:
	$(call build,$(MACOS_BUILDDIR),$(MACOS_BUILD_CMD))

windows:
	$(call build,$(WINDOWS_BUILDDIR),$(WINDOWS_BUILD_CMD))

clean:
	rm -rf $(BUILDDIR)
	cd $(TUN2SOCKS_MOD_PATH) && make clean || true
	$(call undo_modularize)
