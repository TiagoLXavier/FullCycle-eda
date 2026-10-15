FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /out/walletcore ./cmd/walletcore
RUN CGO_ENABLED=1 GOOS=linux go build -o /out/balances ./cmd/balances

FROM debian:bookworm-slim AS walletcore

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/walletcore /usr/local/bin/walletcore

EXPOSE 8080

CMD ["walletcore"]

FROM debian:bookworm-slim AS balances

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/balances /usr/local/bin/balances

EXPOSE 3003

CMD ["balances"]
