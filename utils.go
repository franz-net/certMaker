package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

func encodeCertKey(key *rsa.PrivateKey, cert []byte) (*bytes.Buffer, *bytes.Buffer) {
	pemCert := new(bytes.Buffer)
	pem.Encode(pemCert, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})

	pemKey := new(bytes.Buffer)
	pem.Encode(pemKey, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	return pemCert, pemKey
}

func writeCertKey(key *bytes.Buffer, cert *bytes.Buffer, certType string) {
	certPath := ""
	keyPath := ""

	if certType == "ca" {
		certPath = "./caCert.pem"
		keyPath = "./caKey.pem"
	} else {
		certPath = "./sCert.pem"
		keyPath = "./sKey.pem"
	}

	_ = ioutil.WriteFile(keyPath, key.Bytes(), 0775)
	_ = ioutil.WriteFile(certPath, cert.Bytes(), 0775)
}
