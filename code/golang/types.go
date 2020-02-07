package main

// Command ...
type Command int

// Question ...
type Question struct {
	ID             int
	Q              string
	NextQuestionID int
	Answer         Answer
}

// Answer ...
type Answer int

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
	exitCmd             Command = iota
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
)
