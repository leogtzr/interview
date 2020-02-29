package main

import (
	"regexp"
	"time"

	"github.com/muesli/termenv"
)

// Config ...
type Config struct {
	selectedTopic          string
	ps1                    string
	interviewTopicsDir     string
	hasStarted             bool
	questionIndex          int
	usingInterviewFile     bool
	topicQuestionsLevel    Level
	levelIndex             int
	ignoreLevelChecking    bool
	individualLevelIndexes []int
	levels                 [3]Level
	colorProfile           termenv.Profile
	rgxQuestions           regexp.Regexp
	interview              Interview
}

// Command ...
type Command int

// Question ...
type Question struct {
	ID     int
	Q      string
	Answer Answer
	Level  Level
}

// Answer ...
type Answer int

// Level ...
type Level int

// Interview ...
type Interview struct {
	Interviewee string
	Date        time.Time
	Topics      map[string][]Question
}

const (
	// NotAnsweredYet ...
	NotAnsweredYet Answer = 1
	// OK ...
	OK Answer = 2
	// Wrong ...
	Wrong Answer = 3
	// Neutral ...
	Neutral Answer = 4
)

const (
	// AssociateOrProgrammer ...
	AssociateOrProgrammer Level = 1
	// ProgrammerAnalyst ...
	ProgrammerAnalyst Level = 2
	// SrProgrammer ...
	SrProgrammer Level = 3
	// All ...
	All Level = 4
)

// Commands:
const (
	exitCmd                        Command = iota
	exitInterviewFileCmd           Command = iota
	topicsCmd                      Command = iota
	helpCmd                        Command = iota
	useCmd                         Command = iota
	clearScreenCmd                 Command = iota
	pwdCmd                         Command = iota
	noCmd                          Command = iota
	startCmd                       Command = iota
	printCmd                       Command = iota
	nextQuestionCmd                Command = iota
	previousQuestionCmd            Command = iota
	viewCmd                        Command = iota
	rightAnswerCmd                 Command = iota
	wrongAnswerCmd                 Command = iota
	mehAnswerCmd                   Command = iota
	finishCmd                      Command = iota
	loadCmd                        Command = iota
	increaseLevelCmd               Command = iota
	decreaseLevelCmd               Command = iota
	ignoreLevelCmd                 Command = iota
	showLevelCmd                   Command = iota
	showStatsCmd                   Command = iota
	setAssociateProgrammerLevelCmd Command = iota
	setProgrammerAnalystLevelCmd   Command = iota
	setSRProgrammerLevelCmd        Command = iota
	validateQuestionsCmd           Command = iota
	countCmd                       Command = iota
)
