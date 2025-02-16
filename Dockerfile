FROM golang:1.23-alpine
LABEL authors="kev"
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./
# Download the application dependencies
RUN go mod download
# Copy the application source code
COPY . .
# Build the Go application
RUN go build -o ravenclub cmd/raven-club/main.go
# Expose the port your application listens on
EXPOSE 7770
# Set the command to run the application
CMD ["./ravenclub"]

