package helper

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

// Helper function to get certificate And certificate pool
func GetCertificateAndCertPool(CertFilePath string, KeyFilePath string, CertAuthFilePath string) (*tls.Certificate, *x509.CertPool, error) {
	// Load the client certificates from disk
	certificate, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)
	if err != nil {
		fmt.Print("Could not load client key pair: ", err)
		return nil, nil, err
	}

	// Create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(CertAuthFilePath)
	if err != nil {
		fmt.Print("Could not read ca certificate: ", err)
		return nil, nil, err
	}

	// Append the certificates from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		fmt.Print("Failed to append ca certs")
		return nil, nil, err
	}

	return &certificate, certPool, err
}
