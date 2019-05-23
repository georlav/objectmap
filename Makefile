build:
	go build -o objectmap -mod=vendor -ldflags "-s -w" cmd/objectmap/objectmap.go
test:
	CGO_ENABLED=1 go test ./... -v -race -cover -count=1
lint:
	golangci-lint run