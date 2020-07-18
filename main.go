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

func main() {

	var certInfo Cert
	var certType string

	if len(os.Args) > 1 {
		certInfo, certType = readFlags()
	} else {
		certInfo, certType = loadPromptui()
	}

	if certType == "ca" {
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
		writeCertKey(caPemKey, caPemCert, "ca")
		// create struct
	}

	if certType == "cert" {
		fmt.Println("Generating signed Cert...")
		// retrieve ca information
		caKeyBin, caCertBin, caCertDef := getCA()
		// generate server cert key
		sKeyBin := genPrivateKey()
		// generate cert definition
		sCertDef := genSCertDef(certInfo)
		// Sign cert with CA
		sCertBin := signCert(caKeyBin, caCertBin, sKeyBin, sCertDef, caCertDef)
		// pem encode the scert and return it
		sPemCert, sPemKey := encodeCertKey(sKeyBin, sCertBin)
		// write the cert to file
		writeCertKey(sPemKey, sPemCert, "")
	}
}
