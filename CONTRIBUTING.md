# Contributing

By participating to this project, you agree to abide our [code of conduct](./CODE_OF_CONDUCT.md).

## Setup your machine

This project is written in [Go](https://golang.org/).

Prerequisites:

- Go 1.13+
- `make`

Fork and clone this repository.

Install modules:

```bash
make get
```

A good way of making sure everything is all right is running the following:

```bash
make build
```

## Local development

Run the following commands to start a local development docker:

```bash
make dev_docker
```

Run the following commands to start the `hub` server:

```bash
make dev_hub
```

Run the following commands to start the `indexer` server:

```bash
make dev_indexer
```

## Test your change

When you are satisfied with the changes, we suggest you run:

```
make lint
make test
```

## Pre-commit hook

We use [pre-commit](https://pre-commit.com/) to run tests before each commit.

1. Install pre-commit. See [pre-commit installation](https://pre-commit.com/#installation).
2. Run the following to install the hook:
   ```bash
   pre-commit install
   pre-commit install --hook-type commit-msg
   ```

## Submit a pull request

Push your branch to your fork and open a pull request against the `develop` branch.
