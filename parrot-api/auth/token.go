package auth

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/parrot/parrot-api/model"
)

// TokenProvider holds the Auth Provider's name and Signing Key.
type TokenProvider struct {
	Name       string
	SigningKey []byte
}

// AuthStore is the interface that an Auth Provider implementation requires to retrieve
// and validate credentials in order to issue a token.
type AuthStore interface {
	model.UserStorer
	model.ProjectClientStorer
	Ping() error
	Close() error
}

// CreateToken creates a new token with the provided claims and signs it.
func (p *TokenProvider) CreateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.SigningKey)
}

// ParseAndVerifyToken parses, verifies and validates the claims of the token.
// It returns the claims or an error if the token is not valid or if something went wrong.
func (p *TokenProvider) ParseAndVerifyToken(tokenString string) (jwt.MapClaims, error) {
	claims, err := p.ParseAndExtractClaims(tokenString)
	if err != nil {
		return nil, err
	}
	if err := claims.Valid(); err != nil {
		return nil, err
	}

	return claims, nil
}

// ParseAndExtractClaims parses the claims of the token and its signature without validating the claims.
// It returns the claims or an error.
func (p *TokenProvider) ParseAndExtractClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := parseToken(tokenString, p.SigningKey)
	if err != nil {
		return nil, err
	}

	claims, err := extractClaims(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// parseToken parses and verifies the signature of the token.
// It returns the parsed token or an error.
func parseToken(tokenString string, signingKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}

// extractClaims extracts the claims from the parsed token.
func extractClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
