build:
	go get github.com/aws/aws-lambda-go/lambda
	go get github.com/stretchr/testify/assert
	go get github.com/xeipuuv/gojsonschema
	go get github.com/aws/aws-sdk-go
	go get github.com/satori/go.uuid
	go get -u github.com/aws/aws-xray-sdk-go/...
	env GOOS=linux go build -ldflags="-s -w" -o bin/create_submission create_submission/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/dummy dummy/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/validate validate/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/submit_metadata submit_metadata/main.go