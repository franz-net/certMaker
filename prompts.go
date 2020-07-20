package main

import (
	"net"

	"github.com/manifoldco/promptui"
)

func initPrompt() string {
	prompt := promptui.Select{
		Label: "Select Certificate Type",
		Items: []string{"Ca", "Certificate"},
	}
	_, certType, _ := prompt.Run()

	if certType == "Ca" {
		return "ca"
	} else if certType == "Certificate" {
		return "cert"
	}
	return ""
}

func loadPaths(certType string) (string, string, string) {

	if certType == "ca" {
		destinationPrompt := promptui.Prompt{
			Label: "Provide the output path where the certificates will be placed",
		}
		outputPath, _ := destinationPrompt.Run()
		checkForPath(outputPath)

		caNamePrompt := promptui.Prompt{
			Label: "Provide a name for the CA (this will be the filename)",
		}
		caName, _ := caNamePrompt.Run()

		return outputPath, caName, ""

	} else if certType == "cert" {
		destinationPrompt := promptui.Prompt{
			Label: "Provide the output path where the certificates will be placed",
		}
		outputPath, _ := destinationPrompt.Run()
		checkForPath(outputPath)

		caPathPrompt := promptui.Prompt{
			Label: "Provide the path to the ca cert that will be used to sign new certificates",
		}
		caPath, _ := caPathPrompt.Run()
		checkForPath(caPath)

		caKeyPathPrompt := promptui.Prompt{
			Label: "Provide the path to the private CA key that will be used to sign new certificates",
		}
		caKeyPath, _ := caKeyPathPrompt.Run()
		checkForPath(caKeyPath)

		return outputPath, caPath, caKeyPath
	}

	return "", "", ""

}

func caPrompt() Cert {
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

	caCertObj := Cert{
		Organization:  []string{caOrg},
		Country:       []string{caCountry},
		Province:      []string{caProvince},
		Locality:      []string{caLocality},
		StreetAddress: []string{caStreetAdress},
		PostalCode:    []string{caPostalCode},
	}

	return caCertObj
}

func continueToCertPrompt() bool {
	prompt := promptui.Select{
		Label: "Continue creating a certificate from the CA just created?",
		Items: []string{"Yes", "No"},
	}
	_, continueCert, _ := prompt.Run()

	if continueCert == "Yes" {
		return true
	}

	return false
}

func certPrompt() Cert {
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

	return certObj
}
