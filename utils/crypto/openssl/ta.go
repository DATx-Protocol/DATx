package openssl

import (
	"io/ioutil"
	"os/exec"
)

type TA struct {
	path    string
	content []byte
}

func (o *Openssl) LoadOrCreateTA(filename string) (*TA, error) {
	ta, err := o.LoadTA(filename)
	if err != nil {
		return o.CreateTA(filename)
	}

	return ta, nil
}

func (o *Openssl) LoadTA(filename string) (*TA, error) {
	filename = o.Path + "/" + filename
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	ta := &TA{}
	ta.path = filename
	ta.content = content
	return ta, nil
}

func (o *Openssl) CreateTA(filename string) (*TA, error) {
	var err error
	filename = o.Path + "/" + filename
	ta := &TA{}
	ta.path = filename

	ta.content, err = exec.Command("openvpn", "--genkey", "--secret", "/dev/stdout").Output()
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(filename, ta.content, 0600)

	return ta, err
}

func (ta *TA) GetFilePath() string {
	return ta.path
}

func (ta *TA) String() string {
	return string(ta.content)
}
