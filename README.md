[![godoc](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://pkg.go.dev/github.com/georlav/httprawmock)
[![Tests](https://github.com/georlav/httprawmock/actions/workflows/ci.yml/badge.svg)](https://github.com/georlav/httprawmock/actions/workflows/ci.yml)
[![Golang-CI](https://github.com/georlav/httprawmock/actions/workflows/linter.yml/badge.svg)](https://github.com/georlav/httprawmock/actions/workflows/linter.yml)

# httprawmock
A simple http test server which allow the use of raw http response data for easy mocking APIs and web services. Create end-to-end tests for your http client by setting your own routes using custom patterns and a stub for each.

## Examples
Please check the [examples_test.go](examples_test.go) file for some basic usage examples. Also check some  basic examples from [another project](https://github.com/georlav/bitstamp/blob/master/httpapi_test.go) that uses httprawmock and table driven tests.
```go
    response := `HTTP/1.1 200 OK
Server: nginx/1.14.2
Date: Sat, 20 Nov 2021 13:39:23 GMT
Content-Type: application/json; charset=utf-8
Transfer-Encoding: chunked
Connection: keep-alive
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range
Access-Control-Expose-Headers: Content-Length,Content-Range

{"_id":"6198f9da97069d03e849096d","name":"Sparkle Angel","age":2,"colour":"blue"}`

    // Register your routes
	ts := httprawmock.NewServer(
		httprawmock.NewRoute(http.MethodGet, "/unicorns/{id}", []byte(response)),
	)
	defer ts.Close()

    resp, err := http.DefaultClient.Get(ts.URL + "/unicorns/6198f9da97069d03e849096d")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

    b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(b))
```

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