#!/usr/bin/env bash
git pull
curl -L -s https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 -o $GOPATH/bin/dep
chmod +x $GOPATH/bin/dep
cd serverless
npm install
cd ..
$GOPATH/bin/dep ensure
make