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
		log.Printf("Failed to load .env: %v", err)
	}

	certFilePath := os.Getenv("CERT_FILE_PATH")
	keyFilePath := os.Getenv("CERT_KEY_PATH")

	// TODO: delete log
	log.Println("certFilePath path:", certFilePath)
	log.Println("keyFilePath path:", keyFilePath)

	cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		log.Printf("Error load certificate: %v", err)
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
