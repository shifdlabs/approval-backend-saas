package config

import (
	"Microservice/model"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenDetails struct {
	Token *string
	// Identifier string
	Email     string
	UserID    string
	ExpiresIn *int64
}

func CreateAccessToken(user *model.User, ttl time.Duration, privateKey string) (*TokenDetails, error) {
	now := time.Now().UTC()
	td := &TokenDetails{
		ExpiresIn: new(int64),
		Token:     new(string),
	}
	*td.ExpiresIn = now.Add(ttl).Unix()
	// td.Identifier = uuid.NewV4().String()
	td.UserID = user.ID.String()
	td.Email = user.Email

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		// msg := "Could not decode token private key"
		return nil, err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		// msg := "Could not parse token private key"
		return nil, err
	}

	atClaims := &jwt.MapClaims{
		"data": map[string]string{
			"email": user.Email,
			"id":    td.UserID,
		},
		"exp": td.ExpiresIn,
		"iat": now.Unix(),
		"iss": "expire.com",
	}

	*td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(key)
	if err != nil {
		// msg := "Could not sign token"
		return nil, err
	}

	return td, nil
}

func CreateRefreshToken(user *model.User, ttl time.Duration, privateKey string) (*TokenDetails, error) {
	now := time.Now().UTC()
	td := &TokenDetails{
		ExpiresIn: new(int64),
		Token:     new(string),
	}
	*td.ExpiresIn = now.Add(ttl).Unix()
	td.UserID = user.ID.String()

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		// msg := "Could not decode token private key"
		return nil, err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		// msg := "Could not parse token private key"
		return nil, err
	}

	atClaims := &jwt.MapClaims{
		"data": map[string]string{
			"id": td.UserID,
		},
		"exp": td.ExpiresIn,
		"iat": now.Unix(),
		"iss": "alpha.com",
	}

	*td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(key)
	if err != nil {
		// msg := "Could not sign token"
		return nil, err
	}

	return td, nil
}
