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

generate-webapp:
	yes | rm -rf web/dist/* && cd web && npm i && VITE_SERVER_ADDRESS="http://127.0.0.1/api/" npm run build
