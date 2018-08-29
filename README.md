# Indy Ingest

[![GoDoc][1]][2]
[![GoCard][3]][4]
[![Build Status][5]][6]
[![codecov][7]][8]
[![Codacy Badge][9]][10]

[1]: https://godoc.org/github.com/codetaming/indy-ingest?status.svg
[2]: https://godoc.org/github.com/codetaming/indy-ingest
[3]: https://goreportcard.com/badge/github.com/codetaming/indy-ingest
[4]: https://goreportcard.com/report/github.com/codetaming/indy-ingest
[5]: https://travis-ci.org/codetaming/indy-ingest.svg?branch=master
[6]: https://travis-ci.org/codetaming/indy-ingest
[7]: https://codecov.io/gh/codetaming/indy-ingest/branch/master/graph/badge.svg
[8]: https://codecov.io/gh/codetaming/indy-ingest
[9]: https://api.codacy.com/project/badge/Grade/b75a9233c6064ba4a61c70e44fbaae26
[10]: https://www.codacy.com/app/danielvaughan/indy-ingest?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=codetaming/indy-ingest&amp;utm_campaign=Badge_Grade

*Personal learning project for GoLang and AWS*

Creates and API for submission and validation of JSON via API Gateway and Lambda with storage in data S3 and state storage in DynamoDB.

*Setup and Build*
Requires [npm](https://docs.npmjs.com/cli/install) for serverless framework and [Go](https://golang.org/dl/) with [dep](https://github.com/golang/dep) for code.

Run `./scripts/setup.sh` to resolve dependencies and build.

##Docker Build
Run the make file
```
make
```
Build Docker image
```
docker build -f Dockerfile-arm -t codetaming/ingest-arm .
docker build -f Dockerfile -t codetaming/ingest .
```

Run (Local)
```
docker run --publish 9000:9000 -t codetaming/ingest
```

Run (Pi)
```
docker run --publish 9000:9000 -t codetaming/ingest-arm
```

Push to Docker Hub
```
docker login
docker push codetaming/ingest-arm
docker push codetaming/ingest 
```

## Uses
* Negroni : HTTP Middleware
* Gorilla Mux : URL router and dispatcher
* Envconfig : Configuration management
* JWT Middleware : Token based authentication