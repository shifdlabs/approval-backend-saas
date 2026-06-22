package middleware

import (
	"Microservice/helper"
	"Microservice/pkg/jwks"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWKSAuth validates the incoming SIS-issued JWT against the cached SIS public
// key (no network call per request) and stores the identity claims in the Gin
// context. On a validation failure it refreshes the JWKS once and retries, to
// transparently handle SIS key rotation; if it still fails it returns 401.
func JWKSAuth(client *jwks.JWKSClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractBearer(c)
		if tokenString == "" {
			unauthorized(c, "Missing or malformed Authorization header")
			return
		}

		claims, err := parseAndValidate(client, tokenString)
		if err != nil {
			// Possible key rotation — refresh once, then retry.
			if rErr := client.Refresh(); rErr == nil {
				claims, err = parseAndValidate(client, tokenString)
			}
		}
		if err != nil {
			unauthorized(c, "Invalid or expired token")
			return
		}

		c.Set(helper.ContextKeyUserID, claimString(claims, "sub"))
		c.Set(helper.ContextKeyOrgID, claimString(claims, "org_id"))
		c.Set(helper.ContextKeyOrgRole, claimString(claims, "org_role"))
		c.Set(helper.ContextKeyEmail, claimString(claims, "email"))
		c.Set(helper.ContextKeyName, claimString(claims, "name"))
		c.Set(helper.ContextKeyProducts, claimStringSlice(claims, "products"))

		c.Next()
	}
}

// SubscriptionCheck must run after JWKSAuth. It rejects requests whose token is
// not scoped to an organization, or whose products claim does not include the
// Shifd Approval subscription.
func SubscriptionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if helper.GetOrgID(c) == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Request is not scoped to an organization",
				"code":  "ORG_MISSING",
			})
			return
		}

		if !contains(helper.GetProducts(c), "shifd-approval") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Subscription to Shifd Approval is not active",
				"code":  "SUBSCRIPTION_INACTIVE",
			})
			return
		}

		c.Next()
	}
}

// parseAndValidate verifies the RS256 signature with the cached SIS public key
// and returns the token claims.
func parseAndValidate(client *jwks.JWKSClient, tokenString string) (jwt.MapClaims, error) {
	key := client.GetPublicKey()
	if key == nil {
		return nil, errors.New("no SIS public key available")
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func extractBearer(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	parts := strings.SplitN(header, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func unauthorized(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": msg,
		"code":  "UNAUTHORIZED",
	})
}

func claimString(claims jwt.MapClaims, key string) string {
	if v, ok := claims[key].(string); ok {
		return v
	}
	return ""
}

// claimStringSlice reads a JSON array claim (decoded as []interface{}) into a
// []string, ignoring non-string elements.
func claimStringSlice(claims jwt.MapClaims, key string) []string {
	raw, ok := claims[key].([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func contains(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}
