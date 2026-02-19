package utils

import "testing"

func TestGenerateWhatsAppLink_EncodesMessageAndNormalizesNumber(t *testing.T) {
	link, err := GenerateWhatsAppLink("+212 600-000-000", "iPhone 14 Pro")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	expected := "https://wa.me/212600000000?text=Bonjour%20je%20veux%20plus%20d%27information%20sur%20iPhone%2014%20Pro"
	if link != expected {
		t.Fatalf("unexpected link\nexpected: %s\nactual:   %s", expected, link)
	}
}

func TestGenerateWhatsAppLink_InvalidNumber(t *testing.T) {
	_, err := GenerateWhatsAppLink("abc", "iPhone")
	if err == nil {
		t.Fatalf("expected error for invalid number")
	}
}
