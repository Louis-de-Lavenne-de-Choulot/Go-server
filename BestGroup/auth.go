package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	clientID     = "17d14971a346aceec8b2"
	clientSecret = "71a4215aa63f270ded272a0709495b1539c97f1c"
)

var (
	// var that contains session IDs and the id of the user the session redirects to
	sessions = make(map[string]int)
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

		exist, user := GetUser(t.AccessToken)
		if exist {
			//create randome 20 number string
			sessionID := RandStringBytesRmndr(50)
			//add t.AccessToken and it's UserID to sessions
			sessions[sessionID] = user.Id

			//set t.AccessToken as cookie
			expiration := time.Now().Add(30 * time.Minute)
			cookie := http.Cookie{Name: "sessionID", Value: sessionID, Expires: expiration, Path: "/"}

			//set cookie token
			http.SetCookie(w, &cookie)

			utk := RandStringBytesRmndr(32)
			//check if useradditiontoken[user.id] does not exist
			if _, ok := UserAdditionToken[user.Id]; !ok {
				if user.Rights == 777 {
					//add utk to UserAdditionToken
					UserAdditionToken[user.Id] = utk
				}
			}

			time.AfterFunc(30*time.Minute, func() { delete(sessions, t.AccessToken); delete(UserAdditionToken, user.Id) })
			//redirect to "/"
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			//redirect to "/refused"
			http.Redirect(w, r, "/refused", http.StatusFound)
		}
	})
}

func GetUser(token string) (bool, User) {
	httpClient := &http.Client{}

	//do a request to https://api.github.com/user with headers Authorization: token t.AccessToken and Accept: application/json
	reqURL := "https://api.github.com/user"
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
		return false, User{}
	}

	// We set this header since we want the response as JSON
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "token "+token)

	// Send out the HTTP request
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
		return false, User{}
	}
	defer res.Body.Close()
	//find fied "login" in response and get value
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not read response body: %v", err)
		return false, User{}
	}
	var t map[string]interface{}
	if err := json.Unmarshal(body, &t); err != nil {
		fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
		return false, User{}
	}
	count := 0
	for _, user := range UsersI.Users {
		println(user.Id, t["login"], user.Gthb_identifier)
		if t["login"] != nil && user.Gthb_identifier == t["login"].(string) {
			println("accepted")
			//if username is "" then take t["name"].(string)
			if user.Username == "" && t["name"] != nil {
				UsersI.Users[count].Username = t["name"].(string)
			} else if user.Username == "" && t["login"] != nil {
				UsersI.Users[count].Username = t["login"].(string)
			}
			if user.Avatar_url == "" && t["avatar_url"] != nil {
				UsersI.Users[count].Avatar_url = t["avatar_url"].(string)
			}
			UsersI.Users[count].token = token
			//update registeredUsers.json
			SaveJSON(pwd+"/registeredUsers.json", AutoGenerated{}, UsersI, NodesI{})
			return true, user
		}
		count++
	}
	return false, User{}
}

func GetUserInfos(w http.ResponseWriter, r *http.Request) {

	//Verify session
	cookie, bl := VerifySession(w, r)
	if !bl {
		return
	}

	//get user id from sessions
	user := GetSessionIDOwner(cookie.Value)
	if user.Gthb_identifier != "" {
		user.token = ""
		user.Gthb_identifier = ""
		user.Id = 0
		user.Rights = 0
		user.Files_permissions = nil
		//send user infos
		json.NewEncoder(w).Encode(user)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func GetSessionIDOwner(sessionID string) User {
	userID := sessions[sessionID]
	for _, user := range UsersI.Users {
		if user.Id == userID {
			return user
		}
	}
	return User{}
}
func VerifySession(w http.ResponseWriter, r *http.Request) (*http.Cookie, bool) {
	var cookie *http.Cookie
	//check if cookie "sessionID" exists else redirect to auth
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/auth", http.StatusFound)
		return cookie, false
	}
	//check if cookie "sessionID" is equal to an existing session
	if _, found := sessions[cookie.Value]; !found {
		http.Redirect(w, r, "/auth", http.StatusFound)
		return cookie, false
	}
	return cookie, true
}

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
