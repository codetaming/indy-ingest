#!/usr/bin/env bash
git pull
make
cd serverless
sls deploy -v