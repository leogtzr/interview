package main

import (
	"time"

	"github.com/muesli/termenv"
)

// Config ...
type Config struct {
	selectedTopic          string
	ps1                    string
	hasStarted             bool
	questionIndex          int
	topicQuestionsLevel    Level
	levelIndex             int
	ignoreLevelChecking    bool
	individualLevelIndexes []int
	levels                 [3]Level
	colorProfile           termenv.Profile
	interview              Interview
	intervieweeID          int
	comment                string
}

// Command ...
type Command int

// Question ...
type Question struct {
	ID     int
	Q      string
	Answer string
	Result Result
	Level  Level
}

// Result ...
type Result int

// Level ...
type Level int

// Interview ...
type Interview struct {
	Interviewee string
	Date        time.Time
	Topics      map[string][]Question
}

// Topic ...
type Topic struct {
	ID    int    `json:"id"`
	Topic string `json:"topic"`
}
