package fixtures_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/rclark/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleServer_withFixture() {
	// Generate a fixture file.
	expect := "Lorem ipsum dolor sit amet"
	file, err := os.CreateTemp("", "")
	if err != nil {
		log.Fatal("creating temp file should not error: ", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()
	if _, err := file.WriteString(expect); err != nil {
		log.Fatal("writing to temp file should not error: ", err)
	}
	file.Close()

	// Create a server that serves the fixture file.
	s := fixtures.NewServer(
		fixtures.WithFixture("/data", file.Name()),
	)

	// Start the server and defer stopping it.
	info, stop := s.Listen()
	defer stop()

	// Make a request to the server.
	url := fmt.Sprintf("http://%s/data", info.Addr.String())
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		log.Fatal("request for fixture file should not error: ", err)
	}
	defer resp.Body.Close()

	found, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("reading body should not error: ", err)
	}

	fmt.Println(string(found)) // Output: Lorem ipsum dolor sit amet
}

func ExampleServer_withHandlerFunc() {
	// Create a server that serves a handler function.
	s := fixtures.NewServer(
		fixtures.WithHandlerFunc("/data", func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte("Lorem ipsum dolor sit amet")); err != nil {
				log.Fatal("writing to response should not error: ", err)
			}
		}),
	)

	// Start the server and defer stopping it.
	info, stop := s.Listen()
	defer stop()

	// Make a request to the server.
	url := fmt.Sprintf("http://%s/data", info.Addr.String())
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		log.Fatal("request for handler func should not error: ", err)
	}
	defer resp.Body.Close()

	found, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("reading body should not error: ", err)
	}

	fmt.Println(string(found)) // Output: Lorem ipsum dolor sit amet
}

func ExampleServer_withHandlerAndNotAllowedMethod() {
	// Create a server that serves a handler function. Only GET requests are
	// allowed.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Lorem ipsum dolor sit amet")); err != nil {
			log.Fatal("writing to response should not error: ", err)
		}
	})
	s := fixtures.NewServer(
		fixtures.WithHandler("/data", handler, "GET"),
	)

	// Start the server and defer stopping it.
	info, stop := s.Listen()
	defer stop()

	// Make a request to the server.
	url := fmt.Sprintf("http://%s/data", info.Addr.String())
	resp, err := http.DefaultClient.Post(url, "text/plain", nil)
	if err != nil {
		log.Fatal("request for handler func should not error: ", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode) // Output: 404
}

func ExampleServer_withHandlerAndAllowedMethod() {
	// Create a server that serves a handler function. Only GET requests are
	// allowed.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Lorem ipsum dolor sit amet")); err != nil {
			log.Fatal("writing to response should not error: ", err)
		}
	})
	s := fixtures.NewServer(
		fixtures.WithHandler("/data", handler, "GET"),
	)

	// Start the server and defer stopping it.
	info, stop := s.Listen()
	defer stop()

	// Make a request to the server.
	url := fmt.Sprintf("http://%s/data", info.Addr.String())
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		log.Fatal("request for handler func should not error: ", err)
	}
	defer resp.Body.Close()

	found, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("reading body should not error: ", err)
	}

	fmt.Println(string(found)) // Output: Lorem ipsum dolor sit amet
}

func ExampleServer_providedClient() {
	// Create a server that serves a handler function.
	s := fixtures.NewServer(
		fixtures.WithHandlerFunc("/data", func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte("Lorem ipsum dolor sit amet")); err != nil {
				log.Fatal("writing to response should not error: ", err)
			}
		}),
	)

	// Start the server and defer stopping it.
	info, stop := s.Listen()
	defer stop()

	// Make a request to the server. It doesn't matter what you pass as the host,
	// it will always be directed to the server.
	resp, err := info.Client.Get("https://just.made.this.up.com/data")
	if err != nil {
		log.Fatal("request should not error: ", err)
	}
	defer resp.Body.Close()

	found, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("reading body should not error: ", err)
	}

	fmt.Println(string(found)) // Output: Lorem ipsum dolor sit amet
}

func TestServer(t *testing.T) {
	expect := "Lorem ipsum dolor sit amet"

	// Create a server that serves a handler function.
	s := fixtures.NewServer(
		fixtures.WithHandlerFunc("/data", func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(expect))
			require.NoError(t, err, "writing to response should not error")
		}),
	)

	// Start the server and defer stopping it.
	info := s.TestListen(t)

	// Make a request to the server. It doesn't matter what you pass as the host,
	// it will always be directed to the server.
	resp, err := info.Client.Get("https://just.made.this.up.com/data")
	require.NoError(t, err, "request should not error")
	defer resp.Body.Close()

	found, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "reading body should not error")

	assert.Equal(t, expect, string(found), "response body should match expected")
}
