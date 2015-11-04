package response

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type TestResponse struct {
	Code       int
	Headers    http.Header
	Body       io.Reader
	BodyString string
	BodyBytes  []byte
	Response   *http.Response
}

func (r *TestResponse) UnmarshalBody(v interface{}) {
	err := json.Unmarshal(r.BodyBytes, v)

	if err != nil {
		log.Panic(err)
	}
}

func NewTestResponse(rec *http.Response) (r *TestResponse) {
	bytes, _ := ioutil.ReadAll(rec.Body)
	body := string(bytes[:])

	r = &TestResponse{
		Code:       rec.StatusCode,
		Headers:    rec.Header,
		Body:       rec.Body,
		BodyString: body,
		BodyBytes:  bytes,
		Response:   rec,
	}

	return r
}
