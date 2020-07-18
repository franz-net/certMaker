package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

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
