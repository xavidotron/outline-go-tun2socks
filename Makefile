GOCMD=go
GOMOBILE=gomobile
GOBIND=$(GOMOBILE) bind
GOBUILD=$(GOCMD) build
BUILDDIR=$(shell pwd)/build
ANDROID_BUILDDIR=$(BUILDDIR)/android
ANDROID_ARTIFACT=$(ANDROID_BUILDDIR)/tun2socks.aar
IOS_BUILDDIR=$(BUILDDIR)/ios
IOS_ARTIFACT=$(IOS_BUILDDIR)/Tun2socks.framework
MACOS_BUILDDIR=$(BUILDDIR)/macos
MACOS_BINARY=go-tun2socks-macos
MACOS_ARTIFACT=$(MACOS_BUILDDIR)/$(MACOS_BINARY)
LDFLAGS='-s -w'
IMPORT_PATH=github.com/Jigsaw-Code/go-tun2socks-mobile
TUN2SOCKS_PATH=$(GOPATH)/src/github.com/eycorsican/go-tun2socks

ANDROID_BUILD_CMD="cd $(BUILDDIR) && $(GOBIND) -a -ldflags $(LDFLAGS) -target=android -tags android -o $(ANDROID_ARTIFACT) $(IMPORT_PATH)/android"
IOS_BUILD_CMD="cd $(BUILDDIR) && $(GOBIND) -a -ldflags $(LDFLAGS) -bundleid org.outline.tun2socks -target=ios -tags ios -o $(IOS_ARTIFACT) $(IMPORT_PATH)/ios"
MACOS_BUILD_CMD="cd $(BUILDDIR) && $(GOBUILD) -ldflags $(LDFLAGS) -o $(MACOS_ARTIFACT) $(MACOS_BINARY)"

define build
	mkdir -p $(1)
	cd $(TUN2SOCKS_PATH) && make copy
	eval $(2)
	cd $(TUN2SOCKS_PATH) && make clean
endef

.PHONY: android ios macos clean

all: android ios macos

android:
	$(call build,$(ANDROID_BUILDDIR),$(ANDROID_BUILD_CMD))

ios:
	$(call build,$(IOS_BUILDDIR),$(IOS_BUILD_CMD))

macos:
	$(call build,$(MACOS_BUILDDIR),$(MACOS_BUILD_CMD))

clean:
	rm -rf $(BUILDDIR)
	cd $(TUN2SOCKS_PATH) && make clean
