build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/pipeline cmd/pipeline/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/create_dataset cmd/serverless/create_dataset/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/dummy cmd/serverless/dummy/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/validate cmd/serverless/validate/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/add_metadata cmd/serverless/add_metadata/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/list_datasets cmd/serverless/list_datasets/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get_dataset cmd/serverless/get_dataset/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get_metadata cmd/serverless/get_metadata/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/ingest main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/ic cmd/client/main.go
	env GOOS=darwin GOARCH=amd64  go build -ldflags="-s -w" -o bin/ingest-osx main.go
	env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o bin/ingest-arm main.go
	env GOOS=darwin GOARCH=amd64  go build -ldflags="-s -w" -o pipeline-osx cmd/pipeline/main.go
