[![Go](https://github.com/rclark/fixtures/actions/workflows/go.yml/badge.svg)](https://github.com/rclark/fixtures/actions/workflows/go.yml)

# fixtures

A set of utilities for working with test fixtures in Go.

```go
import "github.com/rclark/fixtures"
```

## Index

- [type Server](<#Server>)
  - [func NewServer\(opts ...ServerOption\) \*Server](<#NewServer>)
  - [func \(s Server\) Listen\(\) \(ServerData, func\(\)\)](<#Server.Listen>)
  - [func \(s Server\) TestListen\(t \*testing.T\) ServerData](<#Server.TestListen>)
- [type ServerData](<#ServerData>)
- [type ServerOption](<#ServerOption>)
  - [func WithFixture\(urlPath, filePath string\) ServerOption](<#WithFixture>)
  - [func WithHandler\(urlPath string, handler http.Handler, methods ...string\) ServerOption](<#WithHandler>)
  - [func WithHandlerFunc\(urlPath string, handler http.HandlerFunc, methods ...string\) ServerOption](<#WithHandlerFunc>)


<a name="Server"></a>
## type Server

Server is an HTTP server that serves static files.

```go
type Server struct {
    // contains filtered or unexported fields
}
```

<details><summary>Example (Provided Client)</summary>
<p>



```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rclark/fixtures"
)

func main() {
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

	fmt.Println(string(found))
}
```

#### Output

```
Lorem ipsum dolor sit amet
```

</p>
</details>

<details><summary>Example (With Fixture)</summary>
<p>



```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/rclark/fixtures"
)

func main() {
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

	fmt.Println(string(found))
}
```

#### Output

```
Lorem ipsum dolor sit amet
```

</p>
</details>

<details><summary>Example (With Handler And Allowed Method)</summary>
<p>



```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rclark/fixtures"
)

func main() {
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

	fmt.Println(string(found))
}
```

#### Output

```
Lorem ipsum dolor sit amet
```

</p>
</details>

<details><summary>Example (With Handler And Not Allowed Method)</summary>
<p>



```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rclark/fixtures"
)

func main() {
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

	fmt.Println(resp.StatusCode)
}
```

#### Output

```
404
```

</p>
</details>

<details><summary>Example (With Handler Func)</summary>
<p>



```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rclark/fixtures"
)

func main() {
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

	fmt.Println(string(found))
}
```

#### Output

```
Lorem ipsum dolor sit amet
```

</p>
</details>

<a name="NewServer"></a>
### func NewServer

```go
func NewServer(opts ...ServerOption) *Server
```

NewServer creates a new Server with the given options. Use WithFixture to add routes on the server that will return static files.

<a name="Server.Listen"></a>
### func \(Server\) Listen

```go
func (s Server) Listen() (ServerData, func())
```

Listen starts the server and returns information about it. The caller must call the returned function to stop the server. If it can't find a port to listen on, or if the server closes unexpectedly, it will panic.

<a name="Server.TestListen"></a>
### func \(Server\) TestListen

```go
func (s Server) TestListen(t *testing.T) ServerData
```

TestListen starts the server and returns information about it. The server will be closed when the test ends. If it can't find a port to listen on, or if the server closes unexpectedly, it will fail the test.

<a name="ServerData"></a>
## type ServerData

ServerData contains information about a running Server.

```go
type ServerData struct {
    Addr   *net.TCPAddr // Addr is the address of the server.
    Client *http.Client // Client is an HTTP client that will only make requests to the server.
}
```

<a name="ServerOption"></a>
## type ServerOption

ServerOption is a functional option for configuring a Server.

```go
type ServerOption func(*Server)
```

<a name="WithFixture"></a>
### func WithFixture

```go
func WithFixture(urlPath, filePath string) ServerOption
```

WithFixture sets the Server to serve a static file at the given URL path.

<a name="WithHandler"></a>
### func WithHandler

```go
func WithHandler(urlPath string, handler http.Handler, methods ...string) ServerOption
```

WithHandler allows you to provide a custom handler for a given URL path. The handler will only be called if the request method is one of the methods provided. If no methods are provided, the handler will be called for all requests.

<a name="WithHandlerFunc"></a>
### func WithHandlerFunc

```go
func WithHandlerFunc(urlPath string, handler http.HandlerFunc, methods ...string) ServerOption
```

WithHandlerFunc allows you to provide a custom handler function for a given URL path. The handler function will only be called if the request method is one of the methods provided. If no methods are provided, the handler function will be called for all requests.

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
