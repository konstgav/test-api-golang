package oauth

import (
	"log"
	"net/http"
	"strings"

	jwt "github.com/golang-jwt/jwt/v4"
)

func IsAuthorized(next func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerAuthorization := r.Header["Authorization"]
		if headerAuthorization == nil {
			log.Println("No Authorization Token provided")
			http.Error(w, "No Authorization Token provided", http.StatusUnauthorized)
			return
		}
		tokenString := strings.Split(headerAuthorization[0], "Bearer ")[1]
		secret, err := ReadFileToByteSlice("oauth/secret.pem")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			log.Println("the token is not valid")
			http.Error(w, "the token is not valid", http.StatusUnauthorized)
			return
		}
		iss := "https://login.xsolla.com"
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			log.Println("the issuer is wrong")
			http.Error(w, "the issuer is wrong", http.StatusUnauthorized)
			return
		}
		privilegedGroup := "readers"
		tokenClaims := token.Claims.(jwt.MapClaims)
		userGroups := tokenClaims["groups"].([]interface{})
		checkGroups := false
		for i := range userGroups {
			if userGroups[i].(map[string]interface{})["name"] == privilegedGroup {
				checkGroups = true
				break
			}
		}
		if !checkGroups {
			log.Println("the group is wrong, you have not enough privileges")
			http.Error(w, "the group is wrong, you have not enough privileges", http.StatusUnauthorized)
			return
		}
		next(w, r)
	})
}
