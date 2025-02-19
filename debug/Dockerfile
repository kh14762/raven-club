FROM golang:1.23-alpine AS build

# Build Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /
# Copy the Go module files
COPY go.mod go.sum ./
# Download the application dependencies
RUN go mod download
# Copy the application source code
COPY . .
# Build the Go application
RUN go build -gcflags="all=-N -l" -o ravenclub cmd/raven-club/main.go
# Expose the port your application listens on
EXPOSE 7771

FROM alpine AS debug
COPY --from=build /go/bin/dlv /dlv
COPY --from=build /ravenclub /

CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app"]