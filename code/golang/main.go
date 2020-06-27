// TODO: save date in the DB
// TODO: code to save comments ... (open a dialog to write up the comment or read it from stdin)
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
			err = listTopics(db)
			if err != nil {
				panic(err)
			}
		case helpCmd:
			printHelp()
		case clearScreenCmd:
			clearScreen()
		case pwdCmd:
			fmt.Println(termenv.String(config.selectedTopic).Bold())
		case useCmd:
			err = setTopic(options, &config, db)
			if err != nil {
				panic(err)
			}
		case startCmd:
			if config.hasStarted {
				printWithColorln("Interview has already started.", yellow, &config)
				break
			}
			fmt.Printf("Interviewee name: ")
			if name, ok := readIntervieweeName(os.Stdin); !ok {
				break
			} else {
				// Persist the info in the DB ...
				id, err := saveIntervieweeName(name, db)
				if err != nil {
					panic(err)
				}
				config.intervieweeID = id
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
				err := setAnswerAsOK(&config, db)
				if err != nil {
					panic(err)
				}
			} else {
				err := answerAs(&config, OK, green, db)
				if err != nil {
					panic(err)
				}
			}

		case wrongAnswerCmd:
			if !config.hasStarted {
				printWithColorln("Interview has not yet started.", yellow, &config)
				break
			}

			if config.ignoreLevelChecking {
				err := setAnswerAsWrong(&config, db)
				if err != nil {
					panic(err)
				}
			} else {
				err := answerAs(&config, Wrong, red, db)
				if err != nil {
					panic(err)
				}
			}

		case mehAnswerCmd:
			if !config.hasStarted {
				printWithColorln("Interview has not yet started.", yellow, &config)
				break
			}

			if config.ignoreLevelChecking {
				err := setAnswerAsNeutral(&config, db)
				if err != nil {
					panic(err)
				}
			} else {
				err := answerAs(&config, Neutral, yellow, db)
				if err != nil {
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
			showStats(&config)
		case setAssociateProgrammerLevelCmd:
			setLevel(AssociateOrProgrammer, &config)
		case setProgrammerAnalystLevelCmd:
			setLevel(ProgrammerAnalyst, &config)
		case setSRProgrammerLevelCmd:
			setLevel(SrProgrammer, &config)
		case countCmd:
			showCounts(&config)
		}
	}

}
