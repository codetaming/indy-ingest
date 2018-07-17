#!/usr/bin/env bash
git pull
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
dep ensure
make
sls deploy -v