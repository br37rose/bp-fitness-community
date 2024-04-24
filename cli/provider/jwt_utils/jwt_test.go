package jwt_utils

import (
	"testing"
	"time"
)

// func GenerateJWTTokenPair(hmacSecret []byte, uuid string, d time.Duration) (string, string, error) {

func TestGenerateJWTTokenPair(t *testing.T) {
	sampleHMACSecret := []byte("123secret")
	sampleUUID := "xxx-xxx-xxx-xxx"
	sampleAccessDuration := 100 * time.Second
	sampleRefreshDuration := sampleAccessDuration + 72*time.Hour

	actualAccessToken, _, actualRefreshToken, _, err := GenerateJWTTokenPair(sampleHMACSecret, sampleUUID, sampleAccessDuration, sampleRefreshDuration)
	if err != nil {
		t.Errorf("received an error %v", err)
	}
	if actualAccessToken == "" {
		t.Error("access token does not exist")
	}
	if actualRefreshToken == "" {
		t.Error("refresh token does not exist")
	}
}

func TestProcessJWTToken(t *testing.T) {
	sampleHMACSecret := []byte("123secret")
	sampleUUID := "xxx-xxx-xxx-xxx"
	sampleAccessDuration := 100 * time.Second
	sampleRefreshDuration := sampleAccessDuration + 72*time.Hour

	actualAccessToken, _, actualRefreshToken, _, err := GenerateJWTTokenPair(sampleHMACSecret, sampleUUID, sampleAccessDuration, sampleRefreshDuration)
	if err != nil {
		t.Errorf("received an error %v", err)
	}
	if actualAccessToken == "" {
		t.Error("access token does not exist")
	}
	if actualRefreshToken == "" {
		t.Error("refresh token does not exist")
	}

	actualUUID, err := ProcessJWTToken(sampleHMACSecret, actualAccessToken)
	if err != nil {
		t.Errorf("received an error %v", err)
	}
	if actualUUID == "" {
		t.Error("access token does not exist")
	}
	if sampleUUID != actualUUID {
		t.Errorf("jwt claim is wrong, got %v but was expecting %v", actualUUID, sampleUUID)
	}
}
