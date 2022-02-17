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

	fmt.Print("Enter something: ")

	var line string
	reader := bufio.NewReader(os.Stdin)

	line, _ = reader.ReadString('\n')

	fmt.Println("got:", line)
	fmt.Println("okay try again in a sec")
	time.Sleep(1 * time.Second)

	fmt.Print("okay go: ")
	reader = bufio.NewReader(os.Stdin)
	line, _ = reader.ReadString('\n')

	fmt.Println("nice, just got:", line)

	fmt.Print("again: ")
	reader = bufio.NewReader(os.Stdin)
	line, _ = reader.ReadString('\n')
	fmt.Println("alright this is it:", line)

	fmt.Println("see ya!")
}
