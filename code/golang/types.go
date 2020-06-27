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

// UserAnswer ...
type UserAnswer struct {
	ID         int    `json:"id"`
	Result     int    `json:"result"`
	Comment    string `json:"comment"`
	QuestionID int    `json:"question_id"`
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

// Topic ...
type Topic struct {
	ID    int    `json:"id"`
	Topic string `json:"topic"`
}
