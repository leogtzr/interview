package main

import (
	"bufio"
	"fmt"
	"os"
)

// Command ...
type Command int

const (
	exitCmd   Command = iota
	topicsCmd Command = iota
	noCmd     Command = iota
)

func main() {

	userInput := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		text, _ := userInput.ReadString('\n')
		cmd := userInputToCmd(text)
		fmt.Println(cmd)
	}

}
