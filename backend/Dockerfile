FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/worldwright ./cmd/server

FROM alpine:3.20 AS runtime

RUN apk add --no-cache wget ca-certificates && \
    addgroup -S app && adduser -S -G app app && \
    mkdir -p /data /data/assets && chown -R app:app /data

COPY --from=build /out/worldwright /usr/local/bin/worldwright

USER app
WORKDIR /data
ENV WORLDWRIGHT_LISTEN_ADDR=":8080" \
    WORLDWRIGHT_DATABASE_PATH="/data/worldwright.db" \
    WORLDWRIGHT_ASSETS_DIR="/data/assets"

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/worldwright"]
