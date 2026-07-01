# STAGE 1: Build the application
FROM golang:1.25-alpine AS builder

# Add metadata to the image
LABEL maintainer="oremi64@gmail.com"
LABEL version="1.0"
LABEL description="ASCII Art Web Service"

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files first to leverage Docker cache
COPY go.mod ./
RUN go mod download

# Copy the source code (including banners, templates, and static files)
COPY . .

# Build the application with optimizations
# -o main: names the output binary
# CGO_ENABLED=0: ensures a static binary for compatibility
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# STAGE 2: Create the final lightweight image
FROM alpine:latest

# Add metadata to the final container object
LABEL type="production"

# Set working directory
WORKDIR /root/

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the app runs on
EXPOSE 8080

# Run the binary
CMD ["./main"]
