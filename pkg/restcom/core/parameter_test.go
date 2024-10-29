package core


import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitNamespaceAndParameter(t *testing.T) {
	namespace, parameter := splitNameSpaceAndParameter("namespace.parameter")
	assert.Equal(t, "namespace", namespace)
	assert.Equal(t, "parameter", parameter)

	namespace, parameter = splitNameSpaceAndParameter("parameter")
	assert.Equal(t, "", namespace)
	assert.Equal(t, "parameter", parameter)
}