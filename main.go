package main

import (
	"fmt"
	"net/http"
	"backend-assignment/handler"
)

func main() {
	http.HandleFunc("/issue", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateIssue(w, r)
		default:
			http.Error(w, "허용되지 않은 메서드입니다.", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/issues", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.ListIssues(w, r)
		} else {
			http.Error(w, "허용되지 않은 메서드입니다.", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/issue/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetIssue(w, r)
		case http.MethodPatch:
			handler.UpdateIssue(w, r)
		default:
			http.Error(w, "허용되지 않은 메서드입니다.", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
