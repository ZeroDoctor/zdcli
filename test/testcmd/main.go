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

	// var proc *sys.Process
	// proc, err = sys.NewProcess("ping www.google.com", nil, "")
	// if err != nil {
	// 	fmt.Println("[ERROR]", err.Error())
	// 	return
	// }
	//
	// done := make(chan bool, 1)
	// r, w, err := os.Pipe()
	// if err != nil {
	// 	fmt.Println("[ERROR]", err.Error())
	// 	return
	// }
	//
	// proc.SetInWriter(os.Stdin)
	// proc.SetOutReader(w)
	// proc.SetErrReader(os.Stderr)
	//
	// err = proc.Start()
	// if err != nil {
	// 	fmt.Println("[ERROR]", err.Error())
	// 	return
	// }
	//
	// // reader := bufio.NewReader(r)
	//
	// time.Sleep(1 * time.Second)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-done:
	// 			r.Close()
	// 			return
	// 		default:
	// 		}
	//
	// 		var buf bytes.Buffer
	// 		_, err := io.Copy(&buf, r)
	// 		if err != nil {
	// 			fmt.Println("[ERROR]", err.Error())
	// 		}
	//
	// 		fmt.Println(buf.String())
	// 	}
	// }()
	//
	// state, err := proc.Wait()
	// if err != nil {
	// 	fmt.Println("[ERROR]", err.Error())
	// 	return
	// }
	// err = w.Close()
	// if err != nil {
	// 	fmt.Println("[ERROR]", err.Error())
	// }
	//
	// close(done)
	// fmt.Println(state)

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
		InFunc: func(w io.WriteCloser) (int, error) {
			var line string

			fmt.Println("waiting for response")
			line, err = reader.ReadString('\n')
			if err != nil {
				return 0, err
			}

			fmt.Println("done here")

			return io.WriteString(w, line+"\r\n")
		},
	}

	err = command.Exec(&info)
	if err != nil {
		fmt.Println("[ERROR]", err.Error())
	}
}
