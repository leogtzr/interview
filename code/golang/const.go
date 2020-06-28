package main

const (
	red     = "#E88388"
	green   = "#A8CC8C"
	yellow  = "#DBAB79"
	blue    = "#71BEF2"
	magenta = "#D290E4"
	cyan    = "#66C2CD"
	gray    = "#B9BFCA"
)

const (
	minNumberOfCharsInIntervieweeName = 10
	interviewFormatLayout             = "2006-01-2 15:04:05"
)

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
	requiredNumberOfFieldsInInterviewHeaderRecord = 2
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
