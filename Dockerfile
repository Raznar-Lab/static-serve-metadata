# Stage 0: Build Go binary
FROM --platform=$TARGETOS/$TARGETARCH golang:1.24-alpine

WORKDIR /app
COPY . ./

# Build the Go binary (adjust the path and output binary name if needed)
RUN go build -o app-binary

# Stage 1: PHP container with built Go binary
FROM --platform=$TARGETOS/$TARGETARCH alpine:3.21

# Set working directory
WORKDIR /app

# Copy app source and built Go binary
COPY --from=0 /app/app-binary /usr/local/bin/app-binary
COPY .github/docker/entrypoint.sh /entrypoint.sh
COPY config.example.yml /app/config.example.yml
RUN chmod +x /usr/local/bin/app-binary

# Expose ports (e.g., for PHP server or reverse proxy)
EXPOSE 80

# Set entrypoint and default command
ENTRYPOINT ["/bin/ash", "/entrypoint.sh"]
