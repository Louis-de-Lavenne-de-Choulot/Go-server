package main

// import (
// 	"encoding/base64"
// 	"encoding/hex"
// 	"io"
// 	"math/rand"
// 	"net/http"
// 	"os"
// 	"sync"
// 	"text/template"
// 	"time" // add this
// )

// const Version = "0.2"

// var (
// 	// db                *sql.DB
// 	// connStr                  = "postgresql://goupb:<password>@<database_ip>/todos?sslmode=enabled"
// 	pwd               string = "."
// 	logfile           string = pwd + "/server.log"
// 	logmutex          sync.Mutex
// 	UserAdditionToken = make(map[int]string)
// 	templates         *template.Template
// 	user              User = User{Id: 0, Username: "admin", Gthb_identifier: "admin", Avatar_url: "https://avatars3.githubusercontent.com/u/14101776?v=3&s=460", Rights: 1, Files_permissions: []int{0}}
// )

// type infos struct {
// 	Title           string
// 	Version         string
// 	Nodes           string
// 	Apps            []App
// 	Rights          int
// 	UserAdditionTkn string
// }

// func main() {
// 	rand.Seed(time.Now().UTC().UnixNano())
// 	// println(Encrypt("Hello"))
// 	NewDB()
// 	// QueryDB("")

// 	templates = template.New("index.html")
// 	//create new template that contains header.html and footer.html
// 	templates.New("header").Parse(readfile(pwd + "/static/header.html"))
// 	templates.New("footer").Parse(readfile(pwd + "/static/footer.html"))
// 	//read the index.html file then serve it
// 	templates.New("index").Parse(readfile(pwd + "/static/index.html"))
// 	//read the lora.html file then serve it
// 	templates.New("lora").Parse(readfile(pwd + "/static/lora.html"))
// 	//read the auth.html file then serve it
// 	templates.New("auth").Parse(readfile(pwd + "/static/auth.html"))
// 	//read the refused.html file then serve it
// 	templates.New("refused").Parse(readfile(pwd + "/static/refused.html"))
// 	//read the schedule.html file then serve it
// 	templates.New("schedule").Parse(readfile(pwd + "/static/schedule.html"))

// 	//init the json_handler package
// 	InitJSON(pwd+"/registeredUsers.json", 1)
// 	InitJSON(pwd+"/end-nodes.json", 2)
// 	InitJSON(pwd+"/apps.json", 3)
// 	for _, nodeFile := range NodesFiles.EndNodes {
// 		InitJSON(pwd+"/nodes/"+nodeFile.FileHash+".json", 0)
// 	}
// 	backup()
// 	AuthSupport()
// 	//handle / and give index.html
// 	http.HandleFunc("/", Index)
// 	//handlefunc that servs auth.html
// 	http.HandleFunc("/auth", Auth)
// 	http.HandleFunc("/refused", Refused)
// 	http.HandleFunc("/lora", Lora)
// 	http.HandleFunc("/schedule", Schedule)
// 	http.HandleFunc("/webhooks/post", PostService)
// 	http.HandleFunc("/api/downlink", DownLink)
// 	http.HandleFunc("/api/getuserinfos", GetUserInfos)
// 	http.HandleFunc("/api/adduser", AddUser)
// 	http.HandleFunc("/api/newschedule/", NewSchedule)
// 	http.HandleFunc("/api/getschedule/", GetSchedule)
// 	http.HandleFunc("/api/update/", ApiUpdate)
// 	http.HandleFunc("/favicon.ico", Favicon)
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	println("Server started on port https")
// }

// // [yyyy-mm-dd:hh-mm-ss-mmss] Recieved $URL with $METHOD
// func logrequest(r *http.Request) {
// 	logmutex.Lock()
// 	defer logmutex.Unlock()
// 	if len(r.URL.Path) > 42 {
// 		r.URL.Path = r.URL.Path[:42]
// 	}
// 	ctime := time.Now().Format(time.RFC850)
// 	logtext := "[" + ctime + "] Received " + r.URL.Path + " with " + r.Method + " method \n"
// 	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
// 	if err != nil {
// 		return
// 	}
// 	defer f.Close()
// 	if _, err = f.WriteString(logtext); err != nil {
// 		return
// 	}
// }

// func loganswer(answer string) {
// 	logmutex.Lock()
// 	defer logmutex.Unlock()
// 	ctime := time.Now().Format(time.RFC850)
// 	logtext := "[" + ctime + "] Answered with " + answer + "\n"
// 	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
// 	if err != nil {
// 		return
// 	}
// 	defer f.Close()
// 	if _, err = f.WriteString(logtext); err != nil {
// 		return
// 	}
// }

// func backup() {
// 	//backup the players.json file
// 	Backup(30, pwd+"/save.json", pwd+"/save2.json", pwd+"/save3.json")
// }

// // func to read file with ioutil and return as string
// func readfile(path string) string {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return ""
// 	}
// 	defer file.Close()
// 	fi, err := file.Stat()
// 	if err != nil {
// 		return ""
// 	}
// 	data := make([]byte, fi.Size())
// 	_, err = file.Read(data)
// 	if err != nil {
// 		return ""
// 	}
// 	return string(data)
// }

// func base64Toutf8(b64 string) string {
// 	//get base64 string and return utf8 string
// 	decoded, err := base64.StdEncoding.DecodeString(b64)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	return string(decoded)
// }

// func base64tohex(b64 string) string {
// 	//get base64 string and return hex string
// 	decoded, err := base64.StdEncoding.DecodeString(b64)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	return hex.EncodeToString(decoded)
// }

// func hexToBase64(h string) (string, error) {
// 	//get hex string and return base64 string
// 	decoded, err := hex.DecodeString(h)
// 	if err != nil {
// 		return "", err
// 	}
// 	return base64.StdEncoding.EncodeToString(decoded), nil
// }

// func Favicon(w http.ResponseWriter, r *http.Request) {}

// func Index(w http.ResponseWriter, r *http.Request) {
// 	logrequest(r)
// 	//Verify session
// 	_, bl := VerifySession(w, r)
// 	if !bl {
// 		return
// 	}
// 	w.Header().Add("Content Type", "text/html")

// 	//Get user to see the rights
// 	// user := GetSessionIDOwner(cookie.Value)
// 	utk := ""
// 	if _, found := UserAdditionToken[user.Id]; found {
// 		utk = UserAdditionToken[user.Id]
// 	}
// 	context := infos{
// 		Title:           "Best-IA",
// 		Version:         Version,
// 		Apps:            AppsI.Apps,
// 		Rights:          1,
// 		UserAdditionTkn: utk,
// 	}
// 	servTemplates(w, []string{"index"}, context)
// 	loganswer("index.html")
// }

// func Lora(w http.ResponseWriter, r *http.Request) {
// 	logrequest(r)
// 	//Verify session
// 	_, bl := VerifySession(w, r)
// 	if !bl {
// 		return
// 	}
// 	w.Header().Add("Content Type", "text/html")

// 	// Get user to see the rights
// 	// user := GetSessionIDOwner(cookie.Value)
// 	utk := ""
// 	if _, found := UserAdditionToken[user.Id]; found {
// 		utk = UserAdditionToken[user.Id]
// 	}
// 	context := infos{
// 		Title:           "LoRawan Server",
// 		Version:         Version,
// 		Nodes:           GetNodes([]int{0}),
// 		Rights:          1,
// 		UserAdditionTkn: utk,
// 	}
// 	servTemplates(w, []string{"lora"}, context)
// 	loganswer("lora.html")
// }

// func Auth(w http.ResponseWriter, r *http.Request) {
// 	//read the auth.html file then serve it
// 	w.Header().Add("Content Type", "text/html")

// 	context := infos{
// 		Title:   "Authentification",
// 		Version: Version,
// 	}
// 	servTemplates(w, []string{"auth"}, context)
// }

// func Refused(w http.ResponseWriter, r *http.Request) {
// 	//read the refused.html file then serve it
// 	w.Header().Add("Content Type", "text/html")

// 	context := infos{
// 		Title:   "Refused",
// 		Version: Version,
// 	}
// 	servTemplates(w, []string{"refused"}, context)
// 	loganswer("refused.html")
// }

// func servTemplates(wr io.Writer, tpltes []string, context interface{}) {
// 	err := templates.ExecuteTemplate(wr, "header", context)
// 	if err != nil {
// 		println(err.Error())
// 		loganswer(err.Error())
// 		return
// 	}
// 	//serve templates
// 	for _, tpl := range tpltes {
// 		err = templates.ExecuteTemplate(wr, tpl, context)
// 		if err != nil {
// 			println(err.Error())
// 			loganswer(err.Error())
// 			return
// 		}
// 	}
// 	err = templates.ExecuteTemplate(wr, "footer", context)
// 	if err != nil {
// 		println(err.Error())
// 		loganswer(err.Error())
// 		return
// 	}
// }
