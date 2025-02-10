package network

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// If want to use this need to enable volumes in docker-compose.production.yml
// volumes:
// - ${HOST_CERT_DIRECTORY}:${APP_CERT_DIRECTORY}:ro // Not work if use with Coolify.
// - /data/coolify/proxy/caddy/data/caddy/certificates/acme-v02.api.letsencrypt.org-directory/shiftwave-dev-b.mijio.app:/app/ssl/certs:ro
func LoadCertificate() tls.Certificate {
	_, err := godotenv.Read(".env")
	if err != nil {
		log.Printf("Failed to load .env: %v", err)
	}

	certFilePath := os.Getenv("CERT_FILE_PATH")
	keyFilePath := os.Getenv("CERT_KEY_PATH")

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
