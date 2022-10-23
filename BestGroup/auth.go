package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	clientID     = "17d14971a346aceec8b2"
	clientSecret = "71a4215aa63f270ded272a0709495b1539c97f1c"
)

func AuthSupport() {
	httpClient := &http.Client{}
	// Create a new redirect route route
	http.HandleFunc("/oauth/redirect/", func(w http.ResponseWriter, r *http.Request) {
		// First, we need to get the value of the `code` query param
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		code := r.FormValue("code")

		// Next, lets for the HTTP request to call the github oauth enpoint
		// to get our access token
		reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
		req, err := http.NewRequest(http.MethodPost, reqURL, nil)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		// We set this header since we want the response
		// as JSON
		req.Header.Set("accept", "application/json")

		// Send out the HTTP request
		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		defer res.Body.Close()

		// Parse the request body into the `OAuthAccessResponse` struct
		var t OAuthAccessResponse
		if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		//set t.AccessToken as cookie
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "token", Value: t.AccessToken, Expires: expiration, Path: "/"}
		//if cookie token exists, delete it
		if _, err := r.Cookie("token"); err == nil {
			cookie.MaxAge = -1
		}

		//set cookie token
		http.SetCookie(w, &cookie)
		println("cookie set")
		//redirect to "/"
		http.Redirect(w, r, "/", http.StatusFound)
	})
}

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}
