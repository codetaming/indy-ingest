#!/usr/bin/env bash
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
~/environment/go/bin/dep ensure
make
cd serverless
sls deploy --stage test -v