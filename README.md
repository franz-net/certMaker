<h1 align="center">
  CertMaker
</h1>

<p align="center">
  <a href="https://github.com/franz-net/certMaker/blob/master/LICENSE">
    <img alt="GPL3" src="https://img.shields.io/github/license/franz-net/certMaker">
  </a>
  <a href="https://goreportcard.com/report/github.com/franz-net/certMaker">
   <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/franz-net/certMaker">
  </a>
  <a href="https://travis-ci.org/github/franz-net/certMaker">
    <img alt="build status" src="https://travis-ci.org/franz-net/certMaker.svg?branch=master">
  </a>  
</p>

Project created with the goal of helping me learn how to write Go code and making self-signed certificate creation.
Empowers users with a single executable to create and store CA certificates, private keys and sign new device certificates.


## Installation (Building the code)

To build and run the program you need: 
* [Golang installed](https://golang.org/doc/install) and [GOPATH configured](https://golang.org/doc/gopath_code.html)
* Install [promptui](https://github.com/manifoldco/promptui) packages
* Clone this repository
* Build the code with `go build...`

 Example of build for linux x86_64: `env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o certMaker certMaker.go`

## Pre-compiled Binaries

To use:
* Download the certMaker executable
* Grant execute permissions `chmod +x certMaker`
* Add to $PATH if necessary
* Run certMaker `./certMaker`

### Releases
* [Version 1.0](https://github.com/franz-net/certMaker/releases/tag/release%2Fv1.0)

## Usage

After building or downloading one of the pre-compiled binaries, run `certMaker` and respond to the prompts
The resulting certificates and keys are generated in the path provided

### Create CA

#### **Flags and Arguments**

  * `-ca-address`: Certificate Authority Street Address (default "101 Market St")
  * `-ca-country`: Certificate Authority Country (2 letter) (default "US")
  * `-ca-identifier`: Internal Identifier for the CA (Required to create a CA or sign a Cert)
  * `-ca-locality`: Certificate Authoritiy Locality (default "San Francisco")
  * `-ca-organization`: Certificate Authority Name (default "Example, Inc")
  * `-ca-province`: Certificate Authority Province (2 letter) (default "CA")
  * `-ca-zipcode`: Certificate Authority Zip Code (default "94016")

#### **Prompt Mode**
![](CA_Prompts.gif)

### Create and Sign Certificate

#### **Flags and Arguments**

  * `-ca-identifier`: ID or name of CA that will be used to signed the Certificate
  * `-server-identifier`: ID or name for the Certificate being created
  * `-server-address`: Address for the Certificate
  * `-server-common-name`: FQDN or Hostname the Certificate will be identifying
  * `-server-country`: Country for the Certificate (2 letters)
  * `-server-locality`: Locality for the Certificate
  * `-server-organization`: Organization for the Certificate
  * `-server-province`: Province for the Certificate
  * `-server-zipcode`: Zip Code for the Certificate
  * `-subject-alterative-hostname`: Alternate hostnames that the Certificate can be idetified as
  * `-subject-alternate-ip`: Alternate IP that the Certificate can be identified as

#### **Prompt Mode**
![](S_Prompts.gif)

### Features
* The "wizard" mode uses the <a href="https://github.com/manifoldco/promptui">promptui</a> package
* Certificates can be verified using OpenSSL ```openssl verify -verbose -CAfile cacert.pem  server.crt```

## Wishlist
* Database to store CAs, Key's and Certs
* Webapp to manage the certificates
* Add tests
* Add silent mode to facilitate automation