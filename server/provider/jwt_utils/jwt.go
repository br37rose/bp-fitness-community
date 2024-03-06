package jwt_utils

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// GenerateJWTToken Generate the `access token` for the secret key.
func GenerateJWTToken(hmacSecret []byte, uuid string, ad time.Duration) (string, time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expiresIn := time.Now().Add(ad)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_uuid"] = uuid
	claims["exp"] = expiresIn.Unix()

	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		return "", expiresIn, err
	}

	return tokenString, expiresIn, nil
}

// GenerateJWTTokenPair Generate the `access token` and `refresh token` for the secret key.
func GenerateJWTTokenPair(hmacSecret []byte, uuid string, ad time.Duration, rd time.Duration) (string, time.Time, string, time.Time, error) {
	//
	// Generate token.
	//
	token := jwt.New(jwt.SigningMethodHS256)
	expiresIn := time.Now().Add(ad)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_uuid"] = uuid
	claims["exp"] = expiresIn.Unix()

	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		return "", time.Now(), "", time.Now(), err
	}

	//
	// Generate refresh token.
	//
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshExpiresIn := time.Now().Add(rd)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["session_uuid"] = uuid
	rtClaims["exp"] = refreshExpiresIn.Unix()

	refreshTokenString, err := refreshToken.SignedString(hmacSecret)
	if err != nil {
		return "", time.Now(), "", time.Now(), err
	}

	return tokenString, expiresIn, refreshTokenString, refreshExpiresIn, nil
}

// ProcessJWTToken validates either the `access token` or `refresh token` and returns either the `uuid` if success or error on failure.
func ProcessJWTToken(hmacSecret []byte, reqToken string) (string, error) {
	token, err := jwt.Parse(reqToken, func(t *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})
	if err == nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			uuid := claims["session_uuid"].(string)
			// m["exp"] := string(claims["exp"].(float64))
			return uuid, nil
		}
		return "", err
	}
	return "", err
}
