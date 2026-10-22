# artifact:Dockerfile | v1.0.0 | 2026-08-19T06:54:00Z | Topic: Delivery | Multi-stage Go build

FROM golang:1.26-alpine AS build

WORKDIR /src

RUN apk add --no-cache \
    ca-certificates \
    git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
      -trimpath \
      -ldflags="-s -w" \
      -o /out/jobflow \
      ./cmd/jobflow

FROM alpine:3.22 AS runtime

RUN apk add --no-cache \
    ca-certificates \
    tzdata

RUN addgroup -S jobflow \
    && adduser -S -G jobflow jobflow

WORKDIR /app

COPY --from=build /out/jobflow /app/jobflow
COPY migrations /app/migrations

USER jobflow

EXPOSE 8080

ENTRYPOINT ["/app/jobflow"]