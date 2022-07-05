package command

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"
)

var TF_Ping = flag.Bool("ping", false, "enable ping testing")

func TestCurlCommand(t *testing.T) {
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

func TestPingCommand(t *testing.T) {
	if TF_Ping == nil || !*TF_Ping {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info := Info{
		// Command: "lua ../lua/build-app.lua test get_name",
		// Dir:     "../lua",
		Command: "ping google.com",
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

			return "", ErrStdInNone
		},
	}

	err := Exec(&info)
	if err != nil {
		fmt.Printf("failed to execute command [error=%s]\n", err.Error())
		t.FailNow()
	}
}
