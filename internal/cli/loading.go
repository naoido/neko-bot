package cli

import (
	"fmt"
	"time"
)

var loadings = []string{"⢿", "⣻", "⣽", "⣾", "⣷", "⣯", "⣟", "⡿"}

func Loading(done chan struct{}) {
	go loading(done)
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
