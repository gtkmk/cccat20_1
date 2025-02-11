# Stage 1: build the Go app
FROM golang:1.23-alpine AS builder
RUN apk update && apk add --no-cache ca-certificates tzdata
WORKDIR /app
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./cccat20_1

# Stage 2: Create a smaller runtime image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/cccat20_1 .
COPY --from=builder /app/.env /root/.env

ENV TZ America/Sao_Paulo

EXPOSE 8089
CMD ["./cccat20_1"]
