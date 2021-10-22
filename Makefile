regenerate:
	go generate ./...
test:
	go test ./...

build:
	go build -o app/api cmd/api/main.go

run:
	go run cmd/api/main.go