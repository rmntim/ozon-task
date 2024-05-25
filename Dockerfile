FROM golang:alpine as builder

WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v -o server /app/cmd/server/main.go

# Use alpine to run the binary
FROM alpine

WORKDIR /app

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /app/server
COPY --from=builder /app/config/default.yaml /app/config.yaml
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

# Run the web service on container startup.
CMD ["/app/server", "--config", "/app/config.yaml"]