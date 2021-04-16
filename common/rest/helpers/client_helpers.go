package rest

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

var jsonPrefix = []byte("{")
var jsonArrayPrefix = []byte("[")

// IsJSONStream scans the provided reader up to size, looking
// for a json prefix indicating this is JSON. It will return the
// bufio.Reader it creates for the consumer. The buffer size will at either be the size of 'size'
// or the size of the Reader passed in 'r', whichever is larger.
// Inspired from https://github.com/kubernetes/apimachinery/blob/v0.21.0/pkg/util/yaml/decoder.go
func IsJSONStream(r io.Reader, size int) (io.Reader, bool) {
	buffer := bufio.NewReaderSize(r, size)
	b, _ := buffer.Peek(size)
	return buffer, containsJSON(b)
}

// containsJSON returns true if the provided buffer appears to start with
// a JSON open brace or JSON open bracket.
// Inspired from https://github.com/kubernetes/apimachinery/blob/v0.21.0/pkg/util/yaml/decoder.go
func containsJSON(buf []byte) bool {
	return hasPrefix(buf, jsonPrefix) || hasPrefix(buf, jsonArrayPrefix)
}

// Return true if the first non-whitespace bytes in buf is
// prefix.
func hasPrefix(buf []byte, prefix []byte) bool {
	trim := bytes.TrimLeftFunc(buf, unicode.IsSpace)
	return bytes.HasPrefix(trim, prefix)
}
