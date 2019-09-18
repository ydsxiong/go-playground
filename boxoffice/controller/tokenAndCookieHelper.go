package controller

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("box_office_secret_key")

type Claims struct {
	GuestName string `json:"username"`
	jwt.StandardClaims
}

/*
 * generate a token with an expiry from a new claim
 */
func generateToken(claims *Claims) (*string, *time.Time) {
	expirationTime := time.Now().Add(reservationTime)
	claims.ExpiresAt = expirationTime.Unix()
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, nil
	}
	return &tokenString, &expirationTime
}

/*
 * extract a claim out of an alive and authentic token
 */
func validateToken(tokenCookie *http.Cookie) *Claims {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenCookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	//err may be a jwt.ErrSignatureInvalid one, or the token could be an invalid one
	if err != nil || !tkn.Valid {
		return nil
	}
	return claims
}

func clearoutTokenFromCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now(),
	})
}

func pushTokenIntoClientCookie(w http.ResponseWriter, tokenString *string, expiry *time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   *tokenString,
		Expires: *expiry,
	})
}

func pullTokenFromClientCookie(req *http.Request) *http.Cookie {
	// try to obtain the session token from the requests cookies, which comes with every request
	tokenCookie, err := req.Cookie("token")
	if err != nil || tokenCookie == nil {
		return nil
	}
	return tokenCookie
}

func pushValidClaimIntoContext(r *http.Request, claims *Claims) {
	context.Set(r, "claimFromCookie", claims)
}

func retriveValidClaimsFromContext(r *http.Request) *Claims {
	cookieclaims := context.Get(r, "claimFromCookie")
	claims, ok := cookieclaims.(*Claims)
	if ok {
		return claims
	}
	return nil
}

func pushSessionTimeoutCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:  "sessionexpired",
		Value: "true",
	})
}

func checkAndRemoveCookieForSessionTimeout(w http.ResponseWriter, r *http.Request) string {
	yes, err := r.Cookie("sessionexpired")
	if err == nil && "true" == yes.Value {
		http.SetCookie(w, &http.Cookie{
			Name:  "sessionexpired",
			Value: "",
		})
		return "You session was timed out, please start over and try again!"
	}
	return ""
}

func pushUnavailableCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:  "unavailable",
		Value: "true",
	})
}

func checkAndRemoveCookieForUnavailable(w http.ResponseWriter, r *http.Request) string {
	yes, err := r.Cookie("unavailable")
	if err == nil && "true" == yes.Value {
		http.SetCookie(w, &http.Cookie{
			Name:  "unavailable",
			Value: "",
		})
		return "Sorry, it's too late, tickets have just been sold out!"
	}
	return ""
}
