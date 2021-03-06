package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aktsk/kalvados/server"
	"github.com/aktsk/kalvados/version"
)

const name = "kalvados-server"

var GitCommit string

func main() {
	var (
		port         int
		keyFileName  string
		certFileName string
		versionFlag  bool
	)

	flag.IntVar(&port, "port", 8000, "Port to listen")
	flag.StringVar(&keyFileName, "keyFile", "key.pem", "Private Key file")
	flag.StringVar(&certFileName, "certFile", "cert.pem", "Cetificate file")
	flag.BoolVar(&versionFlag, "version", false, "print version string")

	flag.Parse()

	if versionFlag {
		fmt.Printf("%s version: %s (rev: %s)", name, version.Get(), GitCommit)
		os.Exit(0)
	}

	keyFile, err := os.Open(keyFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer keyFile.Close()

	keyPEM, err := ioutil.ReadAll(keyFile)
	if err != nil {
		log.Fatal(err)
	}

	keyDER, _ := pem.Decode(keyPEM)
	key, err := x509.ParsePKCS1PrivateKey(keyDER.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	certFile, err := os.Open(certFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer certFile.Close()

	certPEM, err := ioutil.ReadAll(certFile)
	if err != nil {
		log.Fatal(err)
	}

	certDER, _ := pem.Decode(certPEM)
	cert, err := x509.ParseCertificate(certDER.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	server.Serve(port, key, cert)
}
