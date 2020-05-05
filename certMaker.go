/*
certMaker.go
Version: 1.0
Author: Franz Ramirez
Description: Application that allows to generate self generated CAs and the corresponding server certificate
*/
package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"time"
)

func genPrivateKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err.Error())
	}

	return privateKey
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

	_ = ioutil.WriteFile(certPath, key.Bytes(), 0775)
	_ = ioutil.WriteFile(keyPath, cert.Bytes(), 0775)
}

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
func genCaCertDef() *x509.Certificate {
	// Generates x509 certificate
	caDef := &x509.Certificate{
		SerialNumber: big.NewInt(2020),
		Subject: pkix.Name{
			Organization:  []string{"Example, Inc."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Market St. #2"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	return caDef
}

func genCA(caKeyBin *rsa.PrivateKey, caDef *x509.Certificate) []byte {
	// Generates CA cert binary

	caCertBin, err := x509.CreateCertificate(rand.Reader, caDef, caDef, &caKeyBin.PublicKey, caKeyBin)
	if err != nil {
		panic(err.Error())
	}

	return caCertBin
}

func genSCertDef() *x509.Certificate {
	// Generates x509 certificate
	sCertDef := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"Example, Inc."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Market St. #2"},
			PostalCode:    []string{"94016"},
			CommonName:    "goliath",
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	return sCertDef
}

func signCert(caKey *rsa.PrivateKey, caCert []byte, sKeyBin *rsa.PrivateKey, sCertDef *x509.Certificate, caCertDef *x509.Certificate) []byte {

	sCertBin, err := x509.CreateCertificate(rand.Reader, sCertDef, caCertDef, &caKey.PublicKey, caKey)
	if err != nil {
		panic(err.Error())
	}

	return sCertBin
}

func main() {
	fmt.Println("Generating CA...")
	// generate ca key
	caKeyBin := genPrivateKey()
	// generate cert definition
	caCertDef := genCaCertDef()
	// generate ca cert and return it
	caCertBin := genCA(caKeyBin, caCertDef)
	// pem encode the cacert and return it
	caPemCert, caPemKey := encodeCertKey(caKeyBin, caCertBin)
	// Write the cert to file
	writeCertKey(caPemKey, caPemCert, "ca")

	fmt.Println("Generating signed Cert...")
	// generate server cert key
	sKeyBin := genPrivateKey()
	// generate cert definition
	sCertDef := genSCertDef()
	// Sign cert with CA
	sCertBin := signCert(caKeyBin, caCertBin, sKeyBin, sCertDef, caCertDef)
	// pem encode the scert and return it
	sPemCert, sPemKey := encodeCertKey(sKeyBin, sCertBin)
	// write the cert to file
	writeCertKey(sPemKey, sPemCert, "")

	fmt.Println("Done!")

}
