package main

import (
	"fmt"
	"net/http"

	jwtmiddleware "github.com/Akagi201/middleware/jwt"
	jwt "github.com/dgrijalva/jwt-go"
)

var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*jwt.Token)
	fmt.Printf("user: %+v", user)
	fmt.Fprint(w, "This is an authenticated request")
	fmt.Fprint(w, "Claim content:\n")
	for k, v := range user.Claims.(jwt.MapClaims) {
		fmt.Fprintf(w, "%s :\t%#v\n", k, v)
	}
})

func main() {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		Debug: true,
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("supersecret"), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	app := jwtMiddleware.Handler(myHandler)
	http.ListenAndServe("0.0.0.0:3000", app)
}
