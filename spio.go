package jim

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"time"
)

func change_light(state int64) {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
	}
	defer rpio.Close()

	pins := make(map[int64]int64)
	pins[5] = 11
	pins[4] = 7
	pins[3] = 10
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

// func main() {
// 	change_light(5)
// }
