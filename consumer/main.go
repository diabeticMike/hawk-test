// sample server
package main

import (
	"fmt"
	"github.com/hiyosi/hawk"
	"net/http"
	"time"
)

type credentialStore struct{}

func (c *credentialStore) GetCredential(id string) (*hawk.Credential, error) {
	return &hawk.Credential{
		ID:  id,
		Key: "key1",
		Alg: hawk.SHA256,
	}, nil
}

var testCredStore = &credentialStore{}

func hawkHandler(w http.ResponseWriter, r *http.Request) {
	s := hawk.NewServer(testCredStore)

	// authenticate client request
	cred, err := s.Authenticate(r)
	if err != nil {
		w.Header().Set("WWW-Authenticate", "Hawk")
		w.WriteHeader(401)
		fmt.Println(err)
		return
	}

	opt := &hawk.Option{
		TimeStamp: time.Now().Unix(),
		Ext:       "response-specific",
	}

	// build server response header
	h, _ := s.Header(r, cred, opt)

	w.Header().Set("Server-Authorization", h)
	w.WriteHeader(200)
	w.Write([]byte("Hello, " + cred.ID))
}

func main() {
	http.HandleFunc("/resource", hawkHandler)
	http.ListenAndServe(":8080", nil)
}
