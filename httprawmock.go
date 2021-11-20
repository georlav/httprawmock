package httprawmock

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	*httptest.Server
	router *chi.Mux
}

// NewServer starts and returns a new Server.
// The caller should call Close when finished, to shut it down.
func NewServer(routes ...Route) *Server {
	s := NewUnstartedServer(routes...)
	s.Start()

	return s
}

// NewUnstartedServer returns a new Server but doesn't start it.
//
// After changing its configuration, the caller should call Start or
// StartTLS.
//
// The caller should call Close when finished, to shut it down.
func NewUnstartedServer(routes ...Route) *Server {
	r := chi.NewRouter()
	for i := range routes {
		if routes[i].Method == "" {
			r.HandleFunc(routes[i].Pattern, createHandlerFunc(routes[i]))
		} else {
			r.MethodFunc(routes[i].Method, routes[i].Pattern, createHandlerFunc(routes[i]))
		}
	}

	return &Server{
		Server: httptest.NewUnstartedServer(r),
		router: r,
	}
}

// GetRoutes returns registered routes
func (s *Server) GetRoutes() ([]string, error) {
	if s.router == nil {
		return nil, errors.New("router is not initialized")
	}

	var routes []string

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		routes = append(routes, fmt.Sprintf("Method: %s, Pattern: %s", method, route))
		return nil
	}

	if err := chi.Walk(s.router, walkFunc); err != nil {
		return nil, errors.New("failed to walk routes")
	}

	return routes, nil

}

// SetCustomNotFoundHandler change the default not found router handler
func (s *Server) SetCustomNotFoundHandler(h http.HandlerFunc) {
	s.router.NotFound(h)
}

// SetCustomMethodNotAllowedHandler change the default method not allowed router handler
func (s *Server) SetCustomMethodNotAllowedHandler(h http.HandlerFunc) {
	s.router.MethodNotAllowed(h)
}

// SetNotFoundHandler change the default not found router handler
func (s *Server) SetNotFoundHandler(h http.HandlerFunc) {
	s.router.NotFound(h)
}

// createHandlerFunc creates a handler func using data from route object
func createHandlerFunc(rt Route) http.HandlerFunc {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		resp, err := readResponse(rt.Response, r)
		if err != nil {
			panic(fmt.Sprintf("failed to read response mock data, %s", err))
		}

		for i := range resp.Header {
			w.Header().Set(i, resp.Header.Get(i))
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(fmt.Sprintf("failed to read response mock body, %s", err))
		}

		if resp.StatusCode != http.StatusOK {
			w.WriteHeader(resp.StatusCode)
		}

		if resp.StatusCode == http.StatusNoContent {
			return
		}

		if _, err := w.Write(b); err != nil {
			panic(fmt.Sprintf("failed to write data to the connection, %s", err))
		}
	}

	return handlerFunc
}

// readResponse parses response data and converts it to a http.Response object
func readResponse(response []byte, request *http.Request) (*http.Response, error) {
	var resp []byte

	scanner := bufio.NewScanner(bufio.NewReader(bytes.NewBuffer(response)))
	for scanner.Scan() {
		if !bytes.Contains(bytes.ToLower(scanner.Bytes()), []byte("transfer-encoding")) {
			resp = append(resp, append(scanner.Bytes(), '\n')...)
		}
	}
	resp = append(resp, '\n')

	parsedResp, err := http.ReadResponse(bufio.NewReader(bytes.NewBuffer(resp)), request)
	if err != nil {
		return nil, fmt.Errorf("failed to read response data, %s", err)
	}

	return parsedResp, nil
}
