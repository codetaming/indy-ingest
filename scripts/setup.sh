#!/usr/bin/env bash
git pull
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
chmod +x $GOPATH/bin/dep
cd serverless
npm install
cd ..
$GOPATH/bin/dep ensure
make