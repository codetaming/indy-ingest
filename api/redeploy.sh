#!/usr/bin/env bash
git pull
dep ensure
make
sls deploy -v