[![godoc](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://pkg.go.dev/github.com/georlav/httprawmock)
[![Tests](https://github.com/georlav/httprawmock/actions/workflows/ci.yml/badge.svg)](https://github.com/georlav/httprawmock/actions/workflows/ci.yml)
[![Golang-CI](https://github.com/georlav/httprawmock/actions/workflows/linter.yml/badge.svg)](https://github.com/georlav/httprawmock/actions/workflows/linter.yml)

# httprawmock
A simple http test server which allow the use of raw http response data for easy mocking APIs and web services. Create end-to-end tests for your http client by setting your own routes using custom patterns and a stub for each.

## Examples
Please check the [examples_test.go](examples_test.go) file for some basic usage examples.

## Install
```bash
go get github.com/georlav/httprawmock@latest
```

## Running tests
To run tests use
```go
go test -race ./... -v
```

## License
Distributed under the MIT License. See `LICENSE` for more information.

## Contact
George Lavdanis - georlav@gmail.com

Project Link: [https://github.com/georlav/httprawmock](https://github.com/georlav/httprawmock)