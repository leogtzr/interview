package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	selectedTopic = ""
	ps1           = "$ "
)

func main() {

	userInput := bufio.NewReader(os.Stdin)
	interviewTopicsDir := os.Getenv("INTERVIEW_DIR")
	if interviewTopicsDir == "" {
		log.Fatal("INTERVIEW_DIR environment variable not defined.")
	}

	for {
		fmt.Print(ps1String(ps1, selectedTopic))
		text, _ := userInput.ReadString('\n')
		text = strings.TrimSpace(text)
		if len(text) == 0 {
			continue
		}
		cmd, options := userInputToCmd(text)

		switch cmd {
		case exitCmd:
			fmt.Println("\tBye ... ")
			os.Exit(0)
		case topicsCmd:
			listTopics(interviewTopicsDir)
		case helpCmd:
			printHelp()
		case clearScreenCommand:
			clearScreen()
		case useCmd:
			setTopic(options)
		}
	}

}
