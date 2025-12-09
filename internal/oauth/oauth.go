// https://auth0.com/docs/get-started/authentication-and-authorization-flow/authorization-code-flow-with-pkce
package oauth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/udaycmd/rdv/utils"
	// "github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

var (
	authPort        string = "5330"
	authRedirectURL string = fmt.Sprintf("http://localhost:%s/callback", authPort)
)

func genRandomBytes(len int) (string, error) {
	if len < 0 {
		return "", fmt.Errorf("byte length must be greater than 0")
	}

	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// https://developers.google.com/identity/protocols/oauth2/native-app#step1-code-verifier
// https://developers.dropbox.com/oauth-guide#implementing-oauth
func pkce() (string, string, error) {
	codeVerifier, err := genRandomBytes(64)
	if err != nil {
		return "", "", err
	}
	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(hash[:])

	return codeVerifier, codeChallenge, nil
}

func Authorize(op OauthProvider) error {
	bc := op.GetCfg()
	oConf := &oauth2.Config{
		ClientID:     bc.ClientId,
		ClientSecret: bc.Secret,
		Endpoint:     bc.Ep,
		RedirectURL:  authRedirectURL,
		Scopes:       bc.Scopes,
	}

	code_verifier, code_challenge, err := pkce()
	if err != nil {
		return err
	}

	state, err := genRandomBytes(16)
	if err != nil {
		return err
	}

	authURL := oConf.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("token_access_type", "offline"),
		oauth2.SetAuthURLParam("code_challenge", code_challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))

	if err := utils.OpenURL(authURL); err != nil {
		return err
	}

	codeChn := make(chan string)
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code := query.Get("code")
		rs := query.Get("state")
		defer close(codeChn)

		if rs != state {
			http.Error(w, "Invalid 'state' parameter!", http.StatusBadRequest)
			codeChn <- ""
			return
		}

		fmt.Fprintln(w, "Authorization successfull! You can close this window now.")
		codeChn <- code
	})

	go http.ListenAndServe(":"+authPort, nil)
	code := <-codeChn

	// exchange with pkce
	token, err := oConf.Exchange(context.Background(), code,
		oauth2.SetAuthURLParam("code_verifier", code_verifier),
	)
	if err != nil {
		return err
	}

	return SetToken(op, token)
}
