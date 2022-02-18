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

## Test your change

When you are satisfied with the changes, we suggest you run:

```
make lint
make test
```

## Submit a pull request

Push your branch to your fork and open a pull request against the `develop` branch.
