package main

import "net/http"

type schedule struct {
	Title     string
	Version   string
	Days      []int
	Role      string
	Tasks     []string
	Notifs    []string
	Important []string
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
		Title:     "Schedule",
		Version:   Version,
		Days:      []int{}, //31, 1, 2, 3, 4, 5, 6},
		Role:      "Project Manager",
		Tasks:     []string{}, //"Lorem ipsum dolor sit amet", "Lorem ipsum dolor sit amet", "Lorem ipsum dolor sit amet", "Lorem ipsum dolor sit amet", "Lorem ipsum dolor sit amet", "Lorem ipsum dolor sit amet", "Lorem ipsum dolor sit amet"},
		Notifs:    []string{"Project Begins", "Ivan has finished the technical"},
		Important: []string{"", "Finish Technical", "", "", "", "", "End of Week Meeting"},
	}

	servTemplates(w, []string{"schedule"}, contextSc)
}
