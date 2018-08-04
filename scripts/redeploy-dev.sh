#!/usr/bin/env bash
git pull
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
$GOPATH/bin/dep ensure
make
cd serverless
sls deploy -v