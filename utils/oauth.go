// OAuth management utility functions
package utils

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"net/http"
	"os"
	"path/filepath"
)

type OauthProvider struct {
	name     string
	clientID string
	// https://developers.google.com/identity/protocols/oauth2/#installed
	clientSec string
	scopes    []string
	endpoint  oauth2.Endpoint
}

var (
	port          string = "50135"
	redirectURL   string = fmt.Sprintf("http://localhost:%s/callback", port)
	tokenFileName string = ".rdv.json"
	TokenFilePath string = filepath.Join(getUserHome(), tokenFileName)
)

var Providers = []OauthProvider{
	{
		name:      "gdrive",
		clientID:  "593200518603-k0ptna6taq593eiulqnd4vfsk1djh0vl.apps.googleusercontent.com",
		clientSec: "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",
		endpoint:  google.Endpoint,
		scopes:    []string{drive.DriveScope},
	},
}

func getProviderConfig(provider *OauthProvider) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     provider.clientID,
		ClientSecret: provider.clientSec,
		Endpoint:     provider.endpoint,
		RedirectURL:  redirectURL,
		Scopes:       provider.scopes,
	}
}

func getUserHome() string {
	res, err := os.UserHomeDir()
	if err != nil {
		ExitOnError(err.Error())
	}
	return res
}

func generateRandomBytes(len int) (string, error) {
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil || len < 0 {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// See: https://developers.google.com/identity/protocols/oauth2/native-app#step1-code-verifier
func pkce() (string, string, error) {
	codeVerifier, err := generateRandomBytes(64)
	if err != nil {
		return "", "", err
	}
	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(hash[:])

	return codeVerifier, codeChallenge, nil
}

func getToken() (*oauth2.Token, error) {
	file, err := os.Open(TokenFilePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	token := &oauth2.Token{}
	if err := json.NewDecoder(file).Decode(token); err != nil {
		return nil, err
	}

	return token, nil
}

func setToken(token *oauth2.Token, path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(token)
}

func authorize(provider *OauthProvider) (*http.Client, error) {
	oauthConfig := getProviderConfig(provider)

	code_verifier, code_challenge, err := pkce()

	if err != nil {
		return nil, err
	}

	state, err := generateRandomBytes(16)
	if err != nil {
		return nil, err
	}

	authURL := oauthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", code_challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))

	if err := OpenURL(authURL); err != nil {
		return nil, err
	}

	codeChn := make(chan string)
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code := query.Get("code")
		rState := query.Get("state")
		defer close(codeChn)

		if rState != state {
			http.Error(w, "Invalid 'state' parameter!", http.StatusBadRequest)
			codeChn <- ""
			return
		}

		fmt.Fprintln(w, "Authorization successful! You can close this window now.")
		codeChn <- code
	})

	go func() {
		http.ListenAndServe(":"+port, nil)
	}()
	code := <-codeChn

	// exchange with pkce
	token, err := oauthConfig.Exchange(context.Background(), code,
		oauth2.SetAuthURLParam("code_verifier", code_verifier),
	)
	if err != nil {
		return nil, err
	}

	if err := setToken(token, TokenFilePath); err != nil {
		return nil, err
	}

	client := oauthConfig.Client(context.Background(), token)
	return client, nil
}

func LogIn(providerName string) (*http.Client, error) {
	var provider *OauthProvider
	for i := range Providers {
		if Providers[i].name == providerName {
			provider = &Providers[i]
		}
	}

	if provider == nil {
		return nil, fmt.Errorf("provider %s is not supported", providerName)
	}

	if token, err := getToken(); err == nil {
		return getProviderConfig(provider).Client(context.Background(), token), nil
	}

	Warn("Initialize authorization flow")
	client, err := authorize(provider)
	if err != nil {
		return nil, err
	}

	return client, err
}
