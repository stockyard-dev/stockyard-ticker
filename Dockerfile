FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/ticker ./cmd/ticker/
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata curl
COPY --from=builder /bin/ticker /usr/local/bin/ticker
ENV PORT="9020" DATA_DIR="/data"
EXPOSE 9020
HEALTHCHECK --interval=30s --timeout=5s CMD curl -sf http://localhost:9020/health || exit 1
ENTRYPOINT ["ticker"]
