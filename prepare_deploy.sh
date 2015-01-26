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

DEPLOY_DIR="avorsmtp_${VER_MAJOR}.${VER_MINOR}.${VER_PATCH}"
DEPLOY_FILE_NAME="avorsmtp_${VER_MAJOR}.${VER_MINOR}.${VER_PATCH}.tar.gz"

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

echo "Deploy build complete."
echo "Live well"
