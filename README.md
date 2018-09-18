# Indy Ingest

[![GoDoc][1]][2]
[![GoCard][3]][4]
[![Build Status][5]][6]
[![codecov][7]][8]
[![Codacy Badge][9]][10]
[![Docker Repository on Quay][11]][12]

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
[11]: https://quay.io/repository/codetaming/indy-ingest/status "Docker Repository on Quay"
[12]: https://quay.io/repository/codetaming/indy-ingest

*Personal learning project for Go with Serverless AWS or Kubernetes*

Provides an API for validation and storage of metadata/data with data storage S3 and state storage in DynamoDB.

Can be deployed to Kubernetes or serverless with AWS API Gateway and AWS Lambda.

## Original Objective
* Create an efficient, generic metadata / data submission system that performs well on minimal infrastructure
* The system must run performantly on a 3-node Kubernetes cluster of Rasperry Pi 3s
* The system must be capable of running as a AWS serverless application at minimal cost by keeping to the free tier whenever possible

![Rasperry Pi Cluster](./images/cluster.jpg =250x)
 
## Setup and Build
Requires [npm](https://docs.npmjs.com/cli/install) for serverless framework and [Go](https://golang.org/dl/) 1.11 with module support.

Run `./scripts/setup.sh` to resolve dependencies and build.

##Docker Build
Run the make file
```
make
```
Build Docker image
```
docker build -f Dockerfile-arm -t codetaming/indy-ingest-arm .
docker build -f Dockerfile -t codetaming/indy-ingest .
```

Run (Local)
```
docker run --publish 9000:9000 -t codetaming/indy-ingest
```

Run (Pi)
```
docker run --publish 9000:9000 -t codetaming/indy-ingest-arm
```

Push to Docker Hub
```
docker login
docker push codetaming/indy-ingest-arm
docker push codetaming/indy-ingest
```

## Uses
* Negroni : HTTP Middleware
* Gorilla Mux : URL router and dispatcher
* Envconfig : Configuration management
* JWT Middleware : Token based authentication