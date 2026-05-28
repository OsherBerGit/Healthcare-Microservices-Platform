package main

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwksURL = "http://localhost:8180/realms/HealthcareRealm/protocol/openid-connect/certs"

func verifyJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization Header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return getPublicKeyFromKeycloak(token)
		})

		if err != nil || !token.Valid {
			log.Printf("Token verification failed: %v", err)
			http.Error(w, "Invalid or expired Token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "user_id", claims["sub"])
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		http.Error(w, "Invalid Token Claims", http.StatusUnauthorized)
	})
}

func getPublicKeyFromKeycloak(token *jwt.Token) (*rsa.PublicKey, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []struct {
			Kid string `json:"kid"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("kid header not found")
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			nBytes, _ := base64.RawURLEncoding.DecodeString(key.N)
			eBytes, _ := base64.RawURLEncoding.DecodeString(key.E)

			n := new(big.Int).SetBytes(nBytes)
			e := int(new(big.Int).SetBytes(eBytes).Int64())

			return &rsa.PublicKey{N: n, E: e}, nil
		}
	}

	return nil, errors.New("public key not found for kid")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming Request: [Method: %s] [Path: %s]", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		next.ServeHTTP(w, r)
	})
}

func idempotencyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idempotencyKey := r.Header.Get("Idempotency-Key")

		if idempotencyKey == "" {
			http.Error(w, "Missing required header: Idempotency-Key", http.StatusBadRequest)
			return
		}

		log.Printf("Received request with Idempotency-Key: %s", idempotencyKey)
		next.ServeHTTP(w, r)
	})
}
