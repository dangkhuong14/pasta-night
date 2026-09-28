# syntax=docker/dockerfile:1

# Build a static binary (no cgo) so it runs on distroless/static.
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server \
 && mkdir -p /out/data/cache

# distroless/static has CA certificates (HTTPS to TMDB) and runs as nonroot (uid 65532).
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/server /app/server
COPY configs /app/configs
# A named volume starts with this directory's ownership, so nonroot can write the cache.
COPY --from=build --chown=65532:65532 /out/data /app/data
ENV CACHE_DIR=/app/data/cache
EXPOSE 8080
VOLUME ["/app/data/cache"]
ENTRYPOINT ["/app/server"]
