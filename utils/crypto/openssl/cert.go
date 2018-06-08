package openssl

import (
	"fmt"
	"io/ioutil"
)

type Cert struct {
	path string
	key  string

	content    []byte
	contentKey []byte
}

func (o *Openssl) LoadOrCreateCert(filename, keyfile, cn string, ca *CA, server bool) (*Cert, error) {
	cert, err := o.LoadCert(filename, keyfile)
	if err != nil {
		return o.CreateCert(filename, keyfile, cn, ca, server)
	}

	return cert, nil
}

func (o *Openssl) LoadCert(filename, keyfile string) (*Cert, error) {
	var err error
	o.Init()

	filename = o.Path + "/" + filename
	keyfile = o.Path + "/" + keyfile

	c := &Cert{}
	c.path = filename
	c.key = keyfile

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

func (o *Openssl) CreateCert(filename, keyfile, cn string, ca *CA, server bool) (*Cert, error) {
	o.Init()

	filename = o.Path + "/" + filename
	keyfile = o.Path + "/" + keyfile

	request, err := o.CreateCSR(cn, server)
	if err != nil {
		return nil, fmt.Errorf("Create csr failed: " + err.Error())
	}

	cert, err := ca.Sign(request)
	if err != nil {
		return nil, fmt.Errorf("Sign csr failed: ", err)
	}

	cert.path = filename
	cert.key = keyfile

	if err = ioutil.WriteFile(cert.path, cert.content, 0600); err != nil {
		return nil, err
	}
	if err = ioutil.WriteFile(cert.key, request.contentKey, 0600); err != nil {
		return nil, err
	}

	return cert, nil
}

func (c *Cert) GetFilePath() string {
	return c.path
}

func (c *Cert) GetKeyPath() string {
	return c.key
}

func (c *Cert) String() string {
	if c != nil {
		return string(c.content)
	}
	return ""
}

func (c *Cert) KeyString() string {
	if c != nil {
		return string(c.contentKey)
	}
	return ""
}
