package stellarutils

import "testing"

func TestURLVariants(t *testing.T) {
	domain := "tacoman.com"
	variants := URLVariants(domain)

	if len(variants) != 3 {
		t.Errorf("Should have 3 variants but had: %v", variants)
	}

	if variants[0] != "https://stellar.tacoman.com/stellar.txt" {
		t.Errorf("Should be stellar.%v", domain)
	}

	if variants[1] != "https://tacoman.com/stellar.txt" {
		t.Errorf("Should be %v", domain)
	}

	if variants[2] != "https://www.tacoman.com/stellar.txt" {
		t.Errorf("Should be www.%v", domain)
	}
}
