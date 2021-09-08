package fwpkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeManifestOk(t *testing.T) {
	s := `
	{
		"name": "cpu01-tty_accdl",
		"version": "1.1.0",
		"file": "fw_varA_1_1_0.json",
		"compatibility": {
		  	"hw": "s101-cpu01",
		  	"major_revs": [
				1,
				2
		  	]
		}
	}
	`
	m, err := decodeManifest([]byte(s))
	assert.Nil(t, err)
	assert.Equal(t, m.Name, "cpu01-tty_accdl")
	assert.Equal(t, m.Compatibility.HW, "s101-cpu01")
	assert.Equal(t, m.Compatibility.MajorRevs, []int{1, 2})
}

func TestDecodeManifestBrokenJSON(t *testing.T) {
	s := `
	{
		"name": "cpu01-tty_accdl",
		"version": "1.1.0",
		"file": "fw_varA_1_1_0.json",
		"compatibility": {
		  	"hw": "s101-cpu01",
		  	"major_revs": [
				1,
				2
	`

	_, err := decodeManifest([]byte(s))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "can't decode")
}

func TestDecodeManifestMissingName(t *testing.T) {
	s := `
	{
		"version": "1.1.0",
		"file": "fw_varA_1_1_0.json",
		"compatibility": {
		  	"hw": "s101-cpu01",
		  	"major_revs": [
				1,
				2
		  	]
		}
	}
	`

	_, err := decodeManifest([]byte(s))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing \"name\"")
}

func TestDecodeManifestMissingRevs(t *testing.T) {
	s := `
	{
		"name": "cpu01-tty_accdl",
		"version": "1.1.0",
		"file": "fw_varA_1_1_0.json",
		"compatibility": {
		  	"hw": "s101-cpu01"
		}
	}
	`

	_, err := decodeManifest([]byte(s))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing \"compatibility.major_revs\"")
}
