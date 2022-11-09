package main

import (
	"net/http"
	"time"
)

type schedule struct {
	Title   string
	Version string
	Date    string
}

func Schedule(w http.ResponseWriter, r *http.Request) {
	logrequest(r)
	//Verify session
	_, bl := VerifySession(w, r)
	if !bl {
		return
	}
	//read the auth.html file then serve it
	w.Header().Add("Content Type", "text/html")

	contextSc := schedule{
		Title:   "Schedule",
		Version: Version,
		Date:    time.Now().Format("2006-01-02"),
	}

	servTemplates(w, []string{"schedule"}, contextSc)
}
