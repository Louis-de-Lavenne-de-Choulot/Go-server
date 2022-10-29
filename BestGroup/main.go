package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time" // add this

	// add this
	_ "github.com/lib/pq" // add this
)

const Version = "0.2"

var (
	db                *sql.DB
	connStr                  = "postgresql://<username>:<password>@<database_ip>/todos?sslmode=enabled"
	logfile           string = "/home/groupb/Documents/BestGroup/server.log"
	logmutex          sync.Mutex
	UserAdditionToken = make(map[int]string)
	templates         *template.Template
)

type infos struct {
	Version         string
	Nodes           string
	Rights          int
	UserAdditionTkn string
}

func main() {
	db

	templates = template.New("index.html")

	//init the json_handler package
	InitJSON("/home/groupb/Documents/BestGroup/registeredUsers.json", 1)
	InitJSON("/home/groupb/Documents/BestGroup/end-nodes.json", 2)
	for _, nodeFile := range NodesFiles.EndNodes {
		InitJSON("/home/groupb/Documents/BestGroup/nodes/"+nodeFile.FileHash+".json", 0)
	}
	backup()
	AuthSupport()
	//handle / and give index.html
	http.HandleFunc("/", Index)
	//handlefunc that servs auth.html
	http.HandleFunc("/auth", Auth)
	http.HandleFunc("/refused", Refused)
	http.HandleFunc("/webhooks/post", PostService)
	http.HandleFunc("/api/downlink", DownLink)
	http.HandleFunc("/api/getuserinfos", GetUserInfos)
	http.HandleFunc("/api/adduser", AddUser)
	http.HandleFunc("/api/update/", ApiUpdate)

	err := http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/bestiaever.ml/fullchain.pem", "/etc/letsencrypt/live/bestiaever.ml/privkey.pem", nil)
	if err != nil {
		panic(err)
	}
	println("Server started on port https")
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
	Backup(30, "/home/groupb/Documents/BestGroup/save.json", "/home/groupb/Documents/BestGroup/save2.json", "/home/groupb/Documents/BestGroup/save3.json")
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

func base64Toutf8(b64 string) string {
	//get base64 string and return utf8 string
	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err.Error()
	}
	return string(decoded)
}

func base64tohex(b64 string) string {
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

func Index(w http.ResponseWriter, r *http.Request) {
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
	templates.New("index").Parse(readfile("/home/groupb/Documents/BestGroup/static/index.html"))

	//Get user to see the rights
	user := GetSessionIDOwner(cookie.Value)
	utk := ""
	if _, found := UserAdditionToken[user.Id]; found {
		utk = UserAdditionToken[user.Id]
	}
	context := infos{
		Version:         Version,
		Nodes:           GetNodes(user.Files_permissions),
		Rights:          user.Rights,
		UserAdditionTkn: utk,
	}
	err = templates.ExecuteTemplate(w, "index", context)
	if err != nil {
		println(err.Error())
		loganswer(err.Error())
		return
	}
	loganswer("index.html")
}

func Auth(w http.ResponseWriter, r *http.Request) {
	//read the auth.html file then serve it
	w.Header().Add("Content Type", "text/html")
	//read the auth.html file then serve it
	templates.New("auth").Parse(readfile("/home/groupb/Documents/BestGroup/static/auth.html"))
	err := templates.ExecuteTemplate(w, "auth", nil)
	if err != nil {
		println(err.Error())
		loganswer(err.Error())
		return
	}
	loganswer("auth.html")
}

func Refused(w http.ResponseWriter, r *http.Request) {
	//read the refused.html file then serve it
	w.Header().Add("Content Type", "text/html")
	//read the refused.html file then serve it
	templates.New("refused").Parse(readfile("/home/groupb/Documents/BestGroup/static/refused.html"))
	err := templates.ExecuteTemplate(w, "refused", nil)
	if err != nil {
		println(err.Error())
		loganswer(err.Error())
		return
	}
	loganswer("refused.html")
}
