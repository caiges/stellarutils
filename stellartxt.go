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
	stellarTxt, err := FetchStellarTxts(urls)
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

func FetchSingleStellarTxt(url string) <-chan StellarTxtResponse {
	responseChannel := make(chan StellarTxtResponse)
	go func() {
		// Fetch the URL.
		body, err := fetch(url)
		if err != nil {
			responseChannel <- StellarTxtResponse{URL: url, Err: err}
			return
		}
		responseChannel <- StellarTxtResponse{URL: url, Body: body}
	}()
	return responseChannel
}

func fanIn(queue StellarTxtQueue) <-chan StellarTxtResponse {
	responsesChannel := make(chan StellarTxtResponse)
	for _, value := range queue.Queue {
		go func(url string) {
			responsesChannel <- <-FetchSingleStellarTxt(url)
		}(value.URL)
	}
	return responsesChannel
}

func FetchStellarTxts(urls []string) (string, error) {
	responseQueue := StellarTxtQueue{}

	for _, url := range urls {
		// Queue up the responses.
		responseQueue.Add(StellarTxtResponse{URL: url})
	}

	responsesChannel := fanIn(responseQueue)

	for {
		select {
		case resp := <-responsesChannel:
			if resp.Body != "" {
				// Set response on queue item. If response satisfies index 0 return it.
				response, i, err := responseQueue.SetResult(resp.URL, resp.Body)
				if err == nil && response != nil && i == 0 {
					return response.Body, nil
				}
			}

			if resp.Err != nil {
				// Remove from queue.
				responseQueue.Remove(resp.URL)
				// Check next item in queue for response and return it, otherwise do nothing.
				next := responseQueue.Head()
				if next != nil && next.Body != "" {
					return next.Body, nil
				}
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
