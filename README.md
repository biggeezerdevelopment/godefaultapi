# Go Default API Client

A simple and flexible Go HTTP client library that provides default functionality for making API requests with support for various authentication methods.

## Features

- Simple and intuitive API
- Support for all major HTTP methods (GET, POST, PUT, DELETE, PATCH)
- Built-in support for Bearer token authentication
- Built-in support for Basic authentication
- Custom header support
- JSON and XML request/response handling
- Context support for request cancellation
- Error handling with detailed error messages
- Configurable HTTP client with timeout

## Installation

```bash
go get github.com/ddfelts/godefaultapi
```

## Usage

### Basic Usage with JSON

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ddfelts/godefaultapi"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Create a new client
	client := godefaultapi.NewClient("https://api.example.com")

	// Set authentication
	client.SetBearerToken("your-token-here")
	// OR
	client.SetBasicAuth("username", "password")

	// Make a GET request
	var user User
	err := client.Get(context.Background(), "/users/1", &user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User: %+v\n", user)

	// Make a POST request
	newUser := User{Name: "John Doe"}
	var createdUser User
	err = client.Post(context.Background(), "/users", newUser, &createdUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created user: %+v\n", createdUser)
}
```

### Using XML

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ddfelts/godefaultapi"
)

type User struct {
	ID   int    `xml:"id"`
	Name string `xml:"name"`
}

func main() {
	// Create a new client
	client := godefaultapi.NewClient("https://api.example.com")

	// Set content type to XML
	client.SetContentType(godefaultapi.ContentTypeXML)

	// Make a GET request with XML response
	var user User
	err := client.Get(context.Background(), "/users/1", &user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User: %+v\n", user)

	// Make a POST request with XML body
	newUser := User{Name: "John Doe"}
	var createdUser User
	err = client.Post(context.Background(), "/users", newUser, &createdUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created user: %+v\n", createdUser)
}
```

### Setting Custom Headers

```go
client := godefaultapi.NewClient("https://api.example.com")
client.SetHeader("X-Custom-Header", "value")
```

### Using Context

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

var result interface{}
err := client.Get(ctx, "/endpoint", &result)
if err != nil {
    log.Fatal(err)
}
```

### Using Empty Interface for Dynamic Responses

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ddfelts/godefaultapi"
)

func main() {
	client := godefaultapi.NewClient("https://api.example.com")

	// Using empty interface to handle unknown response structure
	var response interface{}
	err := client.Get(context.Background(), "/dynamic-endpoint", &response)
	if err != nil {
		log.Fatal(err)
	}

	// Type assertion to access the response data
	switch data := response.(type) {
	case map[string]interface{}:
		// Handle JSON object
		if name, ok := data["name"].(string); ok {
			fmt.Printf("Name: %s\n", name)
		}
		if age, ok := data["age"].(float64); ok {
			fmt.Printf("Age: %d\n", int(age))
		}
	case []interface{}:
		// Handle JSON array
		fmt.Printf("Received array with %d elements\n", len(data))
		for i, item := range data {
			fmt.Printf("Item %d: %v\n", i, item)
		}
	default:
		fmt.Printf("Unexpected response type: %T\n", data)
	}

	// Alternatively, you can marshal the response back to JSON for inspection
	jsonData, _ := json.MarshalIndent(response, "", "  ")
	fmt.Printf("Raw response: %s\n", string(jsonData))
}
```

### Using Empty Interface with XML Responses

```go
package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"

	"github.com/ddfelts/godefaultapi"
)

func main() {
	client := godefaultapi.NewClient("https://api.example.com")
	client.SetContentType(godefaultapi.ContentTypeXML)

	// Using empty interface to handle unknown XML response structure
	var response interface{}
	err := client.Get(context.Background(), "/dynamic-xml-endpoint", &response)
	if err != nil {
		log.Fatal(err)
	}

	// Type assertion to access the XML data
	switch data := response.(type) {
	case map[string]interface{}:
		// Handle XML object
		if name, ok := data["name"].(string); ok {
			fmt.Printf("Name: %s\n", name)
		}
		if age, ok := data["age"].(string); ok {
			fmt.Printf("Age: %s\n", age)
		}
		// Handle nested XML elements
		if address, ok := data["address"].(map[string]interface{}); ok {
			if street, ok := address["street"].(string); ok {
				fmt.Printf("Street: %s\n", street)
			}
		}
	case []interface{}:
		// Handle XML array
		fmt.Printf("Received XML array with %d elements\n", len(data))
		for i, item := range data {
			if itemMap, ok := item.(map[string]interface{}); ok {
				fmt.Printf("Item %d: %v\n", i, itemMap)
			}
		}
	default:
		fmt.Printf("Unexpected XML response type: %T\n", data)
	}

	// To inspect the raw XML response
	xmlData, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling XML: %v", err)
	} else {
		fmt.Printf("Raw XML response:\n%s\n", string(xmlData))
	}

	// Example of handling XML with known structure but unknown fields
	type DynamicXML struct {
		XMLName xml.Name
		Attrs   []xml.Attr `xml:",any,attr"`
		Content []byte     `xml:",innerxml"`
		Children []DynamicXML `xml:",any"`
	}

	var dynamicXML DynamicXML
	err = client.Get(context.Background(), "/xml-endpoint", &dynamicXML)
	if err != nil {
		log.Fatal(err)
	}

	// Process the dynamic XML structure
	fmt.Printf("Root element: %s\n", dynamicXML.XMLName.Local)
	fmt.Printf("Attributes: %v\n", dynamicXML.Attrs)
	fmt.Printf("Content: %s\n", string(dynamicXML.Content))
	for _, child := range dynamicXML.Children {
		fmt.Printf("Child element: %s\n", child.XMLName.Local)
	}
}
```

## Error Handling

The library returns detailed error messages that include:
- HTTP status codes for failed requests
- Response body for failed requests
- JSON/XML marshaling/unmarshaling errors
- Network errors

## License

MIT 