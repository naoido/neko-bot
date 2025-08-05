package cli

import (
	"context"
	"fmt"
	"golang.org/x/term"
	"os"
	"os/signal"
	"syscall"
)

type CLI interface {
	Start(ctx context.Context) error
	Done() <-chan struct{}
}

type cli struct {
	done chan struct{}
}

func NewCLI() CLI {
	return &cli{
		done: make(chan struct{}),
	}
}

func (c *cli) Start(ctx context.Context) error {
	internal, cancel := context.WithCancel(ctx)
	defer cancel()

	// OSシグナルを受信するため
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// キーイベントを受け取るため（TTYが利用可能な場合のみ）
	keyEvent := NewKeyEvent()
	input, t, err := keyEvent.Listening(internal)
	
	var hasKeyInput bool
	if err != nil {
		fmt.Printf("Key input disabled: %v\n", err)
		hasKeyInput = false
	} else {
		hasKeyInput = true
		defer func() {
			if t != nil {
				_ = term.Restore(int(os.Stdin.Fd()), t)
			}
		}()
	}

LOOP:
	for {
		if hasKeyInput {
			select {
			case cmd := <-input:
				err := c.handler(cmd, cancel)
				if err != nil {
					fmt.Println("CLI操作でエラーが発生しました :", err)
					cancel()
				}
			case <-quit:
				break LOOP
			case <-internal.Done():
				break LOOP
			}
		} else {
			select {
			case <-quit:
				break LOOP
			case <-internal.Done():
				break LOOP
			}
		}
	}

	close(c.done)
	return nil
}

func (c *cli) Done() <-chan struct{} {
	return c.done
}

func (c *cli) handler(input string, cancel context.CancelFunc) error {
	switch input {
	case "q":
		cancel()
	}

	return nil
}
