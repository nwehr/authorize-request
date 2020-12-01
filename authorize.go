package authorize

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// Require is a http.HandlerFunc decorator that checks for a valid web token
func Require(endpoint http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := getToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err := authorize(token); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		endpoint(w, r)
	}
}

// SetAuthorizeFunc sets the function that is used to authorize the token
func SetAuthorizeFunc(fn func(tokenString string) error) {
	authorize = fn
}

// SetKeyFunc sets the function that returns the key to be used for validating the token
func SetKeyFunc(fn func(token *jwt.Token) (interface{}, error)) {
	key = fn
}

var authorize = func(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, key)
	if err != nil {
		return fmt.Errorf("could not parse claim: %v", err)
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

var key = func(token *jwt.Token) (interface{}, error) {
	return []byte("secret-key"), nil
}

func getToken(r *http.Request) (string, error) {
	if authHeader := r.Header.Get("Authorization"); len(authHeader) > 0 {
		return authHeader[7:], nil
	}

	if c, err := r.Cookie("Authorization"); err == nil {
		return c.Value, nil
	}

	return "", fmt.Errorf("no auth token")
}
