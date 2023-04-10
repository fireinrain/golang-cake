package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func main() {
	// Replace these with the hostname and port of the server you want to check
	hostname := "www.metrodeal.com"
	port := "443"

	tlsConfig := &tls.Config{
		// Only enable TLSv1.3
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,
	}

	// Create an HTTP client with the custom TLS config
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Make a request to the server
	resp, err := client.Get(fmt.Sprintf("https://%s:%s/", hostname, port))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Check if the server supports H2
	if resp.ProtoMajor == 2 {
		fmt.Println("Server supports HTTP/2 (H2)")
	} else {
		fmt.Println("Server does not support HTTP/2 (H2)")
	}

	// Check if the server supports TLSv1.3
	if resp.TLS.Version == tls.VersionTLS13 {
		fmt.Println("Server supports TLSv1.3")
	} else {
		fmt.Println("Server does not support TLSv1.3")
	}
}
