FROM golang:1.27.2-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/jobflow ./cmd/jobflow && CGO_ENABLED=0 go build -o /out/migrate ./cmd/migrate
FROM alpine:3.22.2
RUN adduser -D -u 10001 jobflow
WORKDIR /app
COPY --from=build /out/ /app/
COPY migrations/ /app/migrations/
ENV PATH="/app:${PATH}"
USER jobflow
EXPOSE 8080
CMD ["/bin/sh", "-c", "/app/migrate && exec /app/jobflow"]
