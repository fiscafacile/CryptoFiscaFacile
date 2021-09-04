#!/bin/bash

PATH=$GOPATH/bin:$PATH

VER=`git describe --tags`

BASENAME="CryptoFiscaFacile"

GOOS=`go env | grep GOOS | cut -d"=" -f2`
GOARCH=`go env | grep GOARCH | cut -d"=" -f2`

echo $GOOS/$GOARCH

if [ $GOOS = "windows" ]
then
	EXT=".exe"
else
	EXT=""
fi
go build -ldflags "-X main.version=$VER" && mv $BASENAME$EXT $BASENAME-$VER-$GOOS-$GOARCH$EXT
