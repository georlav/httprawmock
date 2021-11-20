package httprawmock_test

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/georlav/httprawmock"
)

func Test_MockJSONAPI(t *testing.T) {
	type input struct {
		method     string
		url        string
		urlPattern string
		body       io.Reader
	}

	testCases := []struct {
		description        string
		input              input
		responseFile       string
		expectedStatusCode int
	}{
		{
			description: "Should create a unicorn",
			input: input{
				method:     http.MethodPost,
				url:        "/unicorns",
				urlPattern: "/unicorns",
				body:       strings.NewReader(`{"name":"Sparkle Angel","age":2,"colour":"pink"}`),
			},
			responseFile:       "testdata/post_unicorn.txt",
			expectedStatusCode: http.StatusCreated,
		},
		{
			description: "Should fail to create a unicorn due to validation error",
			input: input{
				method:     http.MethodPost,
				url:        "/unicorns",
				urlPattern: "/unicorns",
				body:       strings.NewReader(`{"name":"Sparkle Angel","age":2,"colour":"pink"}`),
			},
			responseFile:       "testdata/post_unicorn_400.txt",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "Should retrieve a unicorn",
			input: input{
				method:     http.MethodGet,
				url:        "/unicorns/6198f9da97069d03e849096d",
				urlPattern: "/unicorns/{id}",
			},
			responseFile:       "testdata/get_unicorn.txt",
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "Should fail to retrieve a unicorn",
			input: input{
				method:     http.MethodGet,
				url:        "/unicorns/6198f9da97069d03e849096c",
				urlPattern: "/unicorns/{id}",
			},
			responseFile:       "testdata/get_unicorn_404.txt",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			description: "Should retrieve multiple unicorns",
			input: input{
				method:     http.MethodGet,
				url:        "/unicorns",
				urlPattern: "/unicorns",
			},
			responseFile:       "testdata/get_unicorns.txt",
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "Should update a unicorn",
			input: input{
				method:     http.MethodPut,
				url:        "/unicorns/6198f9da97069d03e849096c",
				urlPattern: "/unicorns/{id}",
			},
			responseFile:       "testdata/put_unicorn.txt",
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "Should update a unicorn (no method is set)",
			input: input{
				method:     "",
				url:        "/unicorns/6198f9da97069d03e849096c",
				urlPattern: "/unicorns/{id}",
			},
			responseFile:       "testdata/put_unicorn.txt",
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "Should delete a unicorn",
			input: input{
				method:     "",
				url:        "/unicorns/6198f9da97069d03e849096c",
				urlPattern: "/unicorns/{id}",
			},
			responseFile:       "testdata/delete_unicorn.txt",
			expectedStatusCode: http.StatusOK,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.description, func(t *testing.T) {
			b, err := os.ReadFile(tc.responseFile)
			if err != nil {
				t.Fatalf("failed to load testfile, %s", err)
			}

			// Register route (one per case)
			ts := httprawmock.NewServer(
				httprawmock.NewRoute(tc.input.method, tc.input.urlPattern, b),
			)
			t.Cleanup(ts.Close)

			if tc.input.method != "" {
				rts, err := ts.GetRoutes()
				if err != nil {
					t.Fatal(err)
				}

				if cntRoutes := len(rts); cntRoutes != 1 {
					t.Fatalf("Expected to have 1 registered routes got %d. Routes: %+v", cntRoutes, rts)
				}
			}

			// build request
			req, err := http.NewRequest(tc.input.method, ts.URL+tc.input.url, nil)
			if err != nil {
				t.Fatalf("failed to create request, %s", err)
			}

			// send the request
			resp, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}

			// check code
			if tc.expectedStatusCode != resp.StatusCode {
				t.Fatalf("Expected to have status %d got %d", tc.expectedStatusCode, resp.StatusCode)
			}

			_ = resp
			// t.Logf("%+v", resp)
		})
	}
}

func Test_MockJSONRegisterMultipleRoutesAtOnce(t *testing.T) {
	type input struct {
		method string
		url    string
		body   io.Reader
	}

	testCases := []struct {
		description        string
		input              input
		expectedStatusCode int
	}{
		{
			description: "Should create a unicorn",
			input: input{
				method: http.MethodPost,
				url:    "/unicorns",
				body:   strings.NewReader(`{"name":"Sparkle Angel","age":2,"colour":"pink"}`),
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			description: "Should retrieve a unicorn",
			input: input{
				method: http.MethodGet,
				url:    "/unicorns/6198f9da97069d03e849096d",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "Should retrieve multiple unicorns",
			input: input{
				method: http.MethodGet,
				url:    "/unicorns",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "Should update a unicorn",
			input: input{
				method: http.MethodPut,
				url:    "/unicorns/6198f9da97069d03e849096c",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "Should fail to meet a unicorn (endpoint does not exists)",
			input: input{
				method: http.MethodGet,
				url:    "/unicorns/6198f9da97069d03e849096c/meet",
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	// Register route (one per case)
	ts := httprawmock.NewServer(
		httprawmock.NewRoute(http.MethodGet, "/unicorns", []byte(`HTTP/1.1 200 OK
Server: nginx/1.14.2
Date: Sat, 20 Nov 2021 13:45:22 GMT
Content-Type: application/json; charset=utf-8
Transfer-Encoding: chunked
Connection: keep-alive
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range
Access-Control-Expose-Headers: Content-Length,Content-Range

[{"_id":"6198f9cc97069d03e849096b","name":"Sparkle Angel","age":2,"colour":"blue"},{"_id":"6198f9ce97069d03e849096c","name":"Sparkle Angel","age":2,"colour":"blue"},{"_id":"6198f9da97069d03e849096d","name":"Sparkle Angel","age":2,"colour":"blue"}]`)),
		httprawmock.NewRoute(http.MethodGet, "/unicorns/{id}", []byte(`HTTP/1.1 200 OK
Server: nginx/1.14.2
Date: Sat, 20 Nov 2021 13:39:23 GMT
Content-Type: application/json; charset=utf-8
Transfer-Encoding: chunked
Connection: keep-alive
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range
Access-Control-Expose-Headers: Content-Length,Content-Range

{"_id":"6198f9da97069d03e849096d","name":"Sparkle Angel","age":2,"colour":"blue"}
`)),
		httprawmock.NewRoute(http.MethodPost, "/unicorns", []byte(`HTTP/1.1 201 Created
Server: nginx/1.14.2
Date: Sat, 20 Nov 2021 13:36:26 GMT
Content-Type: application/json; charset=utf-8
Transfer-Encoding: chunked
Connection: keep-alive
Location: /api/8777798e0dc14ddd924c4c3b1567f4c4/unicorns/6198f9da97069d03e849096d
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range
Access-Control-Expose-Headers: Content-Length,Content-Range

{"name":"Sparkle Angel","age":2,"colour":"blue","_id":"6198f9da97069d03e849096d"}`)),
		httprawmock.NewRoute(http.MethodPut, "/unicorns/{id}", []byte(`HTTP/1.1 200 OK
Server: nginx/1.14.2
Date: Sat, 20 Nov 2021 13:39:23 GMT
Content-Type: application/json; charset=utf-8
Transfer-Encoding: chunked
Connection: keep-alive
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range
Access-Control-Expose-Headers: Content-Length,Content-Range`)),
		httprawmock.NewRoute(http.MethodDelete, "/unicorns/{id}", []byte(`HTTP/1.1 200 OK
Server: nginx/1.14.2
Date: Sat, 20 Nov 2021 19:02:19 GMT
Content-Length: 0
Connection: keep-alive
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range
Access-Control-Expose-Headers: Content-Length,Content-Range`)),
	)
	t.Cleanup(ts.Close)

	rts, err := ts.GetRoutes()
	if err != nil {
		t.Fatal(err)
	}

	if cntRoutes := len(rts); cntRoutes != 5 {
		t.Fatalf("Expected to have 5 registered routes got %d", cntRoutes)
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.description, func(t *testing.T) {
			// build request
			req, err := http.NewRequest(tc.input.method, ts.URL+tc.input.url, nil)
			if err != nil {
				t.Fatalf("failed to create request, %s", err)
			}

			// send the request
			resp, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}

			// check code
			if tc.expectedStatusCode != resp.StatusCode {
				t.Fatalf("Expected to have status %d got %d", tc.expectedStatusCode, resp.StatusCode)
			}

			_ = resp
			// t.Logf("%+v", resp)
		})
	}
}
