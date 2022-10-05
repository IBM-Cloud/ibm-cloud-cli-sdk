package rest

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var chunkSize int

func TestIsJSONStreamWithJSON(t *testing.T) {
	assert := assert.New(t)

	chunkSize = 1024
	jsonBytes := []byte(`{
		"foo": {
			"bar": 1
			"otherStuff": 2
		}
	}`)
	jsonBytesAsString := string(jsonBytes)

	r, isJSON := IsJSONStream(bytes.NewReader(jsonBytes), chunkSize)

	assert.NotNil(r)

	raw, _ := ioutil.ReadAll(r)

	assert.NotNil(raw)
	assert.True(isJSON)
	assert.Equal(jsonBytesAsString, string(raw))
}

func TestIsJSONStreamWithJSONArray(t *testing.T) {
	assert := assert.New(t)

	chunkSize = 1024
	jsonBytes := []byte(`[
		{
			"foo": {
				"bar": 1
				"otherStuff": 2
			}
		},
		{
			"foo2": {
				"bar": 1
				"otherStuff": 2
			}
		}
	]`)
	jsonBytesAsString := string(jsonBytes)

	r, isJSON := IsJSONStream(bytes.NewReader(jsonBytes), chunkSize)

	assert.NotNil(r)

	raw, _ := ioutil.ReadAll(r)

	assert.NotNil(raw)
	assert.True(isJSON)
	assert.Equal(jsonBytesAsString, string(raw))
}

func TestIsJSONStreamWithSpaceJSON(t *testing.T) {
	assert := assert.New(t)

	chunkSize = 1024
	jsonBytes := []byte("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n" +
		"\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n" +
		"\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\t\t\t\t\t\t\t\t\t" +
		"\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t" +
		"\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t" +
		"\t\t\t\t\t\t\t\t\t\t{\"foo\": \"bar\"}")
	jsonBytesAsString := string(jsonBytes)

	r, isJSON := IsJSONStream(bytes.NewReader(jsonBytes), chunkSize)

	assert.NotNil(r)

	raw, _ := ioutil.ReadAll(r)

	assert.NotNil(raw)
	assert.True(isJSON)
	assert.Equal(jsonBytesAsString, string(raw))
}

func TestIsJSONStreamWithInvalidJSON(t *testing.T) {
	assert := assert.New(t)

	chunkSize = 1024
	jsonBytes := []byte(`
		"foo": {
			"bar": 1
			"otherStuff": 2
		}
	`)
	jsonBytesAsString := string(jsonBytes)

	r, isJSON := IsJSONStream(bytes.NewReader(jsonBytes), chunkSize)

	assert.NotNil(r)

	raw, _ := ioutil.ReadAll(r)

	assert.NotNil(raw)
	assert.False(isJSON)
	assert.Equal(jsonBytesAsString, string(raw))
}

func TestIsJSONStreamWithString(t *testing.T) {
	assert := assert.New(t)

	chunkSize = 1024
	jsonBytes := []byte("foobar")
	jsonBytesAsString := string(jsonBytes)

	r, isJSON := IsJSONStream(bytes.NewReader(jsonBytes), chunkSize)

	assert.NotNil(r)

	raw, _ := ioutil.ReadAll(r)

	assert.NotNil(raw)
	assert.False(isJSON)
	assert.Equal(jsonBytesAsString, string(raw))
}

func TestIsJSONStreamWithYAML(t *testing.T) {
	assert := assert.New(t)

	chunkSize = 1024
	jsonBytes := []byte(`---
	foo:
	- bar: jon_snow
	  description: "The king of the north"
	  created: 2016-01-14
	  updated: 2019-08-16
	  company: IBM
	`)
	jsonBytesAsString := string(jsonBytes)

	r, isJSON := IsJSONStream(bytes.NewReader(jsonBytes), chunkSize)

	assert.NotNil(r)

	raw, _ := ioutil.ReadAll(r)

	assert.NotNil(raw)
	assert.False(isJSON)
	assert.Equal(jsonBytesAsString, string(raw))
}

func TestIsJSONStreamWithYamlArray(t *testing.T) {
	assert := assert.New(t)

	chunkSize = 1024
	jsonBytes := []byte(`
	items:
	  - things:
		  thing1: huey
		  things2: dewey
		  thing3: louie
	  - other things:
		  key: value
	`)
	jsonBytesAsString := string(jsonBytes)

	r, isJSON := IsJSONStream(bytes.NewReader(jsonBytes), chunkSize)

	assert.NotNil(r)

	raw, _ := ioutil.ReadAll(r)

	assert.NotNil(raw)
	assert.False(isJSON)
	assert.Equal(jsonBytesAsString, string(raw))
}

func TestIsJSONStreamWithEmptyString(t *testing.T) {
	assert := assert.New(t)

	chunkSize = 1024
	jsonBytes := []byte("")
	jsonBytesAsString := string(jsonBytes)

	r, isJSON := IsJSONStream(bytes.NewReader(jsonBytes), chunkSize)

	assert.NotNil(r)

	raw, _ := ioutil.ReadAll(r)

	assert.NotNil(raw)
	assert.False(isJSON)
	assert.Equal(jsonBytesAsString, string(raw))
}
