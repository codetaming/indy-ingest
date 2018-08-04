build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/create_dataset cmd/create_dataset/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/dummy cmd/dummy/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/validate cmd/validate/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/add_metadata cmd/add_metadata/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/list_datasets cmd/list_datasets/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get_dataset cmd/get_dataset/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get_metadata cmd/get_metadata/main.go