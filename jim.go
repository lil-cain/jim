package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func check_auth(h, Header) {
	return true
}

func lighthandler(w http.ResponseWriter, r *http.Request) {
	authenticated := check_auth(r.Header)
	if !authenticated {
		http.Error(w, "Not authenticated", 403)
	}
	light_str = r.URL.Path[len("/api/light/")]
	light_int, err = strconv.Atoi(light)
	if err != nil {
		http.Error(w, err.Error, 500)
	}
	err = change_light(light_int)
	if err != nil {
		http.Error(w, err.Error, 500)
	}
	fmt.Print(w, "ok")
}
