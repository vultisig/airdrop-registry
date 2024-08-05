FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build both binaries
RUN go build -o airdrop-server ./cmd/server/main.go && \
    go build -o airdrop-worker ./cmd/worker/main.go

EXPOSE 8080

# Use an entrypoint script to allow choosing which binary to run
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
