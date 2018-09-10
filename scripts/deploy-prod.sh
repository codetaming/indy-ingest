#!/usr/bin/env bash
make
cd serverless
sls deploy --stage prod -v