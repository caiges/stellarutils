package stellarutils

import (
	"bufio"
	"fmt"
	"net/http"
)

// Resolve the federation URL from the Stellar.txt file.
func ResolveFederationURL(url string) string {
	// Fetch Stellar.txt
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
