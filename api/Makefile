build:
	go get github.com/aws/aws-lambda-go/lambda
	go get github.com/stretchr/testify/assert
	go get github.com/xeipuuv/gojsonschema
	env GOOS=linux go build -ldflags="-s -w" -o bin/dummy dummy/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/validate validate/main.go