package stellarutils

import "testing"

func TestAdd(t *testing.T) {
	queue := StellarTxtQueue{}
	index, err := queue.Add(StellarTxtResponse{URL: "1", Body: "blars"})
	if err != nil {
		t.Error(err)
	}

	if index > 0 {
		t.Errorf("Should have 1 item in the queue but had: %v", queue.Queue)
	}
}

func TestAddDuplicate(t *testing.T) {
	queue := StellarTxtQueue{}
	queue.Add(StellarTxtResponse{URL: "1", Body: "blars"})
	_, err := queue.Add(StellarTxtResponse{URL: "1", Body: "blars"})
	if err == nil {
		t.Errorf("Should have returned an error")
	}
}

func TestExists(t *testing.T) {
	queue := StellarTxtQueue{}
	queue.Add(StellarTxtResponse{URL: "1", Body: "blars"})

	exists := queue.Exists("1")

	if !exists {
		t.Errorf("Should have returned true but had: %v", exists)
	}
}

func TestAddOverCapacity(t *testing.T) {
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

func TestGet(t *testing.T) {
	queue := StellarTxtQueue{}
	queue.Add(StellarTxtResponse{URL: "1", Body: "tacos"})
	response, err := queue.Get("1")
	if err != nil {
		t.Error(err)
	}

	if response.Body != "tacos" {
		t.Errorf("Should have returned correct body but had: %v", response.Body)
	}

	response, err = queue.Get("tacoshrimp")
	if err == nil {
		t.Errorf("Error should have been returned")
	}
}

func TestStellarTxtQueueSetResult(t *testing.T) {
	queue := StellarTxtQueue{}
	queue.Add(StellarTxtResponse{URL: "1"})
	queue.SetResult("1", "eats tacos")
	response, err := queue.Get("1")

	if err != nil {
		t.Error(err)
	}

	if response.Body != "eats tacos" {
		t.Errorf("Body should have been set but had: %v", response)
	}
}
