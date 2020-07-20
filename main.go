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

func ca(caInfo Cert, outputPath, caName string) {
	fmt.Println("Generating CA...")
	// generate ca key
	caKeyBin := genPrivateKey()
	// generate cert definition
	caCertDef := genCaCertDef(caInfo)
	// generate ca cert and return it
	caCertBin := genCA(caKeyBin, caCertDef)
	// pem encode the cacert and return it
	caPemCert, caPemKey := encodeCertKey(caKeyBin, caCertBin)
	// Write the cert to file
	caPath, caKeyPath := writeCertKey(caPemKey, caPemCert, "ca", outputPath, caName)
	// Return the path to the keys
	fmt.Println("Ca can be found: " + caPath + "\nKey can be found: " + caKeyPath + "/nKeep the key safe!")
	// ask if the user wants to create a cert based on the CA
	for {
		if !continueToCertPrompt() {
			break
		}
		certInfo := certPrompt()
		cert(certInfo, outputPath, caPath, caKeyPath)
	}
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
	certPath, keyPath := writeCertKey(sPemKey, sPemCert, "", outputPath, certInfo.CommonName)
	// Return paths to the user
	fmt.Println("Cert can be found: " + certPath + "\nKey can be found: " + keyPath)
}

func main() {

	if len(os.Args) > 1 {
		//certInfo, certType := readFlags()
		// WILL ADD CLI FLAGS HERE EVENTUALLY
	} else {
		certType := initPrompt()
		if certType == "ca" {
			outputPath, caName, _ := loadPaths(certType)
			certInfo := caPrompt()

			ca(certInfo, outputPath, caName)

		} else if certType == "cert" {
			outputPath, caPath, caKeyPath := loadPaths(certType)
			certInfo := certPrompt()

			cert(certInfo, outputPath, caPath, caKeyPath)
		}
	}
}
