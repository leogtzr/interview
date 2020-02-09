package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/muesli/termenv"
)

var (
	selectedTopic           = ""
	ps1                     = "$ "
	interviewTopicsDir      = ""
	hasStarted         bool = false
	questionIndex           = 0
	colorProfile            = termenv.ColorProfile()
	rgxQuestions            = regexp.MustCompile("^\\d+@.+@(\\d+)?$")
	interview               = Interview{Topics: make(map[string]Questions)}
)

const (
	red                               = "#E88388"
	green                             = "#A8CC8C"
	yellow                            = "#DBAB79"
	blue                              = "#71BEF2"
	magenta                           = "#D290E4"
	cyan                              = "#66C2CD"
	gray                              = "#B9BFCA"
	minNumberOfCharsInIntervieweeName = 10
)

func main() {

	userInput := bufio.NewReader(os.Stdin)
	interviewTopicsDir = os.Getenv("INTERVIEW_DIR")
	if interviewTopicsDir == "" {
		log.Fatal("INTERVIEW_DIR environment variable not defined.")
	}

	for {
		fmt.Print(ps1String(ps1, selectedTopic, interview.Interviewee))
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
			fmt.Println(termenv.String(selectedTopic).Bold())
		case useCmd:
			setTopic(options)
		case startCmd:
			if hasStarted {
				printWithColorln("Interview has already started.", yellow)
				break
			}
			fmt.Printf("Interviewee name: ")
			if name, ok := readIntervieweeName(os.Stdin); !ok {
				break
			} else {
				interview.Interviewee = name
				interview.Date = time.Now()
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
			if !hasStarted {
				printWithColorln("Interview has not yet started.", yellow)
				break
			}
			interview.Topics[selectedTopic][questionIndex].Answer = OK
			printWithColorln(fmt.Sprintf("Answer has saved as '%s'", OK), green)

		case wrongAnswerCmd:
			if !hasStarted {
				printWithColorln("Interview has not yet started.", yellow)
				break
			}

			interview.Topics[selectedTopic][questionIndex].Answer = Wrong
			printWithColorln(fmt.Sprintf("Answer has saved as '%s'", Wrong), red)

		case mehAnswerCmd:
			if !hasStarted {
				printWithColorln("Interview has not yet started.", yellow)
				break
			}
			interview.Topics[selectedTopic][questionIndex].Answer = Neutral
			printWithColorln(fmt.Sprintf("Answer has saved as '%s'", Neutral), magenta)

		case finishCmd:
			//TODO: pending ...
			break
		}
	}

}
