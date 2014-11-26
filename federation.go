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

func ResolveFederationUser(user string) (*FederationResponse, error) {
	userInfo := strings.Split(user, "@")
	username := userInfo[0]
	domain := userInfo[1]

	// Add required URL params
	params := url.Values{}
	params.Add("destination", username)
	params.Add("domain", domain)
	params.Add("type", "federation")

	stellarTxtURLs := URLVariants(domain)

	federationURL, err := ResolveFederationURL(stellarTxtURLs)

	federationResponse, err := QueryFederationService(federationURL, params)
	if err != nil {
		return nil, err
	}

	return federationResponse, nil
}

func QueryFederationService(url string, params url.Values) (*FederationResponse, error) {
	// Make request to federation service.
	resp, err := http.Get(url + "?" + params.Encode())
	if err != nil {
		fmt.Println("Couldn't query federation service")
		return nil, err
	}
	defer resp.Body.Close()

	// Read response into byte array
	federationBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Invalid data from federation service")
		return nil, err
	}

	federationResponse := FederationResponse{}

	// Unmarshall the byte array into a FederationResponse type
	err = json.Unmarshal(federationBody, &federationResponse)
	if err != nil {
		fmt.Println("Could not unmarshall federation response")
	}
	return &federationResponse, nil
}

func URLVariants(domain string) []string {
	// Search order specified in https://github.com/stellar/docs/blob/master/docs/Stellar.txt.md
	variants := []string{"https://stellar." + domain + "/stellar.txt", "https://" + domain + "/stellar.txt", "https://www." + domain + "/stellar.txt"}
	return variants
}
