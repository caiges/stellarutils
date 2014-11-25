package stellarutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type FederationResponse struct {
	FederationJSON struct {
		Type               string `json:"type"`
		Destination        string `json:"destination"`
		Domain             string `json:"domain"`
		DestinationAddress string `json:"destination_address"`
	} `json:"federation_json"`
}

func ResolveFederationUser(user string) FederationResponse {
	userInfo := strings.Split(user, "@")
	username := userInfo[0]
	domain := userInfo[1]
	params := url.Values{}
	stellarTxtDomains := DomainVariants(domain)

	// Append "/stellar.txt" to the end of each variant.
	for i, _ := range stellarTxtDomains {
		stellarTxtDomains[i] += "/stellar.txt"
	}

	federationURL := ResolveFederationURL(stellarTxtDomains)

	// Add required URL params
	params.Add("destination", username)
	params.Add("domain", domain)
	params.Add("type", "federation")

	// Make request to federation service.
	resp, err := http.Get(federationURL + "?" + params.Encode())
	if err != nil {
		fmt.Println("Couldn't query federation service")
	}
	defer resp.Body.Close()

	// Read response into byte array
	federationBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Invalid data from federation service")
	}

	federationResponse := FederationResponse{}

	// Unmarshall the byte array into a FederationResponse type
	err = json.Unmarshal(federationBody, &federationResponse)
	if err != nil {
		fmt.Println("Could not unmarshall federation response")
	}

	return federationResponse
}

func DomainVariants(domain string) []string {
	// Search order specified in https://github.com/stellar/docs/blob/master/docs/Stellar.txt.md
	variants := []string{"stellar." + domain, domain, "www." + domain}
	return variants
}
