package main

import "time"

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
	exitCmd             Command = iota
	exitInterviewFile   Command = iota
	topicsCmd           Command = iota
	helpCmd             Command = iota
	useCmd              Command = iota
	clearScreenCommand  Command = iota
	pwdCommand          Command = iota
	noCmd               Command = iota
	startCmd            Command = iota
	printCmd            Command = iota
	nextQuestionCmd     Command = iota
	previousQuestionCmd Command = iota
	viewCmd             Command = iota
	rightAnswerCmd      Command = iota
	wrongAnswerCmd      Command = iota
	mehAnswerCmd        Command = iota
	finishCmd           Command = iota
	loadCmd             Command = iota
	increaseLevelCmd    Command = iota
	decreaseLevelCmd    Command = iota
	ignoreLevelCmd      Command = iota
)
