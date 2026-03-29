FROM golang:1.26 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o collector-app ./collector/cmd

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/collector-app .
EXPOSE 50051
CMD ["./collector-app"]