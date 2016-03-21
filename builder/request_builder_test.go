package builder

import (
	"net/http"
	"testing"

	"io/ioutil"
	"log"

	"github.com/vsco/http-test/assert"
	"github.com/zenazn/goji/web"
)

type jsonResponse struct {
	Foo string `json:"foo"`
}

func testHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Etag", "abcde")
	w.Write([]byte(r.Method))
}

func jsonEchoHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func server() *web.Mux {
	s := web.New()
	s.Post("/post", testHandler)
	s.Get("/get", testHandler)
	s.Put("/put", testHandler)
	s.Delete("/delete", testHandler)
	s.Head("/head", testHandler)
	s.Options("/options", testHandler)
	s.Patch("/patch", testHandler)
	s.Post("/json", jsonEchoHandler)

	return s
}

func TestMethods(t *testing.T) {
	s := server()

	bldr := WithMux(s)

	resp := bldr.Post("/post").Do()
	assert.Response(t, resp.Response).Contains("POST").IsOK()

	resp = bldr.Get("/get").Param("name", "foo").Do()
	assert.Response(t, resp.Response).Contains("GET").IsOK()

	resp = bldr.Put("/put").Do()
	assert.Response(t, resp.Response).Contains("PUT").IsOK()

	resp = bldr.Delete("/delete").Do()
	assert.Response(t, resp.Response).Contains("DELETE").IsOK()

	resp = bldr.Head("/head").Do()
	assert.Response(t, resp.Response).IsOK()

	resp = bldr.Options("/options").Do()
	assert.Response(t, resp.Response).Contains("OPTIONS").IsOK()

	resp = bldr.Patch("/patch").Do()
	assert.Response(t, resp.Response).Contains("PATCH").IsOK()
}

func TestJSON(t *testing.T) {
	s := server()

	js := &jsonResponse{
		Foo: "bar",
	}

	resp := WithMux(s).Post("/json").JSON(js).Do()

	assert.Response(t, resp.Response).ContainsJSON(js)
}

func TestJSONCases(t *testing.T) {
	s := server()
	bldr := WithMux(s)

	hsh := map[string]string{
		"foo": "bar",
		"a":   "b",
		"b":   "c",
	}
	resp := bldr.Post("/json").JSON(hsh).Do()
	assert.Response(t, resp.Response).ContainsJSON(hsh)

	ar := [3]string{"a", "b", "c"}
	resp = bldr.Post("/json").JSON(ar).Do()
	assert.Response(t, resp.Response).ContainsJSON(ar)
}

func TestHeaders(t *testing.T) {
	s := server()

	resp := WithMux(s).Get("/get").Do()

	assert.Response(t, resp.Response).ContainsHeaderValue("Etag", "abcde")
}

func TestGetRequestParams(t *testing.T) {
	s := server()

	req := WithMux(s).Get("/get").Param("name", "foo").Do().Request

	if req.URL.RawQuery != "name=foo" {
		t.Errorf("query params do not match. got %s", req.URL.RawQuery)
	}
}
