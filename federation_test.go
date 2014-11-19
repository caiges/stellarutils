package stellarutils

import "testing"

func TestDomainVariants(t *testing.T) {
	domain := "tacoman.com"
	variants := DomainVariants(domain)

	if len(variants) != 3 {
		t.Errorf("Should have 3 variants but had: %v", variants)
	}

	if variants[0] != "stellar.tacoman.com" {
		t.Errorf("Should be stellar.%v", domain)
	}

	if variants[1] != "tacoman.com" {
		t.Errorf("Should be %v", domain)
	}

	if variants[2] != "www.tacoman.com" {
		t.Errorf("Should be www.%v", domain)
	}
}
