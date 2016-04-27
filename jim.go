package jim

import (
	"fmt"
	"net/http"
	"strings"
)

func check_auth(h http.Header) bool {
	return true
}

func router(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request to ", r.URL.Path)
	authenticated := check_auth(r.Header)
	if authenticated {
		fmt.Println("Request Authenticated")
	}

	switch {
	case strings.HasPrefix(r.URL.Path, "/api/light/") && r.Method == "POST":
		setLightState(w, r, authenticated)
	case strings.HasPrefix(r.URL.Path, "/api/state/") && r.Method == "GET":
		getLightState(w, r, authenticated)
	default:
		http.Error(w, "Not Found", 404)
	}

}

func setLightState(w http.ResponseWriter, r *http.Request, authenticated bool) {
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

func getLightState(w http.ResponseWriter, r *http.Request, authenticated bool) {
	if !authenticated {
		http.Error(w, "Not authenticated", 403)
	}
	fmt.Fprintln(w, "1")
}

func main() {
	http.HandleFunc("/", router)
	http.ListenAndServe(":8080", nil)
}
