package openssl

import (
	"io/ioutil"
	"os"
)

type Openssl struct {
	Path string

	Country      string
	Province     string
	City         string
	Organization string
	CommonName   string
	Email        string

	initiated bool
}

func (o *Openssl) Init() {
	if o.initiated {
		return
	}

	o.mkdir(o.Path)
	o.mkdir(o.Path + "/ca")
	o.mkdir(o.Path + "/server")
	o.mkdir(o.Path + "/clients")
	o.mkdir(o.Path + "/common")

	o.write_serial(o.Path + "/common/SERIAL")
	o.write_index(o.Path + "/common/index.txt")

	o.initiated = true
}

func (o *Openssl) GetConfigFile() string {
	o.Init()

	if _, err := os.Stat(o.Path + "/common/openvpn.conf"); os.IsNotExist(err) {
		o.WriteConfigFile(o.Path + "/common/openvpn.conf")
	}

	return o.Path + "/common/openvpn.conf"
}

func (o *Openssl) WriteConfigFile(filename string) {
	content := "# OpenSSL config file\n"
	content += "HOME = \"" + o.Path + "\"\n"
	content += "RANDFILE = $HOME/common/random\n"
	content += "oid_section = new_oids\n"
	content += "\n"
	content += "[ new_oids ]\n"
	content += "[ ca ]\n"
	content += "default_ca = CA_default\n"
	content += "\n"
	content += "[ CA_default ]\n"
	content += "dir = $HOME\n"
	content += "certs = $dir/common\n"
	content += "crl_dir = $dir/common\n"
	content += "database = $dir/common/index.txt\n"
	content += "new_certs_dir = $dir/common\n"
	content += "certificate = $dir/ca/ca.crt\n"
	content += "private_key = $dir/ca/ca.key\n"
	content += "serial = $dir/common/SERIAL\n"
	content += "crl = $dir/common/crl.pem\n"
	content += "RANDFILE = $dir/common/.rand\n"
	content += "x509_extensions = usr_cert\n"
	content += "default_days = 3650\n"
	content += "default_crl_days= 30\n"
	content += "default_md = md5\n"
	content += "preserve = no\n"
	content += "policy = policy_match\n"
	content += "\n"
	content += "[ policy_match ]\n"
	content += "countryName = match\n"
	content += "stateOrProvinceName = match\n"
	content += "organizationName = match\n"
	content += "organizationalUnitName = optional\n"
	content += "commonName = supplied\n"
	content += "emailAddress = optional\n"
	content += "\n"
	content += "[ policy_anything ]\n"
	content += "countryName = optional\n"
	content += "stateOrProvinceName = optional\n"
	content += "localityName = optional\n"
	content += "organizationName = optional\n"
	content += "organizationalUnitName = optional\n"
	content += "commonName = supplied\n"
	content += "emailAddress = optional\n"
	content += "\n"
	content += "[ req ]\n"
	content += "default_bits = 1024\n"
	content += "default_keyfile = privkey.pem\n"
	content += "distinguished_name = req_distinguished_name\n"
	content += "attributes = req_attributes\n"
	content += "x509_extensions = v3_ca\n"
	content += "string_mask = nombstr\n"
	content += "\n"
	content += "[ req_distinguished_name ]\n"
	content += "countryName = Country Name (2 letter code)\n"
	content += "countryName_default = \"" + o.Country + "\"\n"
	content += "countryName_min = 2\n"
	content += "countryName_max = 2\n"
	content += "stateOrProvinceName = State or Province Name (full name)\n"
	content += "stateOrProvinceName_default = \"" + o.Province + "\"\n"
	content += "localityName = Locality Name (eg, city)\n"
	content += "localityName_default = \"" + o.City + "\"\n"
	content += "0.organizationName = Organization Name (eg, company)\n"
	content += "0.organizationName_default = \"" + o.Organization + "\"\n"
	content += "organizationalUnitName = Organizational Unit Name (eg, section)\n"
	content += "commonName = Common Name (eg, your name or your server's hostname)\n"
	content += "commonName_max = 64\n"
	content += "commonName_default = \"" + o.CommonName + "\"\n"
	content += "emailAddress = Email Address\n"
	content += "emailAddress_default = \"" + o.Email + "\"\n"
	content += "emailAddress_max = 40\n"
	content += "\n"
	content += "[ req_attributes ]\n"
	content += "challengePassword = A challenge password\n"
	content += "challengePassword_min = 4\n"
	content += "challengePassword_max = 20\n"
	content += "unstructuredName = An optional company name\n"
	content += "\n"
	content += "[ usr_cert ]\n"
	content += "basicConstraints=CA:FALSE\n"
	content += "nsComment = \"OpenSSL Generated Certificate\"\n"
	content += "subjectKeyIdentifier=hash\n"
	content += "authorityKeyIdentifier=keyid,issuer:always\n"
	content += "\n"
	content += "[ server ]\n"
	content += "basicConstraints=CA:FALSE\n"
	content += "nsCertType = server\n"
	content += "nsComment = \"OpenSSL Generated Server Certificate\"\n"
	content += "subjectKeyIdentifier=hash\n"
	content += "authorityKeyIdentifier=keyid,issuer:always\n"
	content += "\n"
	content += "[ v3_req ]\n"
	content += "basicConstraints = CA:FALSE\n"
	content += "keyUsage = nonRepudiation, digitalSignature, keyEncipherment\n"
	content += "\n"
	content += "[ v3_ca ]\n"
	content += "subjectKeyIdentifier=hash\n"
	content += "authorityKeyIdentifier=keyid:always,issuer:always\n"
	content += "basicConstraints = CA:true\n"
	content += "\n"
	content += "[ crl_ext ]\n"
	content += "authorityKeyIdentifier=keyid:always,issuer:always\n"

	ioutil.WriteFile(filename, []byte(content), 0660)
}

func (o *Openssl) mkdir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0600)
	}
}

func (o *Openssl) AppendPath(filename string) string {
	return o.Path + "/" + filename
}

func (o *Openssl) write_serial(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err = ioutil.WriteFile(filename, []byte("1000"), 0660)
	}
}
func (o *Openssl) write_index(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err = ioutil.WriteFile(filename, []byte(""), 0660)
	}
}
