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
	"database/sql"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/manifoldco/promptui"

	_ "github.com/mattn/go-sqlite3"
)

//Cert struct
type Cert struct {
	ID            string
	CAID          string
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

//DbCert struct
type DbCert struct {
	ID      string
	CAID    string
	PEMKEY  bytes.Buffer
	PEMCERT bytes.Buffer
}

func dbConn() (db *sql.DB) {
	_, err := os.Stat("certs.db")
	if os.IsNotExist(err) {
		fmt.Println("DB not found, creating it...")
		os.Create("certs.db")
	}

	db, err = sql.Open("sqlite3", "certs.db")
	if err != nil {
		panic(err.Error())
	}

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS certs (id TEXT PRIMARY KEY, caId TEXT, pemKey TEXT, pemCert TEXT)")
	statement.Exec()

	return db
}

func writeToDb(certData DbCert, caType string) {
	db := dbConn()
	if caType == "ca" {
		id := certData.ID
		pemkey := certData.PEMKEY.Bytes()
		pemcert := certData.PEMCERT.Bytes()

		ins, err := db.Prepare("INSERT INTO certs(id, caId, pemKey, pemCert) VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		_, err = ins.Exec(id, "", pemkey, pemcert)
		if err != nil {
			fmt.Println("ERROR: could not insert data into the db")
			panic(err.Error())
		}
	}
	if caType == "cert" {
		id := certData.ID
		caid := certData.CAID
		pemkey := certData.PEMKEY.Bytes()
		pemcert := certData.PEMCERT.Bytes()

		ins, err := db.Prepare("INSERT INTO certs(id, caId, pemKey, pemCert) VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		_, err = ins.Exec(id, caid, pemkey, pemcert)
		if err != nil {
			fmt.Println("ERROR: could not insert data into the db")
			panic(err.Error())
		}
	}
	defer db.Close()
}

func getCA(cid string) (*rsa.PrivateKey, []byte, *x509.Certificate) {
	db := dbConn()
	id := cid
	var pemKey string
	var pemCert string
	err := db.QueryRow("SELECT pemKey, pemCert FROM certs WHERE id = ?", id).Scan(&pemKey, &pemCert)
	if err != nil {
		fmt.Println("ERROR: could not query the database")
		panic(err.Error())
	}
	//fmt.Println(pemKey + "\n\n" + pemCert)

	keyBlock, _ := pem.Decode([]byte(pemKey))
	var parsedCaKey interface{}
	parsedCaKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		panic(err.Error())
	}
	caKeyBin, _ := parsedCaKey.(*rsa.PrivateKey)
	caCertBlock, _ := pem.Decode([]byte(pemCert))
	caCertDef, _ := x509.ParseCertificate(caCertBlock.Bytes)
	caCertBin := caCertBlock.Bytes

	return caKeyBin, caCertBin, caCertDef
}

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

	_ = ioutil.WriteFile(keyPath, key.Bytes(), 0775)
	_ = ioutil.WriteFile(certPath, cert.Bytes(), 0775)
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

func genCA(caKeyBin *rsa.PrivateKey, caDef *x509.Certificate) []byte {
	// Generates CA cert binary

	caCertBin, err := x509.CreateCertificate(rand.Reader, caDef, caDef, &caKeyBin.PublicKey, caKeyBin)
	if err != nil {
		panic(err.Error())
	}

	return caCertBin
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

func readFlags() (Cert, string) {

	var certObj Cert
	var certType string

	caCmd := flag.NewFlagSet("ca", flag.ExitOnError)
	// CA Flags
	caID := caCmd.String("ca-identifier", "", "Internal Identifier for the CA (Required to create a CA or sign a Cert)")
	caOrg := caCmd.String("ca-organization", "Example, Inc", "Certificate Authority Name")
	caCountry := caCmd.String("ca-country", "US", "Certificate Authority Country (2 letter)")
	caProvince := caCmd.String("ca-province", "CA", "Certificate Authority Province (2 letter)")
	caLocality := caCmd.String("ca-locality", "San Francisco", "Certificate Authoritiy Locality")
	caStreetAdress := caCmd.String("ca-address", "101 Market St", "Certificate Authority Street Address")
	caPostalCode := caCmd.String("ca-zipcode", "94016", "Certificate Authority Zip Code")

	certCmd := flag.NewFlagSet("cert", flag.ExitOnError)
	// Device Cert Flags
	sOrg := certCmd.String("server-organization", "Example, Inc", "Organization where the signed certificate will be installed")
	sCountry := certCmd.String("server-country", "US", "Country code where certificate will be installed (2 letters)")
	sProvince := certCmd.String("server-province", "CA", "Province where certificate will be installed (2 letter)")
	sLocality := certCmd.String("server-locality", "San Francisco", "Locality where certificate will be installed")
	sStreetAdress := certCmd.String("server-address", "101 Market St", "Address where certificate will be installed")
	sPostalCode := certCmd.String("server-zipcode", "94016", "Zip Code where certificate will be installed")
	sCommonName := certCmd.String("server-common-name", "", "FQDN or hostname to be used when issuing a certificate")
	sAltSubjectIP := certCmd.String("subject-alternate-ip", "127.0.0.1", "IP that can also be used to address the server where the certificate will be installed")
	sAltSubjectName := certCmd.String("subject-alterative-hostname", "localhost", "Other hostname that can identify the server where the certificate will be installed")
	sID := certCmd.String("server-identifier", "", "Internal Identifier for the device certificate (Required to create and sign a Cert)")
	scaID := certCmd.String("ca-identifier", "", "Internal Identifier for the CA (Required to create a CA or sign a Cert)")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'ca' or 'cert' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "ca":
		caCmd.Parse(os.Args[2:])

		if *caID == "" {
			fmt.Println("Error: CA Identifier needed to create CA")
			caCmd.PrintDefaults()
			os.Exit(1)
		}

		certObj := Cert{
			ID:            *caID,
			Organization:  []string{*caOrg},
			Country:       []string{*caCountry},
			Province:      []string{*caProvince},
			Locality:      []string{*caLocality},
			StreetAddress: []string{*caStreetAdress},
			PostalCode:    []string{*caPostalCode},
		}

		return certObj, "ca"

	case "cert":
		certCmd.Parse(os.Args[2:])

		fmt.Println("sid: " + *sID + "\nsCommonName: " + *sCommonName + "\nscaID" + *scaID)

		if *sID == "" || *sCommonName == "" || *scaID == "" {
			fmt.Println("ERROR: missing one or more parameters (server-identifier, ca-identifier or common-name")
			certCmd.PrintDefaults()
			os.Exit(1)
		}

		certObj := Cert{
			ID:            *sID,
			CAID:          *scaID,
			Organization:  []string{*sOrg},
			Country:       []string{*sCountry},
			Province:      []string{*sProvince},
			Locality:      []string{*sLocality},
			StreetAddress: []string{*sStreetAdress},
			PostalCode:    []string{*sPostalCode},
			CommonName:    *sCommonName,
			IPAddresses:   []net.IP{net.ParseIP(*sAltSubjectIP)},
			DNSNames:      []string{*sAltSubjectName},
		}

		return certObj, "cert"

	default:
		fmt.Println("Error: expected 'ca' or 'cert' subcommands")
		fmt.Println("CA flags:")
		caCmd.PrintDefaults()
		fmt.Println("\nCert flags:")
		certCmd.PrintDefaults()
		os.Exit(1)
	}

	return certObj, certType
}

func loadPromptui() (Cert, string) {

	var certObj Cert
	var certType string

	prompt := promptui.Select{
		Label: "Select Certificate Type",
		Items: []string{"CaCert", "SignedCert"},
	}

	_, certType, err := prompt.Run()

	if err != nil {
		fmt.Println("ERROR: Prompt failed", err)
		os.Exit(1)
	}

	if certType == "CaCert" {
		// prompt for CA stuff
		id_prompt := promptui.Prompt{
			Label: "Provide CA Id or Name to identify",
		}
		caID, _ := id_prompt.Run()

		org_prompt := promptui.Prompt{
			Label: "Provide CA Organization Name",
		}
		caOrg, _ := org_prompt.Run()

		co_prompt := promptui.Prompt{
			Label: "Provide CA Country (2 letter)",
		}
		caCountry, _ := co_prompt.Run()

		province_prompt := promptui.Prompt{
			Label: "Provide CA Province Name",
		}
		caProvince, _ := province_prompt.Run()

		locality_prompt := promptui.Prompt{
			Label: "Provide CA Locality Name",
		}
		caLocality, _ := locality_prompt.Run()

		address_prompt := promptui.Prompt{
			Label: "Provide CA Address",
		}
		caStreetAdress, _ := address_prompt.Run()

		zip_prompt := promptui.Prompt{
			Label: "Provide CA Zip Code",
		}
		caPostalCode, _ := zip_prompt.Run()

		certObj := Cert{
			ID:            caID,
			Organization:  []string{caOrg},
			Country:       []string{caCountry},
			Province:      []string{caProvince},
			Locality:      []string{caLocality},
			StreetAddress: []string{caStreetAdress},
			PostalCode:    []string{caPostalCode},
		}

		return certObj, "ca"

	} else if certType == "SignedCert" {
		// prompt for device stuff
		sid_prompt := promptui.Prompt{
			Label: "Provide Id or Name to identify the Signed Certificate",
		}
		sID, _ := sid_prompt.Run()

		id_prompt := promptui.Prompt{
			Label: "Provide the CA Id to sign the Certificate against",
		}
		scaID, _ := id_prompt.Run()

		org_prompt := promptui.Prompt{
			Label: "Provide Organization for Certificate",
		}
		sOrg, _ := org_prompt.Run()

		co_prompt := promptui.Prompt{
			Label: "Provide Country for Certificate (2 letter)",
		}
		sCountry, _ := co_prompt.Run()

		province_prompt := promptui.Prompt{
			Label: "Provide Province for Certificate",
		}
		sProvince, _ := province_prompt.Run()

		locality_prompt := promptui.Prompt{
			Label: "Provide Locality Name for Certificate",
		}
		sLocality, _ := locality_prompt.Run()

		address_prompt := promptui.Prompt{
			Label: "Provide Address for Certificate",
		}
		sStreetAdress, _ := address_prompt.Run()

		zip_prompt := promptui.Prompt{
			Label: "Provide Zip Code for Certificate",
		}
		sPostalCode, _ := zip_prompt.Run()

		common_name_prompt := promptui.Prompt{
			Label: "Provide common name (hostname, web address or IP) to associate with the Certificate",
		}
		sCommonName, _ := common_name_prompt.Run()

		sanIP_prompt := promptui.Prompt{
			Label: "Provide IPv4 alternative address (i.e 127.0.0.1) to associate with the Certificate",
		}
		sAltSubjectIP, _ := sanIP_prompt.Run()

		sanDNS_prompt := promptui.Prompt{
			Label: "Provide alternative hostnames or web address (i.e localhost) to associate with the Certificate",
		}
		sAltSubjectName, _ := sanDNS_prompt.Run()

		certObj := Cert{
			ID:            sID,
			CAID:          scaID,
			Organization:  []string{sOrg},
			Country:       []string{sCountry},
			Province:      []string{sProvince},
			Locality:      []string{sLocality},
			StreetAddress: []string{sStreetAdress},
			PostalCode:    []string{sPostalCode},
			CommonName:    sCommonName,
			IPAddresses:   []net.IP{net.ParseIP(sAltSubjectIP)},
			DNSNames:      []string{sAltSubjectName},
		}

		return certObj, "cert"
	}
	return certObj, certType
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
		certData := DbCert{
			ID:      certInfo.ID,
			PEMKEY:  *caPemKey,
			PEMCERT: *caPemCert,
		}
		// Save CA in db
		writeToDb(certData, "ca")
	}

	if certType == "cert" {
		fmt.Println("Generating signed Cert...")
		// retrieve ca information
		caKeyBin, caCertBin, caCertDef := getCA(certInfo.CAID)
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
