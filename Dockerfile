# Copyright (c) 2026 Flick. All Rights Reserved.

# --- Builder stage ---
FROM golang:1.26.2-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

# Optimisations: Prevent downloading dependencies at every builds
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Set build informations
ARG VERSION=selfbuild
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

ENV CGO_ENABLED=0 GOOS=linux

RUN VERSION="${VERSION}" \
    COMMIT="${COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo unknown)}" \
    BUILD_DATE="${BUILD_DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}" \
    CLI_PKG="github.com/matteoepitech/flick/internal/cli" \
    && go build -trimpath -ldflags="-s -w -X ${CLI_PKG}.CLIVersion=${VERSION} -X ${CLI_PKG}.CLICommit=${COMMIT} -X ${CLI_PKG}.CLIBuildDate=${BUILD_DATE}" \
        -o /out/flick-api ./cmd/api

# --- Runtime stage ---
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

WORKDIR /app

COPY --from=builder /out/flick-api /app/flick-api

EXPOSE 15702/tcp
USER nonroot:nonroot

ENTRYPOINT ["/app/flick-api"]
