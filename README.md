# http_test
Test helpers for the net/http requests.

# Request Builder

```Go
import request "github.com/vsco/http_test/builder"

type json struct {
	Foo string `json:"foo"`
}

// Example Handler and Mux
func testHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.Method))
}

s := web.New()
s.Post("/post", testHandler)

// Sending Requests. 
resp := Post("/post").Use(s).Do()

// Sending POST Params
resp := Post("/post").Use(s).Params(map[string]string{"foo":"bar",}).Do()
	
// Sending JSON POST Bodies
js := &json{
	Foo: "bar",
}
resp := Post("/post").Use(s).JSON(js).Do()
	
// Sending Headers
resp := Post("/post").Use(s).Header("foo", "bar").Do()
resp := Post("/post").Use(s).Headers(map[string]string{"foo":"bar",}).Do()
```

# Asserting Responses

```Go

import (
	"github.com/vsco/http_test/builder"
	"github.com/vsco/http_test/assert"
)

func jsonEchoHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

type json struct {
	Foo string `json:"foo"`
}

func TestResponse(t *testing.T) {
	s := web.New()
	s.Post("/post", jsonEchoHandler)
	
	js := &json{
		Foo: "bar",
	}
	
	req := Post("/post").Use(s).JSON(js).Do()
	
	expected := &json{
		Foo: "bar",
	}

	assert.Response(t, req.Response).
		IsOK().
		IsJSON().
		ContainsJSON(&expected)
}
```
