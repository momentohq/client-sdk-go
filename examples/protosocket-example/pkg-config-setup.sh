#!/bin/bash

set -e

# 1. Determine what platform we're on
PLATFORM=$(uname -s)
ARCH=$(uname -m)
echo "PLATFORM: $PLATFORM"
echo "ARCH: $ARCH"

# 2. Download the correct zip file from the latest momento-protosocket-ffi release
LATEST_VERSION=$(git ls-remote --tags --sort="v:refname" https://github.com/momentohq/momento-protosocket-ffi.git | tail -n 1 | sed 's!.*/v!!')
echo "LATEST_VERSION: $LATEST_VERSION"

if [[ "$PLATFORM" == "Darwin" && "$ARCH" == "x86_64" ]]; then
    echo "Downloading momento-protosocket-ffi.tar.gz for Darwin x86_64"
    curl -L -o momento-protosocket-ffi.tar.gz https://github.com/momentohq/momento-protosocket-ffi/releases/download/v$LATEST_VERSION/momento_protosocket_ffi-$LATEST_VERSION.x86_64-macos.tar.gz
elif [[ "$PLATFORM" == "Darwin" && "$ARCH" == "arm64" ]]; then
    echo "Downloading momento-protosocket-ffi.tar.gz for Darwin arm64"
    curl -L -o momento-protosocket-ffi.tar.gz https://github.com/momentohq/momento-protosocket-ffi/releases/download/v$LATEST_VERSION/momento_protosocket_ffi-$LATEST_VERSION.arm64-macos.tar.gz

elif [[ "$PLATFORM" == "Linux" && "$ARCH" == "x86_64" ]]; then
    echo "Downloading momento-protosocket-ffi.tar.gz for Linux x86_64"
    curl -L -o momento-protosocket-ffi.tar.gz https://github.com/momentohq/momento-protosocket-ffi/releases/download/v$LATEST_VERSION/momento_protosocket_ffi-$LATEST_VERSION.x86_64-linux.tar.gz
elif [[ "$PLATFORM" == "Linux" && "$ARCH" == "arm64" ]]; then
    echo "Downloading momento-protosocket-ffi.tar.gz for Linux arm64"
    curl -L -o momento-protosocket-ffi.tar.gz https://github.com/momentohq/momento-protosocket-ffi/releases/download/v$LATEST_VERSION/momento_protosocket_ffi-$LATEST_VERSION.arm64-linux.tar.gz
else
    echo "Unsupported platform or architecture"
    exit 1
fi

# 3. Unzip the tar.gz file and place in new directory
mkdir -p protosocket-ffi
mv momento-protosocket-ffi.tar.gz protosocket-ffi/
cd protosocket-ffi
tar -xzf momento-protosocket-ffi.tar.gz
CURRENT_DIR=$(pwd)
echo "CURRENT_DIR: $CURRENT_DIR"

# 4. Update the first line to use the correct path in the pkg-config file
sed -i '' "1s|.*|libdir=$CURRENT_DIR|g" ./momento_protosocket_ffi.pc

# # 5. Update PKG_CONFIG_PATH
export PKG_CONFIG_PATH=$CURRENT_DIR
echo "[REQUIRED] Make sure you set the environment variable PKG_CONFIG_PATH=$PKG_CONFIG_PATH"

# # 6. Verify that pkg-config picks up on the library
FOUND=$(pkg-config --list-package-names | grep momento_protosocket_ffi)
if [ -z "$FOUND" ]; then
    echo "ERROR: pkg-config could not find momento_protosocket_ffi"
    exit 1
else
    echo "pkg-config found momento_protosocket_ffi, should be good to go"
fi
