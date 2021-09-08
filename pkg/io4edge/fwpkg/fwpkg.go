package fwpkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// FwManifest contains the fields from the manifest file within the firmware package
type FwManifest struct {
	Name          string
	Version       string
	File          string
	Compatibility struct {
		HW        string
		MajorRevs []int `json:"major_revs"`
	}
}

// FirmwarePackage is a handle to manage firmware package archive files
type FirmwarePackage struct {
	fileName string
	file     *os.File
	manifest *FwManifest
}

// NewFirmwarePackageFromFile creates an object to work with the firmware package in fileName
// The file is opened and the manifest is parsed and checed for validity
func NewFirmwarePackageFromFile(fileName string) (*FirmwarePackage, error) {
	p := new(FirmwarePackage)
	p.fileName = fileName

	f, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("can't open " + fileName + ": " + err.Error())
	}
	p.file = f

	m, err := p.getManifest()
	if err != nil {
		return nil, err
	}
	p.manifest = m
	return p, nil
}

// Manifest returns the parsed manifest as ago struct
func (p *FirmwarePackage) Manifest() (manifest *FwManifest) {
	return p.manifest
}

func (p *FirmwarePackage) getManifest() (*FwManifest, error) {
	mJSON := new(bytes.Buffer)

	err := untarFileContent(p.file, "./manifest.json", mJSON)
	if err != nil {
		return nil, errors.New("can't untar manifest: " + err.Error())
	}
	m, err := decodeManifest(mJSON.Bytes())
	if err != nil {
		return nil, errors.New("error in manifest: " + err.Error())
	}
	return m, nil
}

func decodeManifest(b []byte) (*FwManifest, error) {
	var m *FwManifest

	if m == nil {
		fmt.Printf("m is nil")
	}

	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, errors.New("can't decode manifest: " + err.Error())
	}
	if m == nil {
		fmt.Printf("m is nil")
	}
	if m.Name == "" {
		return nil, errors.New("missing \"name\" in manifest")
	}
	if m.Version == "" {
		return nil, errors.New("missing \"version\" in manifest")
	}
	if m.File == "" {
		return nil, errors.New("missing \"file\" in manifest")
	}
	if m.Compatibility.HW == "" {
		return nil, errors.New("missing \"compatibility.hw\" in manifest")
	}
	if len(m.Compatibility.MajorRevs) == 0 {
		return nil, errors.New("missing \"compatibility.major_revs\" in manifest")
	}
	return m, nil
}
