package stellarutils

import (
	"errors"
	"fmt"
)

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
