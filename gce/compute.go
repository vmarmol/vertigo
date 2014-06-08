package gce

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"

	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/compute/v1"
)

type gcloudCredentialsCache struct {
	Data []gceConfig
}

type gceConfig struct {
	Credential gceCredential
	Key        gceKey
	ProjectId  string `json:"projectId"`
}

type gceCredential struct {
	ClientId     string `json:"Client_Id"`
	ClientSecret string `json:"Client_Secret"`
	RefreshToken string `json:"Refresh_Token"`
}

type gceKey struct {
	Scope string
}

// Gets the OAuth2 token from the current user's gcloud crendentials.
func getOauthToken() (*gceCredential, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	confPath := path.Join(usr.HomeDir, ".config/gcloud/credentials")
	f, err := os.Open(confPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load gcloud credentials: %q", confPath)
	}
	defer f.Close()
	cache := &gcloudCredentialsCache{}
	if err := json.NewDecoder(f).Decode(cache); err != nil {
		return nil, err
	}
	if len(cache.Data) == 0 {
		return nil, fmt.Errorf("no gcloud credentials cached in: %q", confPath)
	}
	return &cache.Data[0].Credential, nil
}

func NewCompute() (*compute.Service, error) {
	// Get Oauth2 token.
	creds, err := getOauthToken()
	if err != nil {
		log.Fatal(err)
	}

	// OAuth2 configuration.
	oAuth2Conf := &oauth.Config{
		ClientId:     creds.ClientId,
		ClientSecret: creds.ClientSecret,
		Scope:        compute.ComputeScope,
		RedirectURL:  "oob",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		AccessType:   "offline",
	}
	transport := &oauth.Transport{
		Config: oAuth2Conf,
		// Make the actual request using the cached token to authenticate.
		Token:     &oauth.Token{RefreshToken: creds.RefreshToken},
		Transport: http.DefaultTransport,
	}

	err = transport.Refresh()
	if err != nil {
		return nil, err
	}
	svc, err := compute.New(transport.Client())
	if err != nil {
		return nil, err
	}
	return svc, nil
}
