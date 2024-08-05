FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build both binaries
RUN go build -o airdrop-server ./cmd/server/main.go && \
    go build -o airdrop-worker ./cmd/worker/main.go

EXPOSE 8080

# Copy the entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
