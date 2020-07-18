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
