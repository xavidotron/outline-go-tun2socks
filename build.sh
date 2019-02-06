#!/bin/bash -eux

BUILD_DIR=build
TUN2SOCKS_FRAMEWORK=Tun2socks.framework

rm -rf $BUILD_DIR
make
pushd $BUILD_DIR/$TUN2SOCKS_FRAMEWORK > /dev/null

# Get the framework in the correct format.
# Remove symlinks
rm Headers Modules Resources Tun2socks
mv Versions/A/* .
rm -rf Versions Resources
pushd > /dev/null

# Add Info.plist
cp Info.plist $BUILD_DIR/$TUN2SOCKS_FRAMEWORK/
