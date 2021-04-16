package rest

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var noContentHandler = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func serveHandler(statusCode int, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		fmt.Fprint(w, message)
	}
}

func TestDoWithContext_WithSuccessV(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(serveHandler(200, "{\"foo\": \"bar\"}"))
	defer ts.Close()

	var res map[string]string
	resp, err := NewClient().DoWithContext(context.TODO(), GetRequest(ts.URL), &res, nil)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal("bar", res["foo"])
}

func TestDoWithContext_WithSuccessYAML(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(serveHandler(200, "\"foo\": \"bar\""))
	defer ts.Close()

	var res map[string]string
	resp, err := NewClient().DoWithContext(context.TODO(), GetRequest(ts.URL), &res, nil)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal("bar", res["foo"])
}

func TestDoWithContext_WithSuccessComplexYAML(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(serveHandler(200, `
---
doe: "a deer, a female deer"
ray: "a drop of golden sun"
pi: 3.14159
xmas: true
french-hens: 3
calling-birds:
  - huey
  - dewey
  - louie
  - fred`))
	defer ts.Close()

	var res map[string]interface{}
	resp, err := NewClient().DoWithContext(context.TODO(), GetRequest(ts.URL), &res, nil)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal("a deer, a female deer", res["doe"])
	assert.Equal("a drop of golden sun", res["ray"])
	assert.Equal(3.14159, res["pi"])
	assert.Equal(true, res["xmas"])
	callingBirds := []interface{}{"huey", "dewey", "louie", "fred"}
	assert.Equal(callingBirds, res["calling-birds"])
}

func TestDoWithContext_NilContext_WithSuccessV(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(serveHandler(200, "{\"foo\": \"bar\"}"))
	defer ts.Close()

	var res map[string]string
	resp, err := NewClient().DoWithContext(nil, GetRequest(ts.URL), &res, nil)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal("bar", res["foo"])
}

func TestDoWithContext_NilContext_WithSuccessYAML(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(serveHandler(200, "\"foo\": \"bar\""))
	defer ts.Close()

	var res map[string]string
	resp, err := NewClient().DoWithContext(nil, GetRequest(ts.URL), &res, nil)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal("bar", res["foo"])
}

func TestDoWithContext_WithoutSuccessV(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(serveHandler(200, "abcedefg"))
	defer ts.Close()

	resp, err := NewClient().Do(GetRequest(ts.URL), nil, nil)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestDo_WithSuccessV(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(serveHandler(200, "{\"foo\": \"bar\"}"))
	defer ts.Close()

	var res map[string]string
	resp, err := NewClient().Do(GetRequest(ts.URL), &res, nil)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal("bar", res["foo"])
}

func TestDo_WithSuccessYAML(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(serveHandler(200, "\"foo\": \"bar\""))
	defer ts.Close()

	var res map[string]string
	resp, err := NewClient().Do(GetRequest(ts.URL), &res, nil)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal("bar", res["foo"])
}

func TestDo_WithoutSuccessV(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(serveHandler(200, "abcedefg"))
	defer ts.Close()

	resp, err := NewClient().Do(GetRequest(ts.URL), nil, nil)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestDo_ServerError_WithoutErrorV(t *testing.T) {
	assert := assert.New(t)

	code := 500
	errResp := "Internal server error."

	ts := httptest.NewServer(serveHandler(code, errResp))
	defer ts.Close()

	var successV interface{}
	_, err := NewClient().Do(GetRequest(ts.URL), &successV, nil)
	assert.Nil(successV)
	assert.Error(err)
	assert.Equal(err, &ErrorResponse{code, errResp})
}

func TestDo_ServerError_WithErrorV(t *testing.T) {
	assert := assert.New(t)

	code := 500
	errResp := "{\"message\": \"Internal server error.\"}"

	ts := httptest.NewServer(serveHandler(code, errResp))
	defer ts.Close()

	var successV interface{}
	var errorV = struct {
		Message string
	}{}

	_, err := NewClient().Do(GetRequest(ts.URL), &successV, &errorV)
	assert.NoError(err)
	assert.Equal(errorV.Message, "Internal server error.")
}

func TestNoContent(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(noContentHandler))
	defer ts.Close()

	req := GetRequest(ts.URL)

	var successV interface{}
	_, err := NewClient().Do(req, &successV, nil)
	assert.Error(err) // empty response body error

	_, err = NewClient().Do(req, nil, nil)
	assert.NoError(err)
}

func TestDownloadFile(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(serveHandler(200, "abcedefg"))
	defer ts.Close()

	f, err := ioutil.TempFile("", "BluemixCliRestTest")
	assert.NoError(err)
	defer f.Close()

	_, err = NewClient().Do(GetRequest(ts.URL), f, nil)
	assert.NoError(err)
	bytes, _ := ioutil.ReadFile(f.Name())
	assert.Equal("abcedefg", string(bytes))
}
