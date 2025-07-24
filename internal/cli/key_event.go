package cli

import (
	"context"
	"fmt"
	"golang.org/x/term"
	"os"
)

type KeyEvent struct {
	keyChan chan string
}

func NewKeyEvent() *KeyEvent {
	return &KeyEvent{
		keyChan: make(chan string),
	}
}

func (k *KeyEvent) Listening(ctx context.Context) (<-chan string, *term.State, error) {
	// ターミナルTTYが使えない環境の場合はエラー
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil, nil, fmt.Errorf("KeyListener skipped: no TTY available")
	}

	// Rawモード以前のターミナル状態を保存しておく
	oldTermState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("failed to set raw mode in terminal: %s\n", err)
	}

	go func() {
		input := make([]byte, 1)

		for {
			n, err := os.Stdin.Read(input)
			if err != nil {
				fmt.Println("")
				return
			}

			select {
			case <-ctx.Done():
				break
			default:
				if n > 0 {
					k.keyChan <- string(input[0])
				}
			}
		}
	}()

	return k.keyChan, oldTermState, nil
}
