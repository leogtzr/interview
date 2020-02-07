package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/muesli/termenv"
)

var (
	selectedTopic      = ""
	ps1                = "$ "
	interviewTopicsDir = ""
	colorProfile       = termenv.ColorProfile()
	rgxQuestions       = regexp.MustCompile("^\\d+@.+@(\\d+)?$")
	questionsPerTopic  []Question
	hasStarted         bool   = false
	questionIndex             = 0
	intervieweeName    string = ""
	red                       = "#E88388"
	green                     = "#A8CC8C"
	yellow                    = "#DBAB79"
	blue                      = "#71BEF2"
	magenta                   = "#D290E4"
	cyan                      = "#66C2CD"
	gray                      = "#B9BFCA"
)

func main() {

	userInput := bufio.NewReader(os.Stdin)
	interviewTopicsDir = os.Getenv("INTERVIEW_DIR")
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
		case pwdCommand:
			printWorkingDirectory()
		case useCmd:
			setTopic(options)
		case startCmd:
			if name, ok := readIntervieweeName(); !ok {
				break
			} else {
				intervieweeName = name
			}
			hasStarted = true
			questionIndex = 0
			printQuestion(questionIndex)
		case printCmd:
			printQuestion(questionIndex)
		case nextQuestionCmd:
			gotoNextQuestion()
			printQuestion(questionIndex)
		case previousQuestionCmd:
			gotoPreviousQuestion()
			printQuestion(questionIndex)
		case viewCmd:
			viewStats()
		case rightAnswerCmd:
			markAnswerAsOK()
		case wrongAnswerCmd:
			markAnswerAsWrong()
		case mehAnswerCmd:
			markAnswerAsNeutral()
		}
	}

}
