package listening

import (
	"fmt"
	"golang.org/x/term"
	"neko-bot/internal/errors"
	"neko-bot/neko"
	"os"
)

func KeyListener() {
	buffer := make([]byte, 1)
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	errors.Catch(err, "failed to set raw mode in terminal")
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	for {
		n, err := os.Stdin.Read(buffer)
		errors.Catch(err, fmt.Sprintf("invalid input error [%s]", err))

		if n > 0 {
			inputChar := buffer[0]
			switch inputChar {
			case 'r', 'R':
				fmt.Println("\rreloading...\r")
				neko.UpdateBot(true)
				fmt.Println("\rfinished!\r")
			case 'q':
				return
			}
			fmt.Print("\r")
		}
	}
}
