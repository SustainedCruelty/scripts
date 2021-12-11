package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

func main() {

	var client_id string
	var secret_key string
	var scopes string

	// Ask the user for all the values that are needed
	fmt.Printf("[?] Please enter your application's client id: ")
	fmt.Scanln(&client_id)

	fmt.Printf("[?] Please enter your application's secret key: ")
	fmt.Scanln(&secret_key)

	fmt.Printf("[?] Please enter the scopes you want to auth: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		scopes = scanner.Text()
	}
	scopes = strings.Replace(scopes, " ", "%20", -1)

	// Generate a random string for the state parameter in the login url
	state := RandomString(16)

	// Construct the URL we'll be logging in through
	base_url := "https://login.eveonline.com/v2/oauth/authorize/?response_type=code&redirect_uri=http%3A%2F%2Flocalhost%2Foauth-callback"
	login_url := base_url + fmt.Sprintf("&client_id=%s&scope=%s&state=%s", client_id, scopes, state)

	fmt.Printf("\n[*] Login using the following URL: %s\n\n", login_url)

	var auth_code []string
	var req_state []string

	// Start a local server the EVE SSO will make a request to
	m := http.NewServeMux()
	s := http.Server{Addr: ":80", Handler: m}
	m.HandleFunc("/oauth-callback", func(w http.ResponseWriter, r *http.Request) {

		// Extract the query string parameters from the url
		auth_code, _ = r.URL.Query()["code"]
		req_state, _ = r.URL.Query()["state"]
		fmt.Fprintf(w, "Hi there, thanks for logging in. \nCheck your terminal for the access and refresh token!")

		// Stop the server once we receive the request
		go func() {
			if err := s.Shutdown(context.Background()); err != nil {
				log.Fatal(err)
			}
		}()

	})
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Construct the authentication header needed to retrieve the auth token
	auth_header := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", client_id, secret_key)))
	tokens := FetchTokens(auth_header, auth_code[0])

	fmt.Printf("ACCESS TOKEN: \t%s\n\n", tokens.AccessToken)
	fmt.Printf("REFRESH TOKEN: \t%s\n", tokens.RefreshToken)

}

// Generates a random string
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// Fetches the access and refresh token from the EVE ESI
func FetchTokens(header string, code string) TokenResponse {

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)

	client := &http.Client{}

	request, err := http.NewRequest("POST", "https://login.eveonline.com/v2/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Authorization", "Basic "+header)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Host", "login.eveonline.com")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatal("The HTTP response status code isn't 200 :(")
	}
	// Deserializing the response
	var result TokenResponse
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(data, &result)

	return result
}
