package main

import (
	"regexp"
	"time"

	"github.com/muesli/termenv"
)

// Config ...
type Config struct {
	selectedTopic      string
	ps1                string
	interviewTopicsDir string
	hasStarted         bool
	questionIndex      int
	// usingInterviewFile     bool
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

// Topic ...
type Topic struct {
	ID    int    `json:"id"`
	Topic string `json:"topic"`
}
