#!/usr/bin/env bash
git pull
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
~/environment/go/bin/dep ensure
make
sls deploy -v