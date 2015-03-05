#!/bin/bash
#
# 
# Description : Prepare deploy avorsmtp. 
# Author : vassilux
# Last modified : 2014-10-10 14:53:54  
#

set -e

VER_MAJOR="1"
VER_MINOR="0"
VER_PATCH="1"

VERSION=$(cat VERSION)

cp main.go main.go.bkp

sed -i "/VERSION = \"X.X.X\"/c\VERSION = \"${VERSION}\"" main.go

make clean

make fmt

make

if [ ! -f ./bin/avorsmtp ]; then
	echo "Can not find compiled  project file ./bin/avorsmtp."
	echo "Please cheque make output."
    exit 1
fi

mv main.go.bkp main.go 

DEPLOY_DIR="avorsmtp_${VERSION}"
DEPLOY_FILE_NAME="avorsmtp_${VERSION}.tar.gz"

if [ -d "$DEPLOY_DIR" ]; then
    rm -rf  "$DEPLOY_DIR"
fi
#
#
mkdir "$DEPLOY_DIR"
mkdir "$DEPLOY_DIR/samples"
mkdir "$DEPLOY_DIR/logs"

cp -aR ./bin/* "$DEPLOY_DIR"
cp -aR ./samples/* "${DEPLOY_DIR}/samples"
#
mkdir "$DEPLOY_DIR/docs"
pandoc -o "$DEPLOY_DIR/docs/INSTALL.html" ./docs/INSTALL.md
pandoc -o "$DEPLOY_DIR/docs/ReleaseNotes.html" ./docs/ReleaseNotes.md
cp "$DEPLOY_DIR/docs/INSTALL.html" .
cp "$DEPLOY_DIR/docs/ReleaseNotes.html" .

tar cvzf "${DEPLOY_FILE_NAME}" "${DEPLOY_DIR}"

if [ ! -f "$DEPLOY_FILE_NAME" ]; then
    echo "Deploy build failed."
    exit 1
fi

rm -rf "$DEPLOY_DIR"

if [ ! -d releases ]; then
	mkdir releases
fi

mv INSTALL.html ./releases
mv ReleaseNotes.html ./releases
mv ${DEPLOY_FILE_NAME} ./releases


echo "Deploy build complete."
echo "Live well"
