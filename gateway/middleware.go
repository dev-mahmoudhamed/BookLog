package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func BasicLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		status := c.Writer.Status()
		log.Printf("%s %s -> %d", c.Request.Method, c.Request.URL.Path, status)
	}
}

// JWTMiddleware validates token and sets "userID" in context if OK.
// Behavior controlled by env:
// - AUTH_INTROSPECT_URL (if set) -> calls introspection endpoint (POST token=...)
// - AUTH_ALGO (HS256 or RS256) and AUTH_HS_SECRET or AUTH_RS_PUBKEY (PEM) -> local verify
func JWTMiddleware() gin.HandlerFunc {
	introspectURL := os.Getenv("AUTH_INTROSPECT_URL")
	algo := os.Getenv("AUTH_ALGO") // e.g. HS256 or RS256

	hsSecret := os.Getenv("AUTH_HS_SECRET")
	rsPubKeyPEM := os.Getenv("AUTH_RS_PUBKEY") // PEM string

	var rsKey any
	if rsPubKeyPEM != "" {
		// parse RSA public key if provided
		var err error
		rsKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(rsPubKeyPEM))
		if err != nil {
			log.Printf("failed parsing RSA public key: %v", err)
			rsKey = nil
		}
	}

	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}
		parts := strings.Fields(auth)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
			return
		}
		tokenStr := parts[1]

		// Option A: introspection
		if introspectURL != "" {
			active, sub, err := introspectToken(introspectURL, tokenStr)
			if err != nil || !active {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token invalid"})
				return
			}
			// set user id (subject) in context
			c.Set("userID", sub)
			c.Next()
			return
		}

		// Option B: local verification
		// Option B: local verification
		var keyFunc jwt.Keyfunc
		switch algo {
		case "HS256":
			if hsSecret == "" {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "HS256 configured but secret missing"})
				return
			}
			keyFunc = func(token *jwt.Token) (interface{}, error) {
				// CRITICAL: Verify the token algorithm is actually HMAC
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(hsSecret), nil
			}
		case "RS256":
			if rsKey == nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "RS256 configured but public key not available"})
				return
			}
			keyFunc = func(token *jwt.Token) (interface{}, error) {
				// CRITICAL: Verify the token algorithm is actually RSA
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return rsKey, nil
			}
		default:
			// "Auto" mode: Check what the token claims it is, and see if we have a matching key
			keyFunc = func(token *jwt.Token) (interface{}, error) {
				// If token is HMAC (HS256) and we have a secret -> OK
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok && hsSecret != "" {
					return []byte(hsSecret), nil
				}

				// If token is RSA (RS256) and we have a pubkey -> OK
				if _, ok := token.Method.(*jwt.SigningMethodRSA); ok && rsKey != nil {
					return rsKey, nil
				}

				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
		}

		// Parse token
		parsed, err := jwt.Parse(tokenStr, keyFunc)
		if err != nil || !parsed.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token invalid", "detail": err.Error()})
			return
		}
		// try get "sub" or "user_id" claim
		if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
			var sub string
			if v, found := claims["sub"]; found {
				sub = toString(v)
			} else if v, found := claims["user_id"]; found {
				sub = toString(v)
			}
			if sub != "" {
				c.Set("userID", sub)
			}
		}
		c.Next()
	}
}

func toString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	case float64:
		// maybe numeric id
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.0f", x), "0"), ".")
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

// introspectToken calls a token introspection endpoint. Expect JSON { active: bool, sub: "..." }
func introspectToken(url string, token string) (active bool, subject string, err error) {
	form := "token=" + token
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return false, "", fmt.Errorf("introspect failed: %s", string(body))
	}
	var out struct {
		Active bool   `json:"active"`
		Sub    string `json:"sub"`
		// many introspect endpoints return "username" or "user_id" etc
		UserID string `json:"user_id"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return false, "", err
	}
	if out.Sub != "" {
		return out.Active, out.Sub, nil
	}
	if out.UserID != "" {
		return out.Active, out.UserID, nil
	}
	return out.Active, "", nil
}
