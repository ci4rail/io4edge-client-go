package fwpkg

import (
	"bytes"
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
	assert.Equal(t, "cpu01-tty_accdl", m.Name)
	assert.Equal(t, "s101-cpu01", m.Compatibility.HW)
	assert.Equal(t, []int{1, 2}, m.Compatibility.MajorRevs)
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

func TestNewFirmwarePackageFromFile(t *testing.T) {
	p, err := NewFirmwarePackageFromFile("testdata/t1.fwpkg")
	assert.Nil(t, err)
	m := p.Manifest()
	assert.Equal(t, "cpu01-tty_accdl", m.Name)

	fwFile := new(bytes.Buffer)
	err = p.File(fwFile)
	assert.Nil(t, err)
	assert.Equal(t, 155688, fwFile.Len())
}
