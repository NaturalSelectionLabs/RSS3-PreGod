FROM golang:1.18-alpine AS BUILDER

# Set the Current Working Directory inside the container
WORKDIR /rss3-pregod

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Install basic packages
RUN apk add \
    gcc \
    g++ \
    git

# Download all the dependencies
RUN go get ./hub/

# Build image
RUN go build -o dist/hub ./hub/

FROM alpine:latest AS RUNNER

COPY --from=builder /rss3-pregod/dist/hub .

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./hub"]
