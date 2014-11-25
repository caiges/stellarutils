package stellarutils

import (
	"io/ioutil"
	"net/http"
	"sync"
)

// Resolve the federation URL from the Stellar.txt file.
/*func ResolveFederationURL(url string) string {
	// Fetch Stellar.txt
	resp, err := http.Get("https://" + url + "/stellar.txt")
	if err != nil {
		fmt.Println("Couldn't fetch Stellar.txt")
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
}*/

type StellarTxtResponse struct {
	URL  string
	Body string
	Err  error
}

type StellarTxtQueue struct {
	Queue struct {
		sync.RWMutex
		Data map[string]string
	}
}

func ResolveFederationURL(domainVariants) (string, error) {
	stellarTxt, err := FetchStellarTxt(domainVariants)
	if err != nil {
		return err
	}
	federationURL, err := ParseFederationURL(stellarTxt)
	if err != nil {
		return err
	}
	return federationURL
}

func FetchStellarTxt(urls []string) (string, error) {
	responseChannel, errorChannel := make(chan StellarTxtResponse), make(chan StellarTxtResponse)

	for _, url := range urls {
		// Launch a goroutine to fetch the URL.
		go func(url string) {
			// Fetch the URL.
			body, err := fetch(url)
			if err != nil {
				errorChannel <- StellarTxtError{URL: url, Err: err}
				return
			}
			responseChannel <- StellarTxtResponse{URL: url, Body: body}
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		select {
		case resp := <-responseChannel:
			// Set response on queue item. If response satisfies index 0 return it.
		case resp := <-errorChannel:
			// Remove from queue.
			// Check next item in queue for response and return it, otherwise do nothing.
		}
	}
}

func fetch(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}
	return string(body), nil
}
