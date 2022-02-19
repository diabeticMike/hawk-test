package main

import (
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"

	hawk "github.com/coreos/hawk-go"
)

func main() {
	hawkClient := getHawkClient("user1", "key1")
	resp, err := hawkClient.Get("http://localhost:8080/resource")
	if err != nil {
		fmt.Println(err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(string(bodyBytes))
}

var DefaultHawkHasher = sha256.New

type HawkRoundTripper struct {
	User          string
	Token         string
	SkipSSLVerify bool
}

func (t *HawkRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	creds := &hawk.Credentials{
		ID:   t.User,
		Key:  t.Token,
		Hash: DefaultHawkHasher,
	}

	auth := hawk.NewRequestAuth(req, creds, 0)
	req.Header.Set("Authorization", auth.RequestHeader())

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: t.SkipSSLVerify},
	}
	return transport.RoundTrip(req)
}

func getHawkClient(user string, key string) *http.Client {
	return &http.Client{
		Transport: &HawkRoundTripper{
			User:          user,
			Token:         key,
			SkipSSLVerify: true,
		},
	}
}
