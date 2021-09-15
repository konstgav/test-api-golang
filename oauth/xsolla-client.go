package oauth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func PostRequestToProductApp(w http.ResponseWriter, r *http.Request) {
	token, err := GetXsollaLoginToken()
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = PostRequestToProductAppWithToken(w, token)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func PostRequestToProductAppWithToken(w http.ResponseWriter, token *TokenInfo) error {
	uri := "http://localhost:8080/product"
	byteValue, err := ReadFileToByteSlice("oauth/product-test-post-request.json")
	if err != nil {
		return err
	}
	payload := strings.NewReader(string(byteValue))
	req, err := http.NewRequest("POST", uri, payload)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Println(string(body))
	w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", res.Header.Get("Content-Length"))
	_, err = w.Write(body)
	if err != nil {
		return err
	}
	return nil
}

//TODO: use oauth2.config
func GetXsollaLoginToken() (*TokenInfo, error) {
	state := "statestate"

	byteValue, err := ReadFileToByteSlice("oauth/xsolla-login-account-credentials.json")
	if err != nil {
		return nil, err
	}
	var acc XsollaLoginAccountCredentials
	err = json.Unmarshal(byteValue, &acc)
	if err != nil {
		return nil, err
	}

	userCredentialsByteValue, err := ReadFileToByteSlice("oauth/xsolla-login-user-credentials.json")
	if err != nil {
		return nil, err
	}

	uri := "https://login.xsolla.com/api/oauth2/login?state=" + state + "&client_id=" + acc.ClientId + "&response_type=code"

	payload := strings.NewReader(string(userCredentialsByteValue))
	req, _ := http.NewRequest("POST", uri, payload)
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var loginUrl LoginResponse
	err = json.Unmarshal(body, &loginUrl)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(loginUrl.LoginUrl)
	if err != nil {
		return nil, err
	}
	values := u.Query()
	code := values["code"][0]

	uri = "https://login.xsolla.com/api/oauth2/token"

	payload = strings.NewReader("grant_type=authorization_code&client_secret=" + acc.ClientSecret + "&client_id=" + acc.ClientId + "&redirect_uri=http://localhost:8080/oauth2callbackxsolla&code=" + code) // strings.NewReader(requestBody)
	req, _ = http.NewRequest("POST", uri, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var tokenInfo TokenInfo
	err = json.Unmarshal(body, &tokenInfo)
	if err != nil {
		return nil, err
	}
	return &tokenInfo, nil
}

func ReadFileToByteSlice(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

type XsollaLoginAccountCredentials struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type TokenInfo struct {
	AccessToken  string `json:"access_token"`
	ExpireIn     int    `json:"expire_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type LoginResponse struct {
	LoginUrl string `json:"login_url"`
}
