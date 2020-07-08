package interview

import (
	"database/sql"
	"fmt"
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

// ResultCount ...
type ResultCount struct {
	Result, Count int
}

func (rc ResultCount) String() string {
	return fmt.Sprintf("{Result: %s, count: %d}", Result(rc.Result), rc.Count)
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

// AnswerView ...
type AnswerView struct {
	ID       int
	Question string
	Result   int
	Comment  sql.NullString
	Topic    string
	Title    string
}

func (av AnswerView) String() string {
	hasComment := av.Comment.Valid
	if hasComment {
		return fmt.Sprintf("%s [%s] [%s] [%s] [%s]",
			av.Question, Result(av.Result), av.Comment.String, av.Topic, av.Title)
	}
	return fmt.Sprintf("%s [%s] [%s] [%s]",
		av.Question, Result(av.Result), av.Topic, av.Title)
}

// CandidateView ...
type CandidateView struct {
	ID   int
	Name string
	Date string
}

func (can CandidateView) String() string {
	return fmt.Sprintf("%d - \"%s\" (%s)", can.ID, can.Name, can.Date)
}
