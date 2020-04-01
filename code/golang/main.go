package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/muesli/termenv"
)

func main() {

	config := NewConfig()
	config.interviewTopicsDir = os.Getenv("INTERVIEW_DIR")
	if config.interviewTopicsDir == "" {
		log.Fatal("INTERVIEW_DIR environment variable not defined.")
	}

	userInput := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(ps1String(config.ps1, config.selectedTopic, config.interview.Interviewee))
		text, _ := userInput.ReadString('\n')
		text = strings.TrimSpace(text)
		if len(text) == 0 {
			continue
		}
		cmd, options := userInputToCmd(text)

		switch cmd {
		case exitCmd:
			printWithColorln("Bye", magenta, &config)
			os.Exit(0)
		case exitInterviewFileCmd:
			printWithColorln("Exiting from interview file ... ", gray, &config)
			resetStatus(&config)
			break
		case topicsCmd:
			if config.usingInterviewFile {
				listTopicsFromInterviewFile(&config.interview.Topics, &config)
				break
			}
			listTopics(config.interviewTopicsDir)
		case helpCmd:
			printHelp()
		case clearScreenCmd:
			clearScreen()
		case pwdCmd:
			fmt.Println(termenv.String(config.selectedTopic).Bold())
		case useCmd:
			if config.usingInterviewFile {
				setTopicFrom(options, &config.interview.Topics, &config)
				break
			}
			setTopicFromFileSystem(options, &config)
		case startCmd:
			if config.hasStarted {
				printWithColorln("Interview has already started.", yellow, &config)
				break
			}
			fmt.Printf("Interviewee name: ")
			if name, ok := readIntervieweeName(os.Stdin); !ok {
				break
			} else {
				config.interview.Interviewee = name
				config.interview.Date = time.Now()
			}
			config.hasStarted = true
			printQuestion(config.questionIndex, &config)
		case printCmd:
			printQuestion(config.questionIndex, &config)
		case nextQuestionCmd:
			gotoNextQuestion(&config)
			printQuestion(config.questionIndex, &config)
		case previousQuestionCmd:
			gotoPreviousQuestion(&config)
			printQuestion(config.questionIndex, &config)
		case viewCmd:
			if !config.ignoreLevelChecking {
				viewQuestionsByLevel(&config)
			} else {
				viewQuestions(&config)
			}
		case rightAnswerCmd:
			if !config.hasStarted {
				printWithColorln("Interview has not yet started.", yellow, &config)
				break
			}
			if config.ignoreLevelChecking {
				qs := config.interview.Topics[config.selectedTopic]
				setAnswerAsOK(&qs, &config)
			} else {
				answerAs(&config, OK, green)
			}

		case wrongAnswerCmd:
			if !config.hasStarted {
				printWithColorln("Interview has not yet started.", yellow, &config)
				break
			}

			if config.ignoreLevelChecking {
				qs := config.interview.Topics[config.selectedTopic]
				setAnswerAsWrong(&qs, &config)
			} else {
				answerAs(&config, Wrong, red)
			}

		case mehAnswerCmd:
			if !config.hasStarted {
				printWithColorln("Interview has not yet started.", yellow, &config)
				break
			}

			if config.ignoreLevelChecking {
				qs := config.interview.Topics[config.selectedTopic]
				setAnswerAsNeutral(&qs, &config)
			} else {
				answerAs(&config, Neutral, yellow)
			}

		case finishCmd:
			err := saveInterview(&config)
			if err != nil {
				panic(err)
			}
			printWithColorln(fmt.Sprintf("Interview for '%s' has been saved.\n\n\tBye ...", config.interview.Interviewee), green, &config)
			os.Exit(1)

		case loadCmd:
			interviewFromFile, err := loadInterview(options, &config)
			if err != nil {
				printWithColorln(err.Error(), red, &config)
				break
			}

			config.usingInterviewFile = true
			printWithColorln("You will now be navigating through an interview file.", green, &config)

			config.interview = interviewFromFile

			for topic, questions := range interviewFromFile.Topics {
				fmt.Printf("[%s]\n", topic)
				for _, q := range questions {
					fmt.Println(q.String())
				}
			}

		case increaseLevelCmd:
			increaseLevel(&config)
		case decreaseLevelCmd:
			decreaseLevel(&config)
		case ignoreLevelCmd:
			toggleLevelChecking(&config)
		case showLevelCmd:
			showLevel(&config)
		case showStatsCmd:
			showStats(&config)
		case setAssociateProgrammerLevelCmd:
			setLevel(AssociateOrProgrammer, &config)
		case setProgrammerAnalystLevelCmd:
			setLevel(ProgrammerAnalyst, &config)
		case setSRProgrammerLevelCmd:
			setLevel(SrProgrammer, &config)
		case validateQuestionsCmd:
			validateQuestions(&config)
		case countCmd:
			showCounts(&config)
		case notesCmd:
			err := createNotes(&config)
			if err != nil {
				printWithColorln(err.Error(), yellow, &config)
			}
		}
	}

}
