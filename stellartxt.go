package stellarutils

import (
	"bufio"
	"bytes"
	"fmt"
)

type StellarTxtResponse struct {
	URL  string
	Body string
	Err  error
}

func ResolveFederationURL(urls []string) (string, error) {
	stellarTxt, err := FetchStellarTxt(urls)
	if err != nil {
		return "", err
	}
	federationURL, err := ParseFederationURL(stellarTxt)
	fmt.Printf("%v", federationURL)
	if err != nil {
		return "", err
	}
	return federationURL, nil
}

func FetchStellarTxt(urls []string) (string, error) {
	responseQueue := StellarTxtQueue{}
	responseChannel, errorChannel := make(chan StellarTxtResponse), make(chan StellarTxtResponse)

	for _, url := range urls {
		// Queue up the responses.
		responseQueue.Add(StellarTxtResponse{URL: url})
		// Launch a goroutine to fetch the URL.
		go func(url string) {
			// Fetch the URL.
			body, err := fetch(url)
			if err != nil {
				errorChannel <- StellarTxtResponse{URL: url, Err: err}
				return
			}
			responseChannel <- StellarTxtResponse{URL: url, Body: body}
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		select {
		case resp := <-responseChannel:
			// Set response on queue item. If response satisfies index 0 return it.
			response, i, err := responseQueue.SetResult(resp.URL, resp.Body)
			if err == nil && response != nil && i == 0 {
				return response.Body, nil
			}
		case resp := <-errorChannel:
			// Remove from queue.
			responseQueue.Remove(resp.URL)
			// Check next item in queue for response and return it, otherwise do nothing.
			next := responseQueue.Head()
			if next != nil && next.Body != "" {
				return next.Body, nil
			}
		}
	}

	return "", nil
}

func ParseFederationURL(stellarTxt string) (string, error) {
	var federationURL string
	// Go through response and find federation URL
	scanner := bufio.NewScanner(bytes.NewBufferString(stellarTxt))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "[federation_url]" {
			fmt.Println("Found federation URL...")
			scanner.Scan()
			federationURL = scanner.Text()
			return federationURL, nil
		}
	}
	return "", nil
}
