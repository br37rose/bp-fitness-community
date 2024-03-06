package controller

import (
	"encoding/base64"
	"strings"
)

// GetBinFromBase64String will convert the string parameter and return a `[]byte` object.
//
// Special thanks:
// (1) https://github.com/tomchristie/django-rest-framework/pull/1268
// (2) https://stackoverflow.com/a/39587386
func GetBinFromBase64String(b64s string) ([]byte, error) {
	if strings.Contains(b64s, "data:") && strings.Contains(b64s, ";base64,") {
		b64s = strings.Split(b64s, ";base64,")[1:][0]
	}

	dec, err := base64.StdEncoding.DecodeString(b64s)
	if err != nil {
		return nil, err
	}
	return dec, err
}

// Base64 encode/decode without padding on golang
// https://stackoverflow.com/q/31971614

func Base64EncodeStrippedFromString(s string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(s))
	return strings.TrimRight(encoded, "=")
}

func Base64EncodeStrippedFromBytes(b []byte) string {
	encoded := base64.StdEncoding.EncodeToString(b)
	return strings.TrimRight(encoded, "=")
}

func Base64DecodeStrippedToString(s string) (string, error) {
	if i := len(s) % 4; i != 0 {
		s += strings.Repeat("=", 4-i)
	}
	decoded, err := base64.StdEncoding.DecodeString(s)
	return string(decoded), err
}

func Base64DecodeStrippedToBytes(s string) ([]byte, error) {
	if i := len(s) % 4; i != 0 {
		s += strings.Repeat("=", 4-i)
	}
	decoded, err := base64.StdEncoding.DecodeString(s)
	return decoded, err
}
