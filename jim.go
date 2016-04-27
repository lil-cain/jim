package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"github.com/stianeikeland/go-rpio"
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
		set_light_state(w, r, authenticated)
	case strings.HasPrefix(r.URL.Path, "/api/state/") && r.Method == "GET":
		get_light_state(w, r, authenticated)
	default:
		http.Error(w, "Not Found", 404)
	}

}

func set_light_state(w http.ResponseWriter, r *http.Request, authenticated bool) {
	if !authenticated {
		http.Error(w, "Not authenticated", 403)
	}
	light := r.URL.Path[len("/api/light/")] - '0'
	err := change_light(int64(light))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	fmt.Fprintln(w, "ok")
}

func get_light_state(w http.ResponseWriter, r *http.Request, authenticated bool) {
	if !authenticated {
		http.Error(w, "Not authenticated", 403)
	}
	fmt.Fprintln(w, "1")
}

func change_light(state int64) error {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
    return err
	}
	defer rpio.Close()

	if state == 0 {
		// stub for switching everything off
		fmt.Println("Turning off everything")
	} else {
		pins := make(map[int64]int64)
		pins[5] = 11
		pins[4] = 7
		pins[3] = 4
		pins[2] = 19
		pins[1] = 22

		for _, v := range pins {
			pin := rpio.Pin(v)
			pin.Output()
			pin.High()
			time.Sleep(100 * time.Millisecond)
		}

		pin := rpio.Pin(pins[state])
		pin.Output()
		pin.Low()
	}

	return nil
}

func main() {
	http.HandleFunc("/", router)
	http.ListenAndServe(":8080", nil)
}
