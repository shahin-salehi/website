# Build.
FROM golang:1.24 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o /entrypoint ./cmd/server

# Deploy.
FROM alpine:latest 
WORKDIR /
COPY --from=build-stage /entrypoint /entrypoint
COPY --from=build-stage /app/static /static
EXPOSE 8080
ENTRYPOINT ["/entrypoint"]