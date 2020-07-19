/*
certMaker.go
Version: 1.0
Author: Franz Ramirez
Description: Application that allows to generate self generated CAs and the corresponding server certificate
*/
package main

import (
	"fmt"
	"net"
	"os"
)

//Cert struct
type Cert struct {
	Organization  []string
	Country       []string
	Province      []string
	Locality      []string
	StreetAddress []string
	PostalCode    []string
	CommonName    string
	IPAddresses   []net.IP
	DNSNames      []string
}

func ca(certInfo Cert, outputPath string) {
	fmt.Println("Generating CA...")
	// generate ca key
	caKeyBin := genPrivateKey()
	// generate cert definition
	caCertDef := genCaCertDef(certInfo)
	// generate ca cert and return it
	caCertBin := genCA(caKeyBin, caCertDef)
	// pem encode the cacert and return it
	caPemCert, caPemKey := encodeCertKey(caKeyBin, caCertBin)
	// Write the cert to file
	writeCertKey(caPemKey, caPemCert, "ca", outputPath)
	// create struct
}

func cert(certInfo Cert, outputPath string, caPath string, caKeyPath string) {
	fmt.Println("Generating signed Cert...")
	// retrieve ca information
	caKeyBin, caCertBin, caCertDef := getCA(caPath, caKeyPath)
	// generate server cert key
	sKeyBin := genPrivateKey()
	// generate cert definition
	sCertDef := genSCertDef(certInfo)
	// Sign cert with CA
	sCertBin := signCert(caKeyBin, caCertBin, sKeyBin, sCertDef, caCertDef)
	// pem encode the scert and return it
	sPemCert, sPemKey := encodeCertKey(sKeyBin, sCertBin)
	// write the cert to file
	writeCertKey(sPemKey, sPemCert, "", outputPath)
}

func main() {

	if len(os.Args) > 1 {
		//certInfo, certType := readFlags()
		// WILL ADD CLI FLAGS HERE EVENTUALLY
	} else {
		certType := initPrompt()
		if certType == "ca" {
			outputPath, _, _ := loadPaths(certType)
			certInfo := caPrompt()

			ca(certInfo, outputPath)

		} else if certType == "cert" {
			outputPath, caPath, caKeyPath := loadPaths(certType)
			certInfo := certPrompt()

			cert(certInfo, outputPath, caPath, caKeyPath)
		}

	}
}
