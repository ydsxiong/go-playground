package controller

import (
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/go-kit/kit/log"
)

type httpHandlerMiddleware func(next http.HandlerFunc) http.HandlerFunc

func LoggingMiddleware(logger log.Logger) httpHandlerMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			defer func(begin time.Time) {
				logger.Log("http_call", runtime.FuncForPC(reflect.ValueOf(next).Pointer()).Name(), "took", time.Since(begin))
			}(time.Now())
			next(w, req)
		}
	}
}

/*
 *  a middle layer service to protect those endpoints that may require authentication
 */
func CreateTokenAuthoringMiddleWare(doAuth func(w http.ResponseWriter, req *http.Request) bool) httpHandlerMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			if doAuth(w, req) {
				next(w, req)
			}
		}
	}
}

/*
 *  or:
 *   a service middleware wrapper to protect those endpoints that may require authentication
 *   for example:
 *  	router.HandleFunc("/charge", tokenAuthoringMiddleWare(PayWithCard)).Methods("POST")

func tokenAuthoringMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		tokenCookie := pullTokenFromClientCookie(req)
		if tokenCookie == nil {
			pushSessionTimeoutCookie(w)
			redirectBackHome(w, req)
			return
		}
		claims := validateToken(tokenCookie)
		if claims == nil {
			redirectBackHome(w, req)
		} else {
			pushValidClaimIntoContext(req, claims)
			next(w, req)
		}
	})
}
*/
