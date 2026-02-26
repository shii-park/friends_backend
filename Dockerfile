# --- build stage ---
FROM golang:1.26-alpine AS builder
WORKDIR /app

# HTTPSや時刻まわり（推奨）
RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /bin/app ./cmd/app

# --- runtime stage ---
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /bin/app /app/app

# 実行に必要なSQLファイルとJSONファイルをコピー
COPY internal/sqlc/schema/schema.sql internal/sqlc/schema/schema.sql
COPY internal/domain/cards.json internal/domain/cards.json

EXPOSE 8080
ENV GIN_MODE=release PORT=8080
ENTRYPOINT ["/app/app"]
