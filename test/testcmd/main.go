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

func main() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	reader := bufio.NewReader(os.Stdin)
	info := command.Info{
		Command: "../test.exe",
		Args: []string{
			"",
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
		InFunc: func(w io.WriteCloser, done <-chan struct{}) (int, error) {
			var line string

			line, err = reader.ReadString('\n')
			if err != nil {
				return 0, err
			}

			return io.WriteString(w, line)
		},
	}

	err = command.Exec(&info)
	if err != nil {
		fmt.Println("[ERROR]", err.Error())
	}

	fmt.Println("end program")
}
