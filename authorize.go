package authorize

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// Require checks for a valid Authorization token
func Require(endpoint func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := claimsFromRequest(r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			endpoint(w, r)
		}
	}
}

var issuerPublicKeys map[string]string

func claimsFromRequest(r *http.Request) (claims *Claims, err error) {
	authorizationHeader := r.Header.Get("Authorization")

	if len(authorizationHeader) == 0 {
		return nil, fmt.Errorf("No authorization header")
	}

	claims = new(Claims)

	token, err := jwt.ParseWithClaims(authorizationHeader[7:], claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret-key"), nil

	})

	if err != nil {
		return nil, fmt.Errorf("Could not parse claim: %+v", err)
	}

	if _, ok := token.Claims.(*Claims); !ok || !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	return claims, err
}

type Claims struct {
	jwt.StandardClaims

	ExternalID string  `json:"id"`
	Username   string  `json:"username"`
	Fullname   string  `json:"fullname"`
	Groups     []Group `json:"groups"`
}

type Group struct {
	ExternalID  string `json:"id"`
	Name        string `json:"name"`
	FacilityIDs []int  `json:"facility_ids"`
}
