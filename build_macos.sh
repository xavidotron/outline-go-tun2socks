#!/bin/bash -eux
#
# Copyright 2019 The Outline Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

BUILD_DIR=build/macos
BIN_DIR=bin/macos
TUN2SOCKS_BINARY=$BUILD_DIR/go-tun2socks-macos
# TODO: use Jigsaw developer certificate
CERT_NAME="Mac Developer: Alberto Lalama (6U3H9CUW4N)"

rm -rf $BUILD_DIR
make clean && make macos
codesign -f --prefix org.outline. --entitlements macos/go-tun2socks-macos.entitlements -s "$CERT_NAME" $TUN2SOCKS_BINARY
cp -R $TUN2SOCKS_BINARY $BIN_DIR/
