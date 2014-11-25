package stellarutils

import (
	"errors"
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
		return -1, fmt.Errorf("Could not add StellarTxtResponse we already have the maximum of 3 %v", queue.Queue)
	}

	if queue.Exists(stellarTxtResponse.URL) {
		return -1, fmt.Errorf("Response already exists in queue: %v", queue.Queue)
	}

	queue.Queue = append(queue.Queue, stellarTxtResponse)
	return len(queue.Queue) - 1, nil
}

func (queue *StellarTxtQueue) Remove(url string) (*StellarTxtResponse, error) {
	for i, value := range queue.Queue {
		if value.URL == url {
			newQueue := queue.Queue[0:i]

			if len(queue.Queue)-1 > i {
				queue.Queue = append(queue.Queue, newQueue...)
				queue.Queue = append(queue.Queue, queue.Queue[i+1:len(queue.Queue)]...)
			} else {
				queue.Queue = newQueue
			}

			return &value, nil
		}
	}

	return nil, nil
}

func (queue *StellarTxtQueue) Head() *StellarTxtResponse {
	if len(queue.Queue) > 0 {
		return &queue.Queue[0]
	}
	return nil
}

func (queue *StellarTxtQueue) Exists(url string) bool {
	for _, value := range queue.Queue {
		if url == value.URL {
			return true
		}
	}
	return false
}

func (queue *StellarTxtQueue) Get(url string) (*StellarTxtResponse, error) {
	for i, value := range queue.Queue {
		if value.URL == url {
			return &queue.Queue[i], nil
		}
	}

	return nil, errors.New("Item not found")
}

func (queue *StellarTxtQueue) SetResult(url string, body string) (*StellarTxtResponse, int, error) {
	response, err := queue.Get(url)
	if err != nil {
		return nil, -1, err
	}

	response.Body = body
	return response, -1, nil
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
			//fmt.Printf("%v", resp)
		case resp := <-errorChannel:
			// Remove from queue.
			responseQueue.Remove(resp.URL)
			// Check next item in queue for response and return it, otherwise do nothing.
			next := responseQueue.Head()
			if next != nil && next.Body != "" {
				return next.Body, nil
			}
			//fmt.Printf("%v", resp)
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
