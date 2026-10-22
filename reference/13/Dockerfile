FROM golang:1.27.2-alpine AS build
WORKDIR /src
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/jobflow ./cmd/jobflow

FROM alpine:3.22.2
RUN adduser -D -u 10001 jobflow
WORKDIR /app
COPY --from=build /out/jobflow /app/jobflow
USER jobflow
EXPOSE 8080
ENTRYPOINT ["/app/jobflow"]
