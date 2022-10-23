package main

import (
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"
)

const Version = "0.2"

var logfile string = "./server.log"
var logmutex sync.Mutex

type infos struct {
	Version string
	Nodes   string
}

func main() {

	templates := template.New("index.html")

	//handle / and give index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logrequest(r)
		//check if cookie "sessionID" exists else redirect to auth
		cookie, err := r.Cookie("sessionID")
		if err != nil {
			http.Redirect(w, r, "/auth", http.StatusFound)
			return
		}
		//check if cookie "sessionID" is equal to an existing session
		if _, found := sessions[cookie.Value]; !found {
			http.Redirect(w, r, "/auth", http.StatusFound)
			return
		}
		w.Header().Add("Content Type", "text/html")
		//read the index.html file then serve it
		templates.New("index").Parse(readfile("./static/index.html"))
		context := infos{
			Version: Version,
			Nodes:   GetNodes(),
		}
		err = templates.ExecuteTemplate(w, "index", context)
		if err != nil {
			println(err.Error())
			loganswer(err.Error())
			return
		}
		loganswer("index.html")
	})

	//handlefunc that servs auth.html
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		//read the auth.html file then serve it
		w.Header().Add("Content Type", "text/html")
		//read the auth.html file then serve it
		templates.New("auth").Parse(readfile("./static/auth.html"))
		err := templates.ExecuteTemplate(w, "auth", nil)
		if err != nil {
			println(err.Error())
			loganswer(err.Error())
			return
		}
		loganswer("auth.html")
	})

	//init the json_handler package
	InitJSON("end-nodes.json", true)
	InitJSON("registeredUsers.json", false)
	backup()
	AuthSupport()
	http.HandleFunc("/webhooks/post", PostService)
	http.HandleFunc("/api/downlink", DownLink)
	http.HandleFunc("/api/getuserinfos", GetUserInfos)
	//handle /api/update and check if parameter "date" is set
	ApiUpdate()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
	println("Server started on port 8080")
}

// [yyyy-mm-dd:hh-mm-ss-mmss] Recieved $URL with $METHOD
func logrequest(r *http.Request) {
	logmutex.Lock()
	defer logmutex.Unlock()
	if len(r.URL.Path) > 42 {
		r.URL.Path = r.URL.Path[:42]
	}
	ctime := time.Now().Format(time.RFC850)
	logtext := "[" + ctime + "] Recieved " + r.URL.Path + " with " + r.Method + " method \n"
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	if _, err = f.WriteString(logtext); err != nil {
		return
	}
}

func loganswer(answer string) {
	logmutex.Lock()
	defer logmutex.Unlock()
	ctime := time.Now().Format(time.RFC850)
	logtext := "[" + ctime + "] Answered with " + answer + "\n"
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	if _, err = f.WriteString(logtext); err != nil {
		return
	}
}

func backup() {
	//backup the players.json file
	Backup(30, "save.json", "save2.json", "save3.json")
}

// func to read file with ioutil and return as string
func readfile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return ""
	}
	data := make([]byte, fi.Size())
	_, err = file.Read(data)
	if err != nil {
		return ""
	}
	return string(data)
}

func base64ToHex(b64 string) string {
	//get base64 string and return hex string
	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err.Error()
	}
	return hex.EncodeToString(decoded)
}

func hexToBase64(h string) (string, error) {
	//get hex string and return base64 string
	decoded, err := hex.DecodeString(h)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(decoded), nil
}
