package main

import (
	"fmt"
	"net"
	"os"

	"github.com/manifoldco/promptui"
)

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
		orgPrompt := promptui.Prompt{
			Label: "Provide CA Organization Name",
		}
		caOrg, _ := orgPrompt.Run()

		coPrompt := promptui.Prompt{
			Label: "Provide CA Country (2 letter)",
		}
		caCountry, _ := coPrompt.Run()

		provincePrompt := promptui.Prompt{
			Label: "Provide CA Province Name",
		}
		caProvince, _ := provincePrompt.Run()

		localityPrompt := promptui.Prompt{
			Label: "Provide CA Locality Name",
		}
		caLocality, _ := localityPrompt.Run()

		addressPrompt := promptui.Prompt{
			Label: "Provide CA Address",
		}
		caStreetAdress, _ := addressPrompt.Run()

		zipPrompt := promptui.Prompt{
			Label: "Provide CA Zip Code",
		}
		caPostalCode, _ := zipPrompt.Run()

		certObj := Cert{
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

		orgPrompt := promptui.Prompt{
			Label: "Provide Organization for Certificate",
		}
		sOrg, _ := orgPrompt.Run()

		coPrompt := promptui.Prompt{
			Label: "Provide Country for Certificate (2 letter)",
		}
		sCountry, _ := coPrompt.Run()

		provincePrompt := promptui.Prompt{
			Label: "Provide Province for Certificate",
		}
		sProvince, _ := provincePrompt.Run()

		localityPrompt := promptui.Prompt{
			Label: "Provide Locality Name for Certificate",
		}
		sLocality, _ := localityPrompt.Run()

		addressPrompt := promptui.Prompt{
			Label: "Provide Address for Certificate",
		}
		sStreetAdress, _ := addressPrompt.Run()

		zipPrompt := promptui.Prompt{
			Label: "Provide Zip Code for Certificate",
		}
		sPostalCode, _ := zipPrompt.Run()

		commonNamePrompt := promptui.Prompt{
			Label: "Provide common name (hostname, web address or IP) to associate with the Certificate",
		}
		sCommonName, _ := commonNamePrompt.Run()

		sanIPPrompt := promptui.Prompt{
			Label: "Provide IPv4 alternative address (i.e 127.0.0.1) to associate with the Certificate",
		}
		sAltSubjectIP, _ := sanIPPrompt.Run()

		sanDNSPrompt := promptui.Prompt{
			Label: "Provide alternative hostnames or web address (i.e localhost) to associate with the Certificate",
		}
		sAltSubjectName, _ := sanDNSPrompt.Run()

		certObj := Cert{
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
