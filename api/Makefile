build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/create_dataset create_dataset/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/dummy dummy/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/validate validate/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/add_metadata add_metadata/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/list_datasets list_datasets/main.go