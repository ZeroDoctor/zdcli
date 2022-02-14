package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("start")

	time.Sleep(1 * time.Second)
	fmt.Println("testing...")
	time.Sleep(1 * time.Second)
	fmt.Println("out...")
	time.Sleep(1 * time.Second)
	fmt.Println("output...")

	fmt.Print("Enter something:")

	var line string
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() {
		line = scanner.Text()
	}

	if scanner.Err() != nil {
		fmt.Println(scanner.Err())
	}

	fmt.Println("got:", line)

	fmt.Println("see ya!")
}
