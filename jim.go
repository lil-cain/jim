package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func check_auth(h, Header) boolean {
	return true
}

func change_light(i int) error {
	return nil
}

func lightHandler(w http.ResponseWriter, r *http.Request) {
	authenticated := check_auth(r.Header)
	if !authenticated {
		http.Error(w, "Not authenticated", 403)
	}
	light_str := r.URL.Path[len("/api/light/")]
	light_int, err := strconv.Atoi(light)
	if err != nil {
		http.Error(w, err.Error, 500)
	}
	err := change_light(light_int)
	if err != nil {
		http.Error(w, err.Error, 500)
	}
	fmt.Print(w, "ok")
}

func main() {
	http.HandleFunc("/", lightHandler)
	http.ListenAndServe(":8080", nil)
}
