#!/usr/bin/env bash
cd serverless
npm install
cd ..
dep ensure
make