package stellarutils

import (
	"bufio"
	"fmt"
	"net/http"
)

type FederationResponse struct {
	FederationJSON struct {
		Type               string `json:"type"`
		Destination        string `json:"destination"`
		Domain             string `json:"domain"`
		DestinationAddress string `json:"destination_address"`
	} `json:"federation_json"`
}

func ResolveFederationURL(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Coulnd't fetch Stellar.txt")
	}
	defer resp.Body.Close()

	var federationURL string

	// Go through response and find federation URL
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "[federation_url]" {
			fmt.Println("Found federation URL...")
			scanner.Scan()
			federationURL = scanner.Text()
		}
	}

	return federationURL
}
