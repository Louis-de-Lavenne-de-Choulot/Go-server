package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func PostService(w http.ResponseWriter, r *http.Request) {
	// if len(r.URL.Path) > 50 {
	// 	r.URL.Path = r.URL.Path[:50]
	// }

	logrequest(r)
	// check if the option ?format is present and get the value
	// format := r.URL.Query().Get("format")
	println("called")

	// //print all request
	// for key, value := range r.URL.Query() {
	// 	println(key, value)
	// }
	// for key, value := range r.Header {
	// 	println(key, value)
	// }

	//get or post
	println(r.Method)
	switch r.Method {
	case http.MethodPost:
		//read the body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			loganswer("Error reading body")
			println("Error reading body")
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		mapdata := make(map[string]interface{})
		//parse the body
		err = json.Unmarshal(body, &mapdata)
		if err != nil {
			loganswer("Error parsing body : " + string(body))
			println("Error parsing body : " + string(body))
			http.Error(w, "can't parse body", http.StatusBadRequest)
			return
		}

		//get value of the key "frm_payload" in key "uplink_message"
		if mapdata["uplink_message"] != nil {
			uplink := mapdata["uplink_message"].(map[string]interface{})
			if uplink["frm_payload"] != nil {
				payload := uplink["frm_payload"].(string)
				println(payload)
				//payload conv from base64 to hex
				payload = base64ToHex(payload)
				println(payload)
			}
		}

		loganswer("request validated")
		w.WriteHeader(http.StatusAccepted)
	case http.MethodGet:

		w.WriteHeader(http.StatusMethodNotAllowed)
		println("Method not allowed")
		loganswer("Method not allowed")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		println("Method not allowed")
		loganswer("Method not allowed")
	}
}
