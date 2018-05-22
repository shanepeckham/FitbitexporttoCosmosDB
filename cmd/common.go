package cmd

import (
	"golang.org/x/oauth2"
	"net/http"
	"golang.org/x/oauth2/fitbit"
	"os"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"errors"
	"log"
	"net/url"
)
var cachePath string
var exportDate string

var conf *oauth2.Config
var errExpiredToken = errors.New("expired token")
var err error

type cacherTransport struct {
	Base *oauth2.Transport
	Path string
}

type DataPoint struct {
	Time string `json:"dateTime"`
	Value string `json:"value"`
}

func oauthDance(cachePath string) *http.Client {

	conf := &oauth2.Config{
		ClientID:     os.Getenv("FITBIT_CLIENTID"),
		ClientSecret: os.Getenv("FITBIT_CLIENT_SECRET"), //*clientSecret,
		Scopes:       []string{"sleep", "activity", "heartrate", "profile"},
		Endpoint:     fitbit.Endpoint, // note the call back
	}

	token, err := tokenFromFile(cachePath)
	if err != nil && os.IsNotExist(err) {
		token, err = authorize(conf)
	}
	if err != nil {
		log.Fatal(err)
	}

	c := client(conf, token)
	return c
}

func refreshToken(c *http.Client) {
	// Send a request just to see if it errors. This quickly detects expired
	// tokens and allows us to re-authorize.
	if _, err := c.Get("https://api.fitbit.com/1/user/-/activities.json"); err != nil {
		if urlErr, ok := err.(*url.Error); !ok || urlErr.Err != errExpiredToken {
			log.Fatal(err)
		}
		if token, tokenErr := authorize(conf); tokenErr == nil {
			c = client(conf, token)
		} else {
			log.Fatalf("Request resulted in %v, trying to re-authorize resulted in %v", err, tokenErr)
		}
	}
}


func authorize(conf *oauth2.Config) (*oauth2.Token, error) {
	tokens := make(chan *oauth2.Token)
	errors := make(chan error)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		/* 		if code == "" {
			http.Error(w, "Missing 'code' parameter", http.StatusBadRequest)
			return
		} */
		tok, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			errors <- fmt.Errorf("could not exchange auth code for a token: %v", err)
			return
		}
		tokens <- tok
	})
	go func() {
		// Unfortunately, we need to hard-code this port — when registering
		// with fitbit, full RedirectURLs need to be whitelisted (incl. port).
		errors <- http.ListenAndServe(":7319", nil)
	}()

	//authUrl := "https://www.fitbit.com/oauth2/authorize?response_type=token&client_id=22CSMZ&redirect_uri=http%3A%2F%2Flocalhost%3A7319%2F&response_type=code&scope=activity%20heartrate%20sleep&expires_in=604800"

	authUrl := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Println("Please visit the following URL to authorize:")
	fmt.Println(authUrl)
	select {
	case err := <-errors:
		return nil, err
	case token := <-tokens:
		return token, nil
	}
}

// Like oauth2.Config.Client(), but using cacherTransport to persist tokens.
func client(config *oauth2.Config, token *oauth2.Token) *http.Client {
	return &http.Client{
		Transport: &cacherTransport{
			Path: cachePath,
			Base: &oauth2.Transport{
				Source: config.TokenSource(oauth2.NoContext, token),
			},
		},
	}
}

func tokenFromFile(path string) (*oauth2.Token, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var t oauth2.Token
	if err := json.Unmarshal(content, &t); err != nil {
		return nil, err
	}
	return &t, nil
}


func (c *cacherTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if c.Path == "" {
		c.Path = tokenMount
	}
	cachedToken, err := tokenFromFile(c.Path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if _, err := c.Base.Source.Token(); err != nil {
		return nil, errExpiredToken
	}
	resp, err = c.Base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	newTok, err := c.Base.Source.Token()
	if err != nil {
		// While we’re unable to obtain a new token, the request was still
		// successful, so let’s gracefully handle this error by not caching a
		// new token. In either case, the user will need to re-authenticate.
		return resp, nil
	}
	if cachedToken == nil ||
		cachedToken.AccessToken != newTok.AccessToken ||
		cachedToken.RefreshToken != newTok.RefreshToken {
		bytes, err := json.Marshal(&newTok)
		if err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(c.Path, bytes, 0600); err != nil {
			return nil, err
		}
	}
	return resp, nil
}