package main

import (
	"bytes"
	"net/http"
	"os"
)

func DownLink(w http.ResponseWriter, r *http.Request) {
	key := os.Getenv("API_KEY")
	// post to 'https://thethings.example.com/api/v3/as/applications/algosup-group8-v1/webhooks/wh1/devices/eui-2cf7f1202490140b/down/push' what r contains
	client := &http.Client{}
	//encode r.FormValue("frm_payload") if error return
	frm_p, err := hexToBase64(r.FormValue("frm_payload"))
	if err != nil {
		http.Error(w, "can't encode frm_payload", http.StatusBadRequest)
		return
	}
	json := `{"downlinks":[{"f_port":` + r.FormValue("f_port") + `,"frm_payload":"` + frm_p + `","priority":"NORMAL"}]}`
	println(json)
	req, err := http.NewRequest("POST", "https://eu1.cloud.thethings.network/api/v3/as/applications/algosup-group8-v1/webhooks/app-solu/devices/eui-2cf7f1202490140b/down/push", bytes.NewBuffer([]byte(json)))
	if err != nil {
		println("Error creating request")
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
