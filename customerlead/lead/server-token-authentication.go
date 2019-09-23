package lead

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// Create the JWT key used to generate signature for a system known user
var jwtKey = []byte("secret_key")

// a shared api key for all api users, or can be unique, one for each user stored in the db
var acceptableUsers = []string{"apiuser"}

type claims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

func (c *claims) isValid() bool {
	for _, user := range acceptableUsers {
		if c.User == user {
			return true
		}
	}
	return false
}

/*
  generate a token for a claim of a specific user
  in this case the shared one: apiUser
*/
func generateApiUserToken(claim *claims) *string {
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil
	}
	return &tokenString
}

func tokenAuthoringMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if doAuth(w, req) {
			next(w, req)
		} else {
			respondError(w, http.StatusForbidden, "A valid access token could not be found found in the request header!")
		}
	})
}

func doAuth(w http.ResponseWriter, req *http.Request) bool {
	token := pullTokenFromRequest(req)
	if token == "" {
		return false
	}
	claims := validateToken(token)
	if claims == nil {
		return false
	}

	// in case the underlying api call service may need to access to a user identity from the valid claim
	pushValidClaimIntoContext(req, claims)
	return true
}

func pullTokenFromRequest(req *http.Request) string {
	// try to obtain the session token from the requests header, which should come with every api request
	var headerToken = req.Header.Get("x-access-token")

	headerToken = strings.TrimSpace(headerToken)

	return headerToken
}

/*
 * extract a claim out of an authentic token
 */
func validateToken(token string) *claims {
	claim := &claims{}

	tkn, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	//err may be a jwt.ErrSignatureInvalid one, or the token could be an invalid one
	if err != nil || !tkn.Valid {
		return nil
	}
	return claim
}

/**
push and pull the authorized claim via request context through chained handlers
*/
func pushValidClaimIntoContext(r *http.Request, claim *claims) {
	context.Set(r, "claimFromXAccessToken", claim)
}

func retriveValidClaimsFromContext(r *http.Request) *claims {
	cookieclaims := context.Get(r, "claimFromXAccessToken")
	claims, ok := cookieclaims.(*claims)
	if ok {
		return claims
	}
	return nil
}
