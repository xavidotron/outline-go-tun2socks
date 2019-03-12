GOMOBILE=gomobile
GOBIND=$(GOMOBILE) bind
BUILDDIR=$(shell pwd)/build
IOS_BUILDDIR=$(BUILDDIR)/ios
ANDROID_BUILDDIR=$(BUILDDIR)/android
IOS_ARTIFACT=$(IOS_BUILDDIR)/Tun2socks.framework
ANDROID_ARTIFACT=$(ANDROID_BUILDDIR)/tun2socks.aar
LDFLAGS='-s -w'
IMPORT_PATH=github.com/Jigsaw-Code/go-tun2socks-mobile
TUN2SOCKS_PATH=$(GOPATH)/src/github.com/eycorsican/go-tun2socks

IOS_BUILD_CMD="cd $(BUILDDIR) && $(GOBIND) -a -ldflags $(LDFLAGS) -bundleid org.outline.tun2socks -target=ios -tags ios -o $(IOS_ARTIFACT) $(IMPORT_PATH)/ios"
ANDROID_BUILD_CMD="cd $(BUILDDIR) && $(GOBIND) -a -ldflags $(LDFLAGS) -target=android -tags android -o $(ANDROID_ARTIFACT) $(IMPORT_PATH)/android"

define build
	mkdir -p $(1)
	cd $(TUN2SOCKS_PATH) && make copy
	eval $(2)
	cd $(TUN2SOCKS_PATH) && make clean
endef

.PHONY: android ios clean

all: android ios

android:
	$(call build,$(ANDROID_BUILDDIR),$(ANDROID_BUILD_CMD))

ios:
	$(call build,$(IOS_BUILDDIR),$(IOS_BUILD_CMD))

clean:
	rm -rf $(BUILDDIR)
	cd $(TUN2SOCKS_PATH) && make clean
