# Generate API docs.
# https://github.com/swaggo/swag

swag fmt --dir $(pwd)/hub
swag init --output $(pwd)/hub/docs --dir $(pwd)/hub,$(pwd)/shared
