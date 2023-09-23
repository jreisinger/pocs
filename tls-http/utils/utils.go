package utils

import (
	"crypto/tls"
	"fmt"
	"log"
)

func getCert(certfile, keyfile string) (c tls.Certificate, err error) {
	if certfile != "" && keyfile != "" {
		c, err = tls.LoadX509KeyPair(certfile, keyfile)
		if err != nil {
			log.Printf("error loading key pair: %v\n", err)
		}
	} else {
		err = fmt.Errorf("i have no certificate")
	}
	return
}

// CertReqFunc returns a function for tlsConfig.GetCertificate
func CertReqFunc(certfile, keyfile string) func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	c, err := getCert(certfile, keyfile)

	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		log.Printf("Received TLS Hello asking for %s: sending certificate\n", hello.ServerName)
		if err != nil || certfile == "" {
			log.Println("I have no certificate")
		} else {
			err := OutputPEMFile(certfile)
			if err != nil {
				log.Printf("%v\n", err)
			}
		}
		Wait()
		return &c, nil
	}
}
