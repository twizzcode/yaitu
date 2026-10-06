# ---------- Build stage ----------
FROM golang:1.27-alpine AS build

WORKDIR /app

# Cache dependency dulu.
COPY go.mod go.sum ./
RUN go mod download

# Build binary statis.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/lapanganku-api .

# ---------- Runtime stage ----------
FROM alpine:3.20 AS runtime

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=build /out/lapanganku-api /app/lapanganku-api

# Folder migrations disalin agar bisa dijalankan manual bila perlu.
COPY --from=build /app/migrations /app/migrations

ENV PORT=8080

EXPOSE 8080

CMD ["/app/lapanganku-api"]
