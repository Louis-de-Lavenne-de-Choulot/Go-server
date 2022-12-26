package main

import (
	"embed"
	"encoding/base64"
	"encoding/hex"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time" // add this
)

const Version = "0.2"

var (
	// db                *sql.DB
	// connStr                  = "postgresql://goupb:<password>@<database_ip>/todos?sslmode=enabled"
	pwd               string = "/home/groupb/Documents/BestGroup"
	logfile           string = pwd + "/server.log"
	logmutex          sync.Mutex
	UserAdditionToken = make(map[int]string)
	templates         *template.Template
)

type infos struct {
	Title           string
	Version         string
	Nodes           string
	Apps            []App
	Rights          int
	UserAdditionTkn string
}

//go:embed registeredUsers.json
var registeredUsers []byte

//go:embed apps.json
var apps []byte

//go:embed end-nodes.json
var ende []byte

//go:embed nodes
var nodesF embed.FS

//go:embed static
var fs embed.FS

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	//	NewDB()

	listFiles := []string{"header", "footer", "index", "lora", "auth", "refused", "schedule"}

	templates = template.New("html templates")
	//create new template that contains the listed file
	for _, file := range listFiles {
		r, _ := fs.ReadFile("static/" + file + ".html")
		templates.New(file).Parse(string(r))
	}

	InitJSON(pwd+"/registeredUsers.json", 1)
	InitJSON(pwd+"/end-nodes.json", 2)
	InitJSON(pwd+"/apps.json", 3)
	for _, nodeFile := range NodesFiles.EndNodes {
		InitJSON(pwd+"/nodes/"+nodeFile.FileHash+".json", 0)
	}

	backup()
	AuthSupport()
	//handle / and give index.html
	http.HandleFunc("/", Index)
	//handlefunc that servs auth.html
	http.HandleFunc("/auth", Auth)
	http.HandleFunc("/refused", Refused)
	http.HandleFunc("/lora", Lora)
	http.HandleFunc("/schedule", Schedule)
	http.HandleFunc("/webhooks/post", PostService)
	http.HandleFunc("/api/downlink", DownLink)
	http.HandleFunc("/api/getuserinfos", GetUserInfos)
	http.HandleFunc("/api/adduser", AddUser)
	http.HandleFunc("/api/newschedule/", NewSchedule)
	http.HandleFunc("/api/getschedule/", GetSchedule)
	http.HandleFunc("/api/update/", ApiUpdate)
	http.HandleFunc("/favicon.ico", Favicon)
	println("Server started on port https")
	// directly redirect with func
	go http.ListenAndServe(":http", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
	}))
	err := http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/bestiaever.ml/fullchain.pem", "/etc/letsencrypt/live/bestiaever.ml/privkey.pem", nil)
	// err := http.ListenAndServe(":8080", nil)
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
	logtext := "[" + ctime + "] Received " + r.URL.Path + " with " + r.Method + " method \n"
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
	Backup(30, pwd+"/save.json", pwd+"/save2.json", pwd+"/save3.json")
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

func Favicon(w http.ResponseWriter, r *http.Request) {}

func Index(w http.ResponseWriter, r *http.Request) {
	logrequest(r)
	//Verify session
	cookie, bl := VerifySession(w, r)
	if !bl {
		return
	}
	w.Header().Add("Content Type", "text/html")

	//Get user to see the rights
	user := GetSessionIDOwner(cookie.Value)
	utk := ""
	if _, found := UserAdditionToken[user.Id]; found {
		utk = UserAdditionToken[user.Id]
	}
	context := infos{
		Title:           "Best-IA",
		Version:         Version,
		Apps:            AppsI.Apps,
		Rights:          1,
		UserAdditionTkn: utk,
	}
	servTemplates(w, []string{"index"}, context)
	loganswer("index.html")
}

func Lora(w http.ResponseWriter, r *http.Request) {
	logrequest(r)
	//Verify session
	cookie, bl := VerifySession(w, r)
	if !bl {
		return
	}
	w.Header().Add("Content Type", "text/html")

	// Get user to see the rights
	user := GetSessionIDOwner(cookie.Value)
	utk := ""
	if _, found := UserAdditionToken[user.Id]; found {
		utk = UserAdditionToken[user.Id]
	}
	context := infos{
		Title:           "LoRawan Server",
		Version:         Version,
		Nodes:           GetNodes([]int{0}),
		Rights:          1,
		UserAdditionTkn: utk,
	}
	servTemplates(w, []string{"lora"}, context)
	loganswer("lora.html")
}

func Auth(w http.ResponseWriter, r *http.Request) {
	//read the auth.html file then serve it
	w.Header().Add("Content Type", "text/html")

	context := infos{
		Title:   "Authentification",
		Version: Version,
	}
	servTemplates(w, []string{"auth"}, context)
}

func Refused(w http.ResponseWriter, r *http.Request) {
	//read the refused.html file then serve it
	w.Header().Add("Content Type", "text/html")

	context := infos{
		Title:   "Refused",
		Version: Version,
	}
	servTemplates(w, []string{"refused"}, context)
	loganswer("refused.html")
}

func servTemplates(wr io.Writer, tpltes []string, context interface{}) {
	err := templates.ExecuteTemplate(wr, "header", context)
	if err != nil {
		println(err.Error())
		loganswer(err.Error())
		return
	}
	//serve templates
	for _, tpl := range tpltes {
		err = templates.ExecuteTemplate(wr, tpl, context)
		if err != nil {
			println(err.Error())
			loganswer(err.Error())
			return
		}
	}
	err = templates.ExecuteTemplate(wr, "footer", context)
	if err != nil {
		println(err.Error())
		loganswer(err.Error())
		return
	}
}
