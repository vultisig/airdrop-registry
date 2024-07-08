test:
	go test -v ./...

dev-worker:
	gow run cmd/worker/main.go

dev-server:
	gow run cmd/server/main.go

worker:
	go run cmd/worker/main.go

server:
	go run cmd/server/main.go
