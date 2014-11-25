package stellarutils

import "testing"

func TestStellarTxtQueueAdd(t *testing.T) {
	queue := StellarTxtQueue{}
	index, err := queue.Add(StellarTxtResponse{URL: "1", Body: "blars"})
	if err != nil {
		t.Error(err)
	}

	if index > 0 {
		t.Errorf("Should have 1 item in the queue but had: %v", queue.Queue)
	}
}

func TestStellarTxtQueueAddOverCapacity(t *testing.T) {
	queue := StellarTxtQueue{}
	queue.Add(StellarTxtResponse{URL: "1", Body: "blars"})
	queue.Add(StellarTxtResponse{URL: "2", Body: "taco"})
	queue.Add(StellarTxtResponse{URL: "3", Body: "man"})

	_, err := queue.Add(StellarTxtResponse{URL: "4", Body: "eats"})
	if err == nil {
		t.Errorf("Should have returned an error: %v", err)
	}

	if len(queue.Queue) > 3 {
		t.Errorf("Should have had 3 items in the queue but had: %v", queue.Queue)
	}
}
