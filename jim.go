package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"net/http"
	"strconv"
	"strings"
	"time"
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
		return
	}
	light, err := strconv.Atoi(r.URL.Path[len("/api/light/"):])
	if err != nil {
		http.Error(w, "Bad URL passed", 400)
		return
	}

	if light < 0 || light > 6 {
		http.Error(w, "Bad state passed", 400)
		return
	}

	err = change_light(int64(light))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	fmt.Fprintln(w, "ok")
}

func get_light_state(w http.ResponseWriter, r *http.Request, authenticated bool) {
	if !authenticated {
		http.Error(w, "Not authenticated", 403)
		return
	}
	light_on := false

	// map states to pins
	pins := make(map[int64]int64)
	pins[5] = 27
	pins[4] = 22
	pins[3] = 9
	pins[2] = 11
	pins[1] = 13
	for light, pin_number := range pins {
		pin := rpio.Pin(pin_number)
		if pin.Read() == rpio.Low {
			fmt.Fprintln(w, light)
		}
	}
	if !light_on {
		fmt.Fprintln(w, 0)
	}
}

func change_light(state int64) error {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		return err
	}
	defer rpio.Close()

	// map states to pins
	pins := make(map[int64]int64)
	pins[5] = 27
	pins[4] = 22
	pins[3] = 9
	pins[2] = 11
	pins[1] = 13

	// initialise lights
	for _, v := range pins {
		pin := rpio.Pin(v)
		pin.Output()
		pin.High()
		time.Sleep(100 * time.Millisecond)
	}

	// non zero state turns a light on
	if state != 0 {
		pin := rpio.Pin(pins[state])
		pin.Output()
		pin.Low()
	}

	return nil
}

func main() {
	http.HandleFunc("/", router)
	http.ListenAndServe(":80", nil)
}
