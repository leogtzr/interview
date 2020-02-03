package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	userInput := bufio.NewReader(os.Stdin)
	interviewTopicsDir := os.Getenv("INTERVIEW_DIR")
	if interviewTopicsDir == "" {
		log.Fatal("INTERVIEW_DIR environment variable not defined.")
	}

	for {
		fmt.Print("$ ")
		text, _ := userInput.ReadString('\n')
		text = strings.TrimSpace(text)
		if len(text) == 0 {
			continue
		}
		cmd := userInputToCmd(text)

		switch cmd {
		case exitCmd:
			fmt.Println("\tBye ... ")
			os.Exit(0)
		case topicsCmd:
			listTopics(interviewTopicsDir)
		case helpCmd:
			printHelp()
		}

	}

}
