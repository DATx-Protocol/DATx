package openssl

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

type CA struct {
	path   string
	crl    string
	key    string
	config string
	root   string

	content    []byte
	contentKey []byte
}

func (o *Openssl) LoadOrCreateCA(filename string, keyfile string) (*CA, error) {
	ca, err := o.LoadCA(filename, keyfile)
	if err != nil {
		return o.CreateCA(filename, keyfile)
	}

	return ca, nil
}

func (o *Openssl) LoadCA(filename string, keyfile string) (*CA, error) {
	var err error
	o.Init()

	filename = o.Path + "/ca/" + filename
	keyfile = o.Path + "/ca/" + keyfile

	c := &CA{}
	c.path = filename
	c.key = keyfile
	c.config = o.GetConfigFile()
	c.crl = o.Path + "/common/crl.pem"

	c.content, err = ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if keyfile != "" {
		c.contentKey, err = ioutil.ReadFile(keyfile)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (o *Openssl) CreateCA(filename string, keyfile string) (*CA, error) {
	o.Init()

	filename = o.Path + "/ca/" + filename
	keyfile = o.Path + "/ca/" + keyfile

	cert := &CA{
		path:   filename,
		key:    keyfile,
		config: o.GetConfigFile(),
		crl:    o.Path + "/common/crl.pem",
	}

	content, err := exec.Command(
		"openssl", "req",
		"-days", "3650",
		"-nodes",
		"-new",
		"-x509",
		"-keyout", "/dev/stdout",
		"-out", "/dev/stdout",
		"-config", cert.config,
		"-batch",
		"-utf8",
		"-subj", "/C="+o.Country+"/ST="+o.Province+"/L="+o.City+"/O="+o.Organization+"/CN="+o.CommonName+"/emailAddress="+o.Email,
	).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("openssl req: " + err.Error() + " (" + string(content) + ")")
	}

	reCert := regexp.MustCompile("(?ms)-----BEGIN CERTIFICATE-----(.+)-----END CERTIFICATE-----")
	reKey := regexp.MustCompile("(?ms)-----BEGIN PRIVATE KEY-----(.+)-----END PRIVATE KEY-----")

	cert.content = reCert.Find(content)
	cert.contentKey = reKey.Find(content)

	if len(cert.content) == 0 {
		err = fmt.Errorf("Generated certificate is 0 long")
		return nil, err
	}

	if len(cert.contentKey) == 0 {
		err = fmt.Errorf("Generated key is 0 long")
		return nil, err
	}

	if err = ioutil.WriteFile(cert.path, cert.content, 0600); err != nil {
		return nil, err
	}
	if err = ioutil.WriteFile(cert.key, cert.contentKey, 0600); err != nil {
		return nil, err
	}

	// Uppdate the CRL (client revoke list)
	content, err = exec.Command(
		"openssl", "ca",
		"-gencrl",
		"-out", "/dev/stdout",
		"-config", cert.config,
		"-batch",
	).CombinedOutput()
	if err != nil {
		return cert, fmt.Errorf("openssl gencrl: " + err.Error() + " (" + string(content) + ")")
	}

	if len(content) == 0 {
		err = fmt.Errorf("Generated CRL is 0 in length")
		return cert, err
	}

	if err = ioutil.WriteFile(cert.crl, content, 0600); err != nil {
		return cert, err
	}

	return cert, nil
}

func (ca *CA) Sign(request *CSR) (*Cert, error) {
	if ca == nil {
		return nil, fmt.Errorf("No CA was supplied")
	}

	cmd := exec.Command(
		"openssl", "ca",
		"-days", "3650",
		"-in", "/dev/stdin",
		"-out", "/dev/stdout",
		"-config", ca.config,
		"-utf8",
		"-batch",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	pipe, _ := cmd.StdinPipe()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	pipe.Write(request.content)
	io.WriteString(pipe, "\n")

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("Exit code ", err, " - ", string(stderr.Bytes()))
	}

	cert := &Cert{
		content:    stdout.Bytes(),
		contentKey: request.contentKey,
	}

	if len(cert.content) == 0 {
		return nil, fmt.Errorf("Generated certificate is 0 long")
	}

	return cert, nil
}

func (ca *CA) Revoke(cert *Cert) error {
	cmd := exec.Command(
		"openssl", "ca",
		"-revoke", "/dev/stdin",
		"-config", ca.config,
		"-utf8",
		"-batch",
	)

	var b bytes.Buffer
	cmd.Stdout = &b
	pipe, _ := cmd.StdinPipe()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Exit code ", err, " - ", string(b.Bytes()))
	}

	pipe.Write(cert.content)
	io.WriteString(pipe, "\n")

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("Exit code ", err, " - ", string(b.Bytes()))
	}

	// Uppdate the CRL (client revoke list)
	content, err := exec.Command(
		"openssl", "ca",
		"-gencrl",
		"-out", "/dev/stdout",
		"-config", ca.config,
		"-batch",
	).Output()
	if err != nil {
		return err
	}

	if len(content) == 0 {
		err = fmt.Errorf("Generated CRL is 0 in length")
		return err
	}

	if err = ioutil.WriteFile(ca.crl, content, 0600); err != nil {
		return err
	}

	if err = os.Remove(cert.GetFilePath()); err != nil {
		return err
	}

	if err = os.Remove(cert.GetKeyPath()); err != nil {
		return err
	}

	return nil
}

func (ca *CA) GetFilePath() string {
	return ca.path
}

func (ca *CA) GetCRLPath() string {
	return ca.crl
}

func (ca *CA) String() string {
	if ca != nil {
		return string(ca.content)
	}

	return ""
}
