FROM golang:alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /airdrop-registry

EXPOSE 8080

CMD ["/airdrop-registry"]
