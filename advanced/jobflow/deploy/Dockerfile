FROM golang:1.27.2-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir /out && for svc in pub sub act ret demo; do CGO_ENABLED=0 go build -trimpath -o /out/$svc ./cmd/$svc; done

FROM alpine:3.22.2
RUN apk add --no-cache ca-certificates && adduser -D -u 10001 jobflow
WORKDIR /app
COPY --from=build /out/ /app/
COPY migrations/ /app/migrations/
ENV PATH="/app:${PATH}"
USER jobflow
EXPOSE 8080
CMD ["pub"]
