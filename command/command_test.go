package command

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestPingCommand(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info := Info{
		Command: "curl www.google.com",
		Ctx:     ctx,

		ErrFunc: func(msg []byte) (int, error) {
			fmt.Print(string(msg))
			return len(msg), nil
		},
		OutFunc: func(msg []byte) (int, error) {
			fmt.Print(string(msg))
			return len(msg), nil
		},
	}

	err := Exec(&info)
	if err != nil {
		fmt.Printf("failed to execute command [error=%s]\n", err.Error())
		t.FailNow()
	}
}

func TestStdInCommand(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info := Info{
		Command: "lua ../lua/build-app.lua test get_name",
		Dir:     "../lua",
		Ctx:     ctx,

		ErrFunc: func(msg []byte) (int, error) {
			fmt.Print(string(msg))
			return len(msg), nil
		},
		OutFunc: func(msg []byte) (int, error) {
			fmt.Print(string(msg))
			return len(msg), nil
		},
		InFunc: func(ctx context.Context) (string, error) {
			time.Sleep(500 * time.Millisecond)

			var line string

			line = "test"

			return line, nil
		},
	}

	err := Exec(&info)
	if err != nil {
		fmt.Printf("failed to execute command [error=%s]\n", err.Error())
		t.FailNow()
	}
}
