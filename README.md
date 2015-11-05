# goji_test
Test helpers for the Goji web framework

# Request Builder

```Go
import request "github.com/vsco/goji_test/builder"

type json struct {
	Foo string `json:"foo"`
}

// Example Handler and Mux
func testHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.Method))
}

s := web.New()
s.Post("/post", testHandler)

bldr := builder.WithMux(s)

// Sending Requests. 
resp := bldr.Post("/post").Do()

// Sending POST Params
resp := bldr.Post("/post").Param("foo", "bar").Do()
resp := bldr.Post("/post").Params(map[string]string{"foo":"bar",}).Do()
	
// Sending JSON POST Bodies
js := &json{
	Foo: "bar",
}
resp := bldr.Post("/post").JSON(js).Do()
	
// Sending Headers
resp := bldr.Post("/post").Header("foo", "bar").Do()
resp := bldr.Post("/post").Headers(map[string]string{"foo":"bar",}).Do()
```

# Asserting Responses

```Go

import (
	"github.com/vsco/goji_test/builder"
	"github.com/vsco/goji_test/assert"
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
	
	req := builder.WithMux(s).Post("/post").JSON(js).Do()
	
	expected := &json{
		Foo: "bar",
	}

	assert.Response(t, req.Response).
		IsOK().
		IsJSON().
		ContainsJSON(&expected)
}
```
