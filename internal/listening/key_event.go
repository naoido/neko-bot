package listening

import (
	"fmt"
	"golang.org/x/term"
	"neko-bot/discord/bot"
	"neko-bot/discord/handler"
	"neko-bot/internal/errors"
	"os"
	"syscall"
)

func KeyListener() {
	buffer := make([]byte, 1)
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	errors.Catch(err, "failed to set raw mode in terminal")
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	exitCode := byte(3)

	for {
		n, err := os.Stdin.Read(buffer)
		errors.Catch(err, fmt.Sprintf("invalid input error [%s]", err))

		if n > 0 {
			inputChar := buffer[0]
			switch inputChar {
			case 'r', 'R':
				Printlr("reloading...")
				err := bot.Update()
				errors.CatchAndPanic(err, "failed to reload")
				Printlr("finished!")
			case 'c':
				commands := handler.GetRegisteredCommands()
				width, _, err := term.GetSize(syscall.Stdin)
				if err != nil {
					errors.Catch(err, "failed to get terminal size")
					return
				}
				printFilledLine(width, '*', '-')
				Printlr("")
				for _, cmd := range commands {
					Printlr("    /" + fmt.Sprintf("%-10s", cmd.Name) +
						"                            :" + cmd.Description + "    ")
					for _, option := range cmd.Options {
						Printlr(
							"               " +
								"    " + fmt.Sprintf("%-20s", option.Name) +
								"    :" + fmt.Sprintf("%-30s", option.Description) +
								"    " + fmt.Sprintf("%-15s", option.Type.String()) +
								"    " + fmt.Sprintf("%t", option.Required))
					}
					Printlr("")
				}
				printFilledLine(width, '*', '-')
			case 'q', exitCode:
				return
			}
		}
	}
}

func printFilledLine(width int, edgeChar byte, fillChar byte) {
	l := make([]byte, width)
	l[0] = edgeChar
	for i := 1; i < (width - 1); i++ {
		l[i] = fillChar
	}
	l[width-1] = edgeChar
	Printlr(string(l))
}

func Printlr(str string) {
	fmt.Println("\r" + str + "\r")
}
