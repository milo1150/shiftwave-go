package network

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func LoadCertificate() tls.Certificate {
	_, err := godotenv.Read(".env")
	if err != nil {
		log.Fatalf("Failed to load certificate: %v", err)
	}

	certFile := os.Getenv("CERT_FILE_PATH")
	keyFile := os.Getenv("KEY_FILE_PATH")
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Error load certificate: %v", err)
	}

	return cert
}

func GetHttpClientWithCert(cert tls.Certificate) *http.Client {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	return httpClient
}
