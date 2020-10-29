package bluemix_test

import (
	"testing"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	v := VersionType{}
	assert.Equal(t, 0, v.Major)
	assert.Equal(t, 0, v.Minor)
	assert.Equal(t, 0, v.Build)
	assert.Empty(t, v.String())

	v = VersionType{Major: 1, Minor: 2, Build: 3}
	assert.Equal(t, "1.2.3", v.String())
}
