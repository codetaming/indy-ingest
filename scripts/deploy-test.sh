#!/usr/bin/env bash
make
cd serverless
sls deploy --stage test -v