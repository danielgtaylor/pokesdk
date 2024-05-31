# Pokemon API SDK for Go

[![CI](https://github.com/danielgtaylor/pokesdk/workflows/CI/badge.svg?branch=main)](https://github.com/danielgtaylor/pokesdk/actions?query=workflow%3ACI+branch%3Amain++) [![codecov](https://codecov.io/gh/danielgtaylor/pokesdk/branch/main/graph/badge.svg)](https://codecov.io/gh/danielgtaylor/pokesdk) [![Docs](https://godoc.org/github.com/danielgtaylor/pokesdk?status.svg)](https://pkg.go.dev/github.com/danielgtaylor/pokesdk?tab=doc) [![Go Report Card](https://goreportcard.com/badge/github.com/danielgtaylor/pokesdk)](https://goreportcard.com/report/github.com/danielgtaylor/pokesdk)

This is a Go SDK for the [Pokémon API](https://pokeapi.co/).

## Installation

```bash
go get github.com/danielgtaylor/pokesdk
```

## Usage Example

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/danielgtaylor/pokesdk"
)

func main() {
	ctx := context.Background()
	sdk := pokesdk.New(pokesdk.Config{})

	// Print pokemon names.
	for result := range sdk.ListPokemon().All(ctx) {
		if result.Error != nil {
			log.Fatalf("Failed to list Pokemon: %v", result.Error)
		}
		fmt.Printf("Pokemon: %s\n", result.Value.Name)
	}

	// Print Pikachu's stat details.
	pika, err := sdk.GetPokemon(ctx, "pikachu")
	if err != nil {
		panic(err)
	}
	fmt.Println("Pikachu stats:")
	fmt.Println(pika.Stats)
}
```

## Development

The project has no dependencies and running the tests is easy:

```sh
$ go test -cover
```

To run a simple integration test, you can use the demo app:

```sh
go run ./cmd/demo
```

### Design

The SDK is designed to be simple and easy to use. All calls return parsed Go structs.

#### Pagination

The SDK uses a paginator pattern to allow for easy iteration over large sets of data by transparently fetching pages and providing items through a Go channel. Each result item contains the page the item came from, the overall index of the item among all pages, the value itself, and any error that occurred.

```go
for result := range sdk.ListPokemon().All(ctx) {
	if result.Error != nil {
		log.Fatalf("Failed to list Pokemon: %v", result.Error)
	}
	fmt.Printf("Pokemon: %s\n", result.Value.Name)
}
```

It's also possible to stop iteration early by using the `AllWithCancel` method and calling the cancel function so the paginator stops processing pages. For example, to stop after 50 items:

```go
iter, cancel := sdk.ListPokemon().AllWithCancel(ctx)
for result := range iter {
	if result.Index > 50 {
		// We are done, so let the paginator know to stop.
		cancel()
		break
	}
	if result.Error != nil {
		log.Fatalf("Failed to list Pokemon: %v", result.Error)
	}
	fmt.Printf("Pokemon: %s\n", result.Value.Name)
}
```

#### Extensible Client

The SDK is also designed to be easily extensible and to allow for easy mocking of the API by using custom clients/transports.

```go
transport := MyCustomTransport{}
client := &http.Client{Transport: transport}
sdk := pokesdk.New(pokesdk.Config{Client: client})
```

Ideas for extending the client:

- Auth
- Caching
- Client-side limiting of concurrent requests
- Adding tracing information to outgoing requests
- And more!

## License

This project is licensed under the MIT License.
