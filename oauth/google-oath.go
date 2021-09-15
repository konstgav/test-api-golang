package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func CreateGoogleAuth() (*GoogleAuth, error) {
	var acc GoogleAccountCredentials
	byteValue, err := ReadFileToByteSlice("oauth/google-account-credentials.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(byteValue, &acc)
	if err != nil {
		fmt.Println("Unmarshall error: ", err.Error())
		return nil, err
	}
	config := &oauth2.Config{
		ClientID:     acc.ClientId,
		ClientSecret: acc.ClientSecret,
		RedirectURL:  acc.RedirectURL,
		Scopes:       acc.Scopes,
		Endpoint:     google.Endpoint,
	}
	return &GoogleAuth{Config: config}, nil
}

type GoogleAccountCredentials struct {
	ClientId     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURL  string   `json:"redirect_URL"`
	Scopes       []string `json:"scopes"`
}

type GoogleAuth struct {
	Config *oauth2.Config
}

func (a GoogleAuth) AuthMiddleware(next func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil || cookie == nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if cookie.Value == "" {
			log.Println("Need to authorizate")
			http.Error(w, "Need to authorizate", http.StatusUnauthorized)
			return
		}
		log.Println("Google access token: ", cookie.Value)
		next(w, r)
	})
}

func (a GoogleAuth) Check(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: check")
}

func (a GoogleAuth) Oauth2callback(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: oauth2callback")

	code, ok := r.URL.Query()["code"]

	if !ok {
		log.Println("Url Param 'code' is missing")
		http.Error(w, "Url Param 'code' is missing", http.StatusInternalServerError)
		return
	}
	ctx := context.Background()

	token, err := a.Config.Exchange(ctx, code[0])
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "error in Exchange code to token "+err.Error(), http.StatusInternalServerError)
		return
	}

	if token == nil {
		log.Println("token is nil")
		http.Error(w, "token is nil", http.StatusInternalServerError)
		return
	}

	log.Println("Success - get Google access token!")
	e, err := json.Marshal(token)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := http.Cookie{
		Name:   "token",
		Value:  string(e),
		MaxAge: 300,
	}
	http.SetCookie(w, &c)
	w.WriteHeader(200)
}

func (a GoogleAuth) Authorize(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: authorize")
	authURL := a.Config.AuthCodeURL("your state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusSeeOther)
}
