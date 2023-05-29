package utils

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
)

// CertificateInfo returns a string describing the certificate
func CertificateInfo(cert *x509.Certificate) string {
	if cert.Subject.CommonName == cert.Issuer.CommonName {
		return fmt.Sprintf("    Self-signed certificate %v\n", cert.Issuer.CommonName)
	}

	s := fmt.Sprintf("    Subject %v\n", cert.DNSNames)
	s += fmt.Sprintf("    Issued by %s\n", cert.Issuer.CommonName)
	return s
}

// OutputPEMFile reads info from a PEM file and displays it
func OutputPEMFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	for len(data) > 0 {
		var block *pem.Block
		block, data = pem.Decode(data)
		log.Printf("Type: %#v\n", block.Type)
		switch block.Type {
		case "CERTIFICATE":
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return err
			}
			log.Printf(CertificateInfo(cert))
		default:
			log.Printf(block.Type)
		}
	}
	return nil
}
