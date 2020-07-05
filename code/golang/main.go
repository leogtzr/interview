package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/muesli/termenv"
)

func main() {

	config := NewConfig()

	dbConfig, err := readConfig("interviews.env", os.Getenv("HOME"), map[string]interface{}{
		"db_user":     os.Getenv("DB_INTERVIEW_USER"),
		"db_password": os.Getenv("DB_INTERVIEW_PASSWORD"),
		"db_name":     os.Getenv("DB_INTERVIEW_NAME"),
		"db_driver":   os.Getenv("DB_DRIVER"),
	})
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	// DB setup ...
	jdbcURL := fmt.Sprintf("%s:%s@/%s", dbConfig.GetString("db_user"), dbConfig.GetString("db_password"), dbConfig.GetString("db_name"))
	db, err := sql.Open(dbConfig.GetString("db_driver"), jdbcURL)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Hour * 3)
	defer db.Close()

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
		case topicsCmd:
			if err = listTopics(db); err != nil {
				panic(err)
			}
		case helpCmd:
			printHelp()
		case clearScreenCmd:
			clearScreen()
		case pwdCmd:
			fmt.Println(termenv.String(config.selectedTopic).Bold())
		case useCmd:
			if err = setTopic(options, &config, db); err != nil {
				panic(err)
			}
		case startCmd:
			if config.hasStarted {
				printWithColorln("Interview has already started.", yellow, &config)
				break
			}

			if len(config.selectedTopic) == 0 {
				printWithColorln("You need to select a topic first.", red, &config)
				break
			}

			fmt.Printf("Interviewee name: ")

			name, ok := readIntervieweeName(os.Stdin)
			if !ok {
				break
			}
			id, err := saveIntervieweeName(name, db)
			if err != nil {
				panic(err)
			}
			config.intervieweeID = id
			config.interview.Interviewee = name
			config.interview.Date = time.Now()
			config.hasStarted = true
			// Message to the user that the interview has started.
			printQuestion(config.questionIndex, &config)
		case printCmd:
			printQuestion(config.questionIndex, &config)
		case nextQuestionCmd:
			config.comment = ""
			gotoNextQuestion(&config)
			printQuestion(config.questionIndex, &config)
		case previousQuestionCmd:
			config.comment = ""
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
				if err := setAnswerAsOK(&config, db); err != nil {
					panic(err)
				}
			} else {
				if err := answerAs(&config, OK, green, db); err != nil {
					panic(err)
				}
			}

		case wrongAnswerCmd:
			if !config.hasStarted {
				printWithColorln("Interview has not yet started.", yellow, &config)
				break
			}

			if config.ignoreLevelChecking {
				if err := setAnswerAsWrong(&config, db); err != nil {
					panic(err)
				}
			} else {
				if err := answerAs(&config, Wrong, red, db); err != nil {
					panic(err)
				}
			}

		case mehAnswerCmd:
			if !config.hasStarted {
				printWithColorln("Interview has not yet started.", yellow, &config)
				break
			}

			if config.ignoreLevelChecking {
				if err := setAnswerAsNeutral(&config, db); err != nil {
					panic(err)
				}
			} else {
				if err := answerAs(&config, Neutral, yellow, db); err != nil {
					panic(err)
				}
			}

		case finishCmd:
			printWithColorln(fmt.Sprintf("Interview for '%s' has been saved.\n\n\tBye ...", config.interview.Interviewee), green, &config)
			os.Exit(0)
		case increaseLevelCmd:
			increaseLevel(&config)
		case decreaseLevelCmd:
			decreaseLevel(&config)
		case ignoreLevelCmd:
			toggleLevelChecking(&config)
		case showLevelCmd:
			showLevel(&config)
		case showStatsCmd:
			if err := showStats(&config, db); err != nil {
				panic(err)
			}
		case setAssociateProgrammerLevelCmd:
			setLevel(AssociateOrProgrammer, &config)
		case setProgrammerAnalystLevelCmd:
			setLevel(ProgrammerAnalyst, &config)
		case setSRProgrammerLevelCmd:
			setLevel(SrProgrammer, &config)
		case countCmd:
			showCounts(&config)
		case createCommentCmd:
			fmt.Printf("Comment: ")
			comment, err := readComment()
			if err != nil {
				panic(err)
			}
			config.comment = comment
		case createQuestionCmd:
			if err := makeQuestion(&config, db); err != nil {
				panic(err)
			}
			printWithColorln("Question created", magenta, &config)
		case viewAnwswerCmd:
			viewAnswer(config.questionIndex, &config)
		}
	}

}
