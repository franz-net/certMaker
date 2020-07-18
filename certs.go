package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

func genPrivateKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err.Error())
	}

	return privateKey
}

func genSCertDef(certInfo Cert) *x509.Certificate {
	// Generates x509 certificate
	sCertDef := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  certInfo.Organization,
			Country:       certInfo.Country,
			Province:      certInfo.Province,
			Locality:      certInfo.Locality,
			StreetAddress: certInfo.StreetAddress,
			PostalCode:    certInfo.PostalCode,
			CommonName:    certInfo.CommonName,
		},
		IPAddresses:  certInfo.IPAddresses, //[]net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:     certInfo.DNSNames,
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
