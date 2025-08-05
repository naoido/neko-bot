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

	// キーイベントを受け取るため
	keyEvent := NewKeyEvent()
	input, t, err := keyEvent.Listening(internal)
	if err != nil {
		return err
	}
	defer func() {
		if t != nil {
			_ = term.Restore(int(os.Stdin.Fd()), t)
		}
	}()

	// OSシグナルを受信するため
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

LOOP:
	for {
		select {
		case cmd := <-input:
			err := c.handler(cmd, cancel)

			if err != nil {
				fmt.Println("CLI操作でエラーが発生しました :", err)
				cancel()
			}
		case <-quit:
		case <-internal.Done():
			break LOOP
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
