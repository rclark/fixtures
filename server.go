package fixtures

import (
	"fmt"
	"net"
	"net/http"
)

// Server is an HTTP server that serves static files.
type Server struct {
	mux *http.ServeMux
}

// ServerOption is a functional option for configuring a Server.
type ServerOption func(*Server)

// WithFixture sets the Server to serve a static file at the given URL path.
func WithFixture(urlPath, filePath string) ServerOption {
	return func(s *Server) {
		s.mux.HandleFunc(urlPath, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filePath)
		})
	}
}

func allowed(handler http.Handler, methods ...string) http.Handler {
	if len(methods) == 0 {
		return handler
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, method := range methods {
			if r.Method == method {
				handler.ServeHTTP(w, r)
				return
			}
		}

		http.NotFound(w, r)
	})
}

// WithHandler allows you to provide a custom handler for a given URL path. The
// handler will only be called if the request method is one of the methods
// provided. If no methods are provided, the handler will be called for all
// requests.
func WithHandler(urlPath string, handler http.Handler, methods ...string) ServerOption {
	return func(s *Server) {
		s.mux.Handle(urlPath, allowed(handler, methods...))
	}
}

// WithHandlerFunc allows you to provide a custom handler function for a given
// URL path. The handler function will only be called if the request method is
// one of the methods provided. If no methods are provided, the handler function
// will be called for all requests.
func WithHandlerFunc(urlPath string, handler http.HandlerFunc, methods ...string) ServerOption {
	return func(s *Server) {
		s.mux.Handle(urlPath, allowed(handler, methods...))
	}
}

// NewServer creates a new Server with the given options. Use WithFixture to add
// routes on the server that will return static files.
func NewServer(opts ...ServerOption) *Server {
	s := &Server{mux: http.NewServeMux()}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type transport struct {
	Transport http.RoundTripper
	Domain    string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Host = t.Domain
	req.URL.Scheme = "http"
	return t.Transport.RoundTrip(req)
}

// ServerData contains information about a running Server.
type ServerData struct {
	Addr   *net.TCPAddr // Addr is the address of the server.
	Client *http.Client // Client is an HTTP client that will only make requests to the server.
}

// Listen starts the server and returns information about it. The caller must
// call the returned function to stop the server.
func (s Server) Listen() (ServerData, func()) {
	server := &http.Server{Handler: s.mux}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	data := ServerData{Addr: listener.Addr().(*net.TCPAddr)}
	data.Client = &http.Client{
		Transport: &transport{
			Transport: http.DefaultTransport,
			Domain:    fmt.Sprintf("%s:%d", data.Addr.IP.String(), data.Addr.Port),
		},
	}

	close := func() {
		defer listener.Close()
		defer server.Close()
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return data, close
}
