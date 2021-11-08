package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTxtValueFromKey(t *testing.T) {
	txt := [][]byte{[]byte("func=abc"), []byte("delta=def")}

	v, err := getTxtValueFromKey("func", txt)
	assert.Nil(t, err)
	assert.Equal(t, "abc", v)

	v, err = getTxtValueFromKey("delta", txt)
	assert.Nil(t, err)
	assert.Equal(t, "def", v)

	_, err = getTxtValueFromKey("gamma", txt)
	assert.NotNil(t, err)
}

func TestAuxPort(t *testing.T) {
	protocol, port, err := auxPort("tcp-10000")
	assert.Nil(t, err)
	assert.Equal(t, protocol, "tcp")
	assert.Equal(t, port, 10000)

	_, _, err = auxPort("10000")
	assert.NotNil(t, err)

	_, _, err = auxPort("tcp-1A000")
	assert.NotNil(t, err)

	_, _, err = auxPort("not_avail")
	assert.NotNil(t, err)

}
