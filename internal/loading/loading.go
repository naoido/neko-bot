package loading

import (
	"fmt"
	"time"
)

var loadings = []string{"⢿", "⣻", "⣽", "⣾", "⣷", "⣯", "⣟", "⡿"}

var s chan struct{}

func Start() {
	s = make(chan struct{})
	go loading(s)
}

func Stop() {
	close(s)
}

func loading(stop chan struct{}) {
loop:
	for {
		for _, l := range loadings {
			select {
			case <-stop:
				break loop
			default:
				fmt.Print("\r" + l)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
	fmt.Print("\r ")
}
