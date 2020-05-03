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

func genCA() (*rsa.PrivateKey, []byte, *x509.Certificate) {
	// Generates CA cert
	cadef := &x509.Certificate{
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

	caKeyBin, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err.Error())
	}

	caCertBin, err := x509.CreateCertificate(rand.Reader, cadef, cadef, &caKeyBin.PublicKey, caKeyBin)
	if err != nil {
		panic(err.Error())
	}

	caCert := new(bytes.Buffer)
	pem.Encode(caCert, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCertBin,
	})

	caKey := new(bytes.Buffer)
	pem.Encode(caKey, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caKeyBin),
	})

	_ = ioutil.WriteFile("./caKey.pem", caKey.Bytes(), 0775)
	_ = ioutil.WriteFile("./caCert.pem", caCert.Bytes(), 0775)

	return caKeyBin, caCertBin, cadef
}

func genSCert() (*rsa.PrivateKey, *x509.Certificate) {
	// Device Cert information
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

	sKeyBin, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err.Error())
	}

	return sKeyBin, sCertDef
}

func signCert(caKey *rsa.PrivateKey, caCert []byte, sKeyBin *rsa.PrivateKey, sCertDef *x509.Certificate, caCertDef *x509.Certificate) []byte {
	sCertBin, err := x509.CreateCertificate(rand.Reader, sCertDef, caCertDef, &caKey.PublicKey, caKey)
	if err != nil {
		panic(err.Error())
	}

	sCert := new(bytes.Buffer)
	pem.Encode(sCert, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: sCertBin,
	})

	sKey := new(bytes.Buffer)
	pem.Encode(sKey, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(sKeyBin),
	})

	_ = ioutil.WriteFile("./sKey.pem", sKey.Bytes(), 0775)
	_ = ioutil.WriteFile("./sCert.pem", sCert.Bytes(), 0775)

	return sCertBin
}

func main() {
	caKey, caCert, caCertDef := genCA()
	sKeyBin, sCertDef := genSCert()
	sCert := signCert(caKey, caCert, sKeyBin, sCertDef, caCertDef)
	if sCert == nil {
		fmt.Println("Error...")
	}
	fmt.Println("Done!")

}
