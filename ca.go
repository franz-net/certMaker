package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

func genCA(caKeyBin *rsa.PrivateKey, caDef *x509.Certificate) []byte {
	// Generates CA cert binary

	caCertBin, err := x509.CreateCertificate(rand.Reader, caDef, caDef, &caKeyBin.PublicKey, caKeyBin)
	if err != nil {
		panic(err.Error())
	}

	return caCertBin
}

func genCaCertDef(certInfo Cert) *x509.Certificate {
	// Generates x509 certificate
	caDef := &x509.Certificate{
		SerialNumber: big.NewInt(2020),
		Subject: pkix.Name{
			Organization:  certInfo.Organization,
			Country:       certInfo.Country,
			Province:      certInfo.Province,
			Locality:      certInfo.Locality,
			StreetAddress: certInfo.StreetAddress,
			PostalCode:    certInfo.PostalCode,
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

func getCA() (*rsa.PrivateKey, []byte, *x509.Certificate) {
	// NEED TO MODIFY THIS TO LOAD A CA FROM FILE
	var pemKey string
	var pemCert string

	keyBlock, _ := pem.Decode([]byte(pemKey))
	var parsedCaKey interface{}
	parsedCaKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		panic(err.Error())
	}
	caKeyBin, _ := parsedCaKey.(*rsa.PrivateKey)
	caCertBlock, _ := pem.Decode([]byte(pemCert))
	caCertDef, _ := x509.ParseCertificate(caCertBlock.Bytes)
	caCertBin := caCertBlock.Bytes

	return caKeyBin, caCertBin, caCertDef
}
