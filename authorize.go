package authorize

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// Require is a http.HandlerFunc decorator that checks for a valid web token
func Require(endpoint func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if len(authorization) == 0 {
			http.Error(w, "no authorization header", http.StatusUnauthorized)
			return
		}

		if err := authorize(authorization[7:]); err != nil {
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
