package assert

import (
	"bytes"
	"net/http"
	"testing"

	"fmt"

	"encoding/json"

	"log"

	"github.com/stretchr/testify/assert"
	"github.com/vsco/goji-test/response"
)

type assertionGroup struct {
	response *response.TestResponse
	failures []string
	T        *testing.T
}

func Response(t *testing.T, r *response.TestResponse) (a *assertionGroup) {
	a = &assertionGroup{
		response: r,
		T:        t,
	}

	return
}

func (a *assertionGroup) IsOK() *assertionGroup {
	if a.response.Code != http.StatusOK {
		a.Errorf("status code was %d expected %d", a.response.Code, http.StatusOK)
	}

	return a
}

func (a *assertionGroup) IsJSON() *assertionGroup {
	a.ContainsContentType("application/json")

	return a
}

func (a *assertionGroup) HasStatusCode(code int) *assertionGroup {
	if a.response.Code != code {
		a.Errorf("status code was %d expected %d", a.response.Code, code)
	}

	return a
}

func (a *assertionGroup) ContainsContentType(et string) *assertionGroup {
	a.ContainsHeaderValue("Content-Type", et)

	return a
}

func (a *assertionGroup) ContainsEtag(et string) *assertionGroup {
	a.ContainsHeaderValue("Etag", et)

	return a
}

func (a *assertionGroup) ContainsHeaderValue(h string, v string) *assertionGroup {
	av := a.response.Response.Header.Get(h)

	if av == "" {
		a.Errorf("%s header is not found", h)
	}

	if av != v {
		a.Errorf("%s was not equal to %s", h, v)
	}

	return a
}

func (a *assertionGroup) Contains(b string) *assertionGroup {
	if a.response.BodyString != b {
		msg := fmt.Sprintf(`body does not match:
			%s
		expected:
			%s`, a.response.BodyString, b)
		a.Errorf(msg)
	}

	return a
}

func (a *assertionGroup) ContainsJSON(s interface{}) *assertionGroup {
	js, err := json.Marshal(s)

	if err != nil {
		log.Fatal(err)
	}

	var expected bytes.Buffer
	err = json.Indent(&expected, js, "", "\t")

	if err != nil {
		log.Fatal(err)
	}

	var actual bytes.Buffer
	err = json.Indent(&actual, a.response.BodyBytes, "", "\t")

	if err != nil {
		log.Fatal(err)
	}

	if string(actual.Bytes()[:]) != string(expected.Bytes()[:]) {
		msg := fmt.Sprintf(`JSON body does not match:
%s
expected:
%s`, string(actual.Bytes()[:]), string(expected.Bytes()[:]))
		a.Errorf(msg)
	}

	return a
}

func (a *assertionGroup) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	assert.Fail(a.T, msg)
}
