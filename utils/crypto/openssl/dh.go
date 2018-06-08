package openssl

import (
	"io/ioutil"
	"os/exec"
	"strconv"
)

type DH struct {
	path    string
	content []byte
}

func (o *Openssl) LoadOrCreateDH(filename string, size int) (*DH, error) {
	dh, err := o.LoadDH(filename)
	if err != nil {
		return o.CreateDH(filename, size)
	}

	return dh, nil
}

func (o *Openssl) LoadDH(filename string) (*DH, error) {
	filename = o.Path + "/" + filename
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	dh := &DH{}
	dh.path = filename
	dh.content = content
	return dh, nil
}

func (o *Openssl) CreateDH(filename string, size int) (*DH, error) {
	var err error
	filename = o.Path + "/" + filename
	dh := &DH{}
	dh.path = filename

	dh.content, err = exec.Command("openssl", "dhparam", strconv.Itoa(size)).Output()
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(filename, dh.content, 0600)

	return dh, err
}

func (dh *DH) GetFilePath() string {
	return dh.path
}

func (dh *DH) String() string {
	return string(dh.content)
}
