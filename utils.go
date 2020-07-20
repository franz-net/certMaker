package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func checkForPath(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("Error, output path not found")
		os.Exit(1)
	}
	//if !os.IsPermission(err) {
	//	fmt.Println("Error, output path is not writable")
	//	os.Exit(1)
	//}

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

func writeCertKey(key *bytes.Buffer, cert *bytes.Buffer, certType, outputPath, certName string) (string, string) {
	certPath := ""
	keyPath := ""

	if certType == "ca" {
		certPath = outputPath + "/caCert" + certName + ".pem"
		keyPath = outputPath + "/caKey" + certName + ".pem"
	} else {
		certPath = outputPath + "/sCert" + certName + ".pem"
		keyPath = outputPath + "/sKey" + certName + ".pem"
	}

	_ = ioutil.WriteFile(keyPath, key.Bytes(), 0775)
	_ = ioutil.WriteFile(certPath, cert.Bytes(), 0775)

	return certPath, keyPath
}

func getCA(inputCA, inputCaKey string) (*rsa.PrivateKey, []byte, *x509.Certificate) {

	pemKey := readCaData(inputCaKey)
	pemCert := readCaData(inputCA)

	keyBlock, _ := pem.Decode(pemKey)
	var parsedCaKey interface{}
	parsedCaKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		panic(err.Error())
	}
	caKeyBin, _ := parsedCaKey.(*rsa.PrivateKey)
	caCertBlock, _ := pem.Decode(pemCert)
	caCertDef, _ := x509.ParseCertificate(caCertBlock.Bytes)
	caCertBin := caCertBlock.Bytes

	return caKeyBin, caCertBin, caCertDef
}

func readCaData(inputPath string) []byte {
	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error, could not read the file: " + inputPath)
		os.Exit(1)
	}

	return data
}
