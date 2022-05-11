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
RUN go get ./reptile/

# Build image
RUN go build -o dist/reptile ./reptile/

FROM alpine:latest AS RUNNER

COPY --from=builder /rss3-pregod/dist/reptile .

# Run the executable
CMD ["./reptile"]
