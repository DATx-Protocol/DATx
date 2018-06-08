package openssl

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
)

type CSR struct {
	//path string
	//key  string

	content    []byte
	contentKey []byte
}

func (o *Openssl) LoadCSR(filename, keyfile string) (*CSR, error) {
	var err error
	o.Init()

	filename = o.Path + "/" + filename
	keyfile = o.Path + "/" + keyfile

	c := &CSR{}

	c.content, err = ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	c.contentKey, err = ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (o *Openssl) CreateCSR(cn string, server bool) (*CSR, error) {
	var err error
	o.Init()
	c := &CSR{}
	args := []string{
		"req",
		"-days", "3650",
		"-nodes",
		"-new",
		"-keyout", "/dev/stdout",
		"-out", "/dev/stdout",
		"-config", o.GetConfigFile(),
		"-batch",
		"-utf8",
		"-subj", "/C=" + o.Country + "/ST=" + o.Province + "/L=" + o.City + "/O=" + o.Organization + "/CN=" + cn + "/emailAddress=" + o.Email,
	}

	if server {
		args = append(args, "-extensions", "server")
	}

	content, err := exec.Command("openssl", args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("openssl req: " + err.Error() + " (" + string(content) + ")")
	}

	reCert := regexp.MustCompile("(?ms)-----BEGIN CERTIFICATE REQUEST-----(.+)-----END CERTIFICATE REQUEST-----")
	reKey := regexp.MustCompile("(?ms)-----BEGIN PRIVATE KEY-----(.+)-----END PRIVATE KEY-----")

	c.content = reCert.Find(content)
	c.contentKey = reKey.Find(content)

	if len(c.content) == 0 {
		err = fmt.Errorf("Generated csr is 0 long")
		return nil, err
	}

	if len(c.contentKey) == 0 {
		err = fmt.Errorf("Generated csr key is 0 long")
		return nil, err
	}
	return c, nil
}

func (csr *CSR) Save(filename string) error {
	if err := ioutil.WriteFile(filename, csr.content, 0600); err != nil {
		return err
	}
	return nil
}

func (csr *CSR) SaveKey(filename string) error {
	if err := ioutil.WriteFile(filename, csr.contentKey, 0600); err != nil {
		return err
	}
	return nil
}

func (csr *CSR) String() string {
	if csr != nil {
		return string(csr.content)
	}
	return ""
}

func (csr *CSR) KeyString() string {
	if csr != nil {
		return string(csr.contentKey)
	}
	return ""
}
