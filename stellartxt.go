package stellarutils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type StellarTxtResponse struct {
	URL  string
	Body string
	Err  error
}

type StellarTxtQueue struct {
	Queue []StellarTxtResponse
}

func (queue *StellarTxtQueue) Add(stellarTxtResponse StellarTxtResponse) (int, error) {
	if len(queue.Queue) == 3 {
		return -1, fmt.Errorf("stellarutils: Could not add StellarTxtResponse we already have the maximum of 3 %v", queue.Queue)
	}
	queue.Queue = append(queue.Queue, stellarTxtResponse)
	return len(queue.Queue) - 1, nil
}

func ResolveFederationURL(domainVariants []string) (string, error) {
	stellarTxt, err := FetchStellarTxt(domainVariants)
	if err != nil {
		return "", err
	}
	federationURL, err := ParseFederationURL(stellarTxt)
	if err != nil {
		return "", err
	}
	return federationURL, nil
}

func FetchStellarTxt(urls []string) (string, error) {
	responseChannel, errorChannel := make(chan StellarTxtResponse), make(chan StellarTxtResponse)

	for _, url := range urls {
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
			fmt.Printf("%v", resp)
		case resp := <-errorChannel:
			// Remove from queue.
			// Check next item in queue for response and return it, otherwise do nothing.
			fmt.Printf("%v", resp)
		}
	}

	return "", nil
}

func ParseFederationURL(stellarTxt string) (string, error) {
	return "", nil
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
