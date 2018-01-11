#!/usr/bin/env bash

PROTODEP_VERSION="v0.0.1-patch"
PROTODEP_TMP_PATH="/tmp/protodep.tar.gz"
PROTODEP_BIN="./node_modules/.bin/protodep"

initArch() {
	ARCH=$(uname -m)
	case $ARCH in
		armv5*) ARCH="armv5";;
		armv6*) ARCH="armv6";;
		armv7*) ARCH="armv7";;
		aarch64) ARCH="arm64";;
		x86) ARCH="386";;
		x86_64) ARCH="amd64";;
		i686) ARCH="386";;
		i386) ARCH="386";;
	esac
	echo "ARCH=$ARCH"
}

initOS() {
	OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')
	echo "OS=$OS"
}

download() {
	rm -rf $PROTODEP_TMP_PATH
  local url="https://github.com/franciscocpg/protodep/releases/download/${PROTODEP_VERSION}/protodep_${OS}_${ARCH}.tar.gz"
  echo "Downloading protodep from $url"
  eval "curl -${SILENT}L "$url" -o $PROTODEP_TMP_PATH"
}

install() {
	if [[ ! -f $PROTODEP_TMP_PATH ]]
	then
		download
	fi

  tar xvf $PROTODEP_TMP_PATH
  chmod +x protodep
  mv protodep ./node_modules/.bin
}

main() {
  if [[ ! -f $PROTODEP_BIN ]]
  then
    initArch
    initOS
    install
  fi
}

main
