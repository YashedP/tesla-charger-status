FROM golang:1.24-alpine AS builder
WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/tesla-charger-status ./cmd/server

FROM alpine:3.20
RUN adduser -D -H -u 10001 app
WORKDIR /app

COPY --from=builder /out/tesla-charger-status /app/tesla-charger-status
COPY scripts /app/scripts

RUN mkdir -p /app/data /app/secrets && chown -R app:app /app
USER app

EXPOSE 5000
ENTRYPOINT ["/app/tesla-charger-status"]
