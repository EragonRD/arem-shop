package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var waDigitsRegex = regexp.MustCompile(`^[1-9][0-9]{7,14}$`)

// GenerateWhatsAppLink construit le lien de contact WhatsApp public d'un produit.
// Le message est encode avec des espaces en %20 (et non +) selon le format attendu.
func GenerateWhatsAppLink(whatsAppNumber, productName string) (string, error) {
	number, err := normalizeAndValidateWhatsAppNumber(whatsAppNumber)
	if err != nil {
		return "", err
	}

	message := fmt.Sprintf("Bonjour je veux plus d'information sur %s", strings.TrimSpace(productName))
	encoded := encodeWithPercentSpaces(message)

	return fmt.Sprintf("https://wa.me/%s?text=%s", number, encoded), nil
}

func normalizeAndValidateWhatsAppNumber(number string) (string, error) {
	trimmed := strings.TrimSpace(number)
	trimmed = strings.ReplaceAll(trimmed, " ", "")
	trimmed = strings.ReplaceAll(trimmed, "-", "")
	trimmed = strings.ReplaceAll(trimmed, "(", "")
	trimmed = strings.ReplaceAll(trimmed, ")", "")
	trimmed = strings.TrimPrefix(trimmed, "+")

	if !waDigitsRegex.MatchString(trimmed) {
		return "", fmt.Errorf("invalid whatsapp number")
	}

	return trimmed, nil
}

func encodeWithPercentSpaces(value string) string {
	encoded := url.QueryEscape(value)
	return strings.ReplaceAll(encoded, "+", "%20")
}
