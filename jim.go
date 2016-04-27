package main

import (
	"fmt"
	"net/http"
	"strings"
)

func check_auth(h http.Header) bool {
	return true
}

func change_light(i int) error {
	return nil
}

func router(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request to ", r.URL.Path)
	authenticated := check_auth(r.Header)
	if authenticated {
		fmt.Println("Request Authenticated")
	}

	if strings.HasPrefix(r.URL.Path, "/api/light/") {
		lightHandler(w, r, authenticated)
	} else {
		http.Error(w, "Not Found", 404)
	}

}

func lightHandler(w http.ResponseWriter, r *http.Request, authenticated bool) {
	if !authenticated {
		http.Error(w, "Not authenticated", 403)
	}
	light := r.URL.Path[len("/api/light/")] - '0'
	err := change_light(int(light))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	fmt.Fprintln(w, "ok")
}

func main() {
	http.HandleFunc("/", router)
	http.ListenAndServe(":8080", nil)
}
