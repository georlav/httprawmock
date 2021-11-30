package httprawmock_test

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/georlav/httprawmock"
)

func ExampleNewServer() {
	resp1, _ := os.ReadFile("testdata/post_unicorn.txt")
	resp2, _ := os.ReadFile("testdata/get_unicorns.txt")
	resp3, _ := os.ReadFile("testdata/get_unicorn.txt")
	resp4, _ := os.ReadFile("testdata/put_unicorn.txt")
	resp5, _ := os.ReadFile("testdata/delete_unicorn.txt")

	// Register your routes according to you needs
	ts := httprawmock.NewServer(
		httprawmock.NewRoute(http.MethodPost, "/unicorns", resp1),
		httprawmock.NewRoute(http.MethodGet, "/unicorns", resp2),
		httprawmock.NewRoute(http.MethodGet, "/unicorns/{id}", resp3),
		httprawmock.NewRoute(http.MethodPut, "/unicorns/{id}", resp4),
		httprawmock.NewRoute(http.MethodDelete, "/unicorns/{id}", resp5),
	)
	defer ts.Close()

	// Retrieve all registered routes (for debuging purposes)
	rts, err := ts.GetRoutes()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i := range rts {
		fmt.Println(rts[i])
	}

	// You can pass the server url (ts.URL) to your client implementation and send the request
	resp, err := http.DefaultClient.Get(ts.URL + "/unicorns/6198f9da97069d03e849096d")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(b))

	// Unordered output:
	// Method: DELETE, Pattern: /unicorns/{id}
	// Method: GET, Pattern: /unicorns
	// Method: GET, Pattern: /unicorns/{id}
	// Method: POST, Pattern: /unicorns
	// Method: PUT, Pattern: /unicorns/{id}
	// {"_id":"6198f9da97069d03e849096d","name":"Sparkle Angel","age":2,"colour":"blue"}
}
