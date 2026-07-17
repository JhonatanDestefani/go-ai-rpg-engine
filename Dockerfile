# ---- Build stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# ---- Runtime stage ----
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/server ./server

USER app

ENTRYPOINT ["./server"]
