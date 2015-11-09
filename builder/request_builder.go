package builder

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"io"

	"bytes"

	"encoding/json"
	"log"

	"github.com/vsco/goji-test/response"
	"github.com/zenazn/goji/web"
)

type RequestBuilder struct {
	method   string
	params   map[string]string
	headers  map[string]string
	path     string
	body     io.Reader
	context  web.C
	mux      *web.Mux
	Request  *http.Request
	Response *response.TestResponse
}

func (r *RequestBuilder) Post(path string) *RequestBuilder {
	r.method = "POST"
	r.path = path

	return r
}

func (r *RequestBuilder) Get(path string) *RequestBuilder {
	r.method = "GET"
	r.path = path

	return r
}

func (r *RequestBuilder) Put(path string) *RequestBuilder {
	r.method = "PUT"
	r.path = path

	return r
}

func (r *RequestBuilder) Delete(path string) *RequestBuilder {
	r.method = "DELETE"
	r.path = path

	return r
}

func (r *RequestBuilder) Head(path string) *RequestBuilder {
	r.method = "HEAD"
	r.path = path

	return r
}

func (r *RequestBuilder) Options(path string) *RequestBuilder {
	r.method = "OPTIONS"
	r.path = path

	return r
}

func (r *RequestBuilder) Patch(path string) *RequestBuilder {
	r.method = "PATCH"
	r.path = path

	return r
}

func WithMux(m *web.Mux) (r *RequestBuilder) {
	r = &RequestBuilder{
		path:    "/",
		method:  "GET",
		mux:     m,
		context: web.C{Env: map[interface{}]interface{}{}},
	}

	return
}

func (r *RequestBuilder) Method(method string) *RequestBuilder {
	r.method = method

	return r
}

func (r *RequestBuilder) Path(path string) *RequestBuilder {
	r.path = path

	return r
}

func (r *RequestBuilder) Params(params map[string]string) *RequestBuilder {
	r.params = params

	return r
}

func (r *RequestBuilder) Param(k string, v string) *RequestBuilder {
	if r.params == nil {
		r.params = make(map[string]string)
	}

	r.params[k] = v

	return r
}

func (r *RequestBuilder) JSON(s interface{}) *RequestBuilder {
	js, err := json.Marshal(s)

	if err != nil {
		log.Fatal(err)
	}

	r.body = bytes.NewReader(js)

	return r
}

func (r *RequestBuilder) Headers(headers map[string]string) *RequestBuilder {
	r.headers = headers

	return r
}

func (r *RequestBuilder) Header(k string, v string) *RequestBuilder {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}

	r.headers[k] = v

	return r
}

func (r *RequestBuilder) Context(c web.C) *RequestBuilder {
	r.context = c

	return r
}

func (r *RequestBuilder) Do() *RequestBuilder {
	ts := httptest.NewServer(r.mux)
	defer ts.Close()

	var err error

	form := url.Values{}
	for k, v := range r.params {
		form.Add(k, v)
	}

	if r.method == "POST" || r.method == "PUT" || r.method == "PATCH" {
		if r.body == nil {
			buf := bytes.NewBufferString(form.Encode())
			r.Request, err = http.NewRequest(r.method, ts.URL+r.path, buf)
			r.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if err != nil {
				panic(err)
			}
		} else {
			r.Request, err = http.NewRequest(r.method, ts.URL+r.path, r.body)
			r.Request.Header.Set("Content-Type", "application/json")
		}
	} else {
		r.Request, err = http.NewRequest(r.method, ts.URL+r.path+"?"+form.Encode(), nil)
		if err != nil {
			panic(err)
		}
	}

	for k, v := range r.headers {
		r.Request.Header.Set(k, v)
	}

	res, err := http.DefaultClient.Do(r.Request)

	if err != nil {
		panic(err)
	}

	r.Response = response.NewTestResponse(res)

	return r
}
