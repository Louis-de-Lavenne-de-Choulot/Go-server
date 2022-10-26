package main

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
)

func DownLink(w http.ResponseWriter, r *http.Request) {
	logrequest(r)

	//get cookie and check if user is logged in
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		loganswer("not logged in")
		return
	}

	//check user id
	user := GetSessionIDOwner(cookie.Value)
	//check if parameter "current" is set
	if r.FormValue("current") == "" {
		http.Error(w, "current not set", http.StatusBadRequest)
		loganswer("current not set")
		return
	}
	current, _ := strconv.Atoi(r.FormValue("current"))
	println(len(user.Files_permissions), current)
	if len(user.Files_permissions) < current {
		http.Error(w, "Modified Request", http.StatusUnauthorized)
		loganswer("Modified Request")
		return
	}

	// key := os.Getenv("API_KEY")
	keyID := user.Files_permissions[current]
	key := ""
	wbksApp := ""
	for _, v := range NodesFiles.EndNodes {
		if v.ID == keyID {
			key = v.APIKey
			wbksApp = v.WebhookAPI
		}
	}
	println(key)
	println(keyID)

	// post to ttn what r contains
	client := &http.Client{}
	//encode r.FormValue("frm_payload") if error return
	frm_p, err := hexToBase64(r.FormValue("frm_payload"))
	if err != nil {
		http.Error(w, "can't encode frm_payload", http.StatusBadRequest)
		return
	}

	//create json
	json := `{"downlinks":[{"f_port":` + r.FormValue("f_port") + `,"frm_payload":"` + frm_p + `","priority":"` + r.FormValue("priority") + `"}]}`

	var req *http.Request

	if len(Nodes[keyID].Nodes) > 0 {
		req, err = http.NewRequest("POST", "https://eu1.cloud.thethings.network/api/v3/as/applications/"+Nodes[keyID].Nodes[0].EndDeviceIds.ApplicationIds.ApplicationID+"/webhooks/"+wbksApp+"/devices/"+Nodes[keyID].Nodes[0].EndDeviceIds.DeviceID+"/down/push", bytes.NewBuffer([]byte(json)))
		if err != nil {
			println("Error creating request")
		}
	} else if r.FormValue("dev_id") != "" {
		req, err = http.NewRequest("POST", "https://eu1.cloud.thethings.network/api/v3/as/applications/"+r.FormValue("app_id")+"/webhooks/"+wbksApp+"/devices/"+r.FormValue("dev_id")+"/down/push", bytes.NewBuffer([]byte(json)))
		if err != nil {
			println("Error creating request")
		}
	} else {
		http.Error(w, "no device id and/or app id", http.StatusBadRequest)
		loganswer("no deviceid  and/or app id")
		return
	}

	req.Header.Add("Authorization", "Bearer "+key)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "my-integration/my-integration-version")
	resp, err := client.Do(req)
	if err != nil {
		println("Error sending request : ", err.Error())

	} else {
		defer resp.Body.Close()
	}
	//w.write resp.body
	rB := new(bytes.Buffer)
	rB.ReadFrom(resp.Body)
	if rB.String() != "" {
		w.Write(rB.Bytes())
	} else {
		//redirect to / after 1s
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// check for update in dynamic mode, uses date which needs to be changed for ID
func ApiUpdate(w http.ResponseWriter, r *http.Request) {
	//get cookie and check get user id
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		loganswer("not logged in")
		return
	}
	//check user id
	user := GetSessionIDOwner(cookie.Value)
	current, _ := strconv.Atoi(r.FormValue("current"))
	if len(user.Files_permissions) <= current {
		http.Error(w, "Modified Request", http.StatusUnauthorized)
		loganswer("Modified Request")
		return
	}

	id := user.Files_permissions[current]
	logrequest(r)
	//check if parameter "date" is set and if latest Nodes.Nodes[0].ReceivedAt date is newer
	if r.URL.Query().Get("date") != "" {
		if Nodes[id].Nodes != nil {
			//replace + by space in date
			firstNode := Nodes[id].Nodes[len(Nodes[id].Nodes)-1].ReceivedAt.String()
			firstNode = strings.Replace(firstNode, "+", " ", 1)
			lastNode := Nodes[id].Nodes[0].ReceivedAt.String()
			lastNode = strings.Replace(lastNode, "+", " ", 1)
			if r.URL.Query().Get("date") != lastNode || r.URL.Query().Get("date") != firstNode {
				//if set, return GetNodes() in body``
				w.Write([]byte(GetNodes(user.Files_permissions)))
				loganswer("api/update")
			} else {
				http.Error(w, "no new data", http.StatusNoContent)
				loganswer("up to date")
			}
		} else {
			//if not set, return http error
			http.Error(w, "parameter 'date' not set", http.StatusBadRequest)
			loganswer("parameter 'date' not set")
		}
	} else {
		//return http status 200
		w.WriteHeader(http.StatusOK)
	}
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	logrequest(r)
	//get cookie and check if user is logged in as admin
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		loganswer("not logged in")
		return
	}
	//check if user is admin
	user := GetSessionIDOwner(cookie.Value)
	if user.Gthb_identifier != "" {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		loganswer("not logged in")
		return
	}

	if user.Rights != 777 {
		http.Error(w, "not admin", http.StatusUnauthorized)
		loganswer("not admin")
		return
	}

	//check if user is logged in
	if r.FormValue("github_identifier") != "" && r.FormValue("user_addition_token") != "" {
		for _, user := range UsersI.Users {
			if user.Gthb_identifier == r.FormValue("github_identifier") {
				http.Error(w, "user already exists", http.StatusBadRequest)
				loganswer("user already exists")
				return
			}
		}
		tkn := r.FormValue("user_addition_token")
		//check if user_addition_token is in the list of tokens UserAdditionTokens
		if _, found := UserAdditionToken[user.Id]; found {
			if UserAdditionToken[user.Id] == tkn {
				rights := 0
				if r.FormValue("rights") != "" {
					rights, _ = strconv.Atoi(r.FormValue("rights"))
				}
				//if yes, add user to UsersI.Users
				UsersI.Users = append(UsersI.Users, User{Gthb_identifier: r.FormValue("github_identifier"), Rights: rights})
				//save UsersI.Users
				SaveJSON("registeredUsers.json", AutoGenerated{}, UsersI, NodesI{})
				loganswer("user added")
				http.Redirect(w, r, "/", http.StatusOK)
				return
			} else {
				http.Error(w, "wrong token", http.StatusBadRequest)
				loganswer("wrong token")
				return
			}
		} else {
			http.Error(w, "user has no token", http.StatusBadRequest)
			loganswer("user has no token")
			return
		}
	} else {
		//return http error
		http.Error(w, "username or password not set", http.StatusBadRequest)
		loganswer("username or password not set")
	}
}
