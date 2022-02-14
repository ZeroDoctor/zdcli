package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"testcmd/command"
	"time"
)

// TODO: look into creating a window/linux process with syscall instead

func main() {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info := command.Info{
		Command: "bash -c",
		Args: []string{
			"/mnt/c/Users/Daniel/Documents/zerodoc/zdcli/test/test",
		},
		Ctx: ctx,

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
				return 0, fmt.Errorf("failed at scanner [error=%w]", scanner.Err())
			}

			return io.WriteString(w, line)
		},
	}

	err = command.Exec(&info)
	if err != nil {
		fmt.Println("[ERROR]", err.Error())
	}
}
