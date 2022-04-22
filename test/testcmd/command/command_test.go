package command

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
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
		Command: "../test/test.exe",
		Ctx:     ctx,

		ErrFunc: func(msg []byte) (int, error) {
			fmt.Print(string(msg))
			return len(msg), nil
		},
		OutFunc: func(msg []byte) (int, error) {
			fmt.Print(string(msg))
			return len(msg), nil
		},
		InFunc: func(w io.WriteCloser) (int, error) {
			var line string
			scanner := bufio.NewScanner(os.Stdin)

			if scanner.Scan() {
				line = scanner.Text()
			}

			if scanner.Err() != nil {
				fmt.Println(scanner.Err())
			}

			return w.Write([]byte(line))
		},
	}

	err := Exec(&info)
	if err != nil {
		fmt.Printf("failed to execute command [error=%s]\n", err.Error())
		t.FailNow()
	}
}
