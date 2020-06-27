package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func setup() {
	dbUser, ok := os.LookupEnv("DB_INTERVIEW_USER")
	if !ok {
		panic("DB_INTERVIEW_USER variable not set")
	}
	dbPassword, ok := os.LookupEnv("DB_INTERVIEW_PASSWORD")
	if !ok {
		panic("DB_INTERVIEW_PASSWORD variable not set")
	}
	dbName, ok := os.LookupEnv("DB_TEST_INTERVIEW_NAME")
	if !ok {
		panic("DB_TEST_INTERVIEW_NAME variable not set")
	}

	jdbcURL := fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName)
	var err error
	db, err = sql.Open("mysql", jdbcURL)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 20)
	// Do something here.
	fmt.Printf("\033[1;36m%s\033[0m", "> Setup completed\n")
}

func teardown() {
	// Do something here.
	err := db.Close()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\033[1;36m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func Test_getTopics(t *testing.T) {
	topics, err := getTopics(db)
	if err != nil {
		t.Error(err.Error())
	}
	if len(topics) < 1 {
		t.Error("Expecting topics in DB.")
	}
}

func Test_getQuestionsByTopicWithLevel(t *testing.T) {
	topics, err := getQuestionsByTopicWithLevel("java", SrProgrammer, db)
	if err != nil {
		t.Error(err.Error())
	}
	if len(topics) < 1 {
		t.Errorf("Expecting topics in DB: %d\n", len(topics))
	}
}

/*
// TODO: fix this ...
func Test_setAnswerAsOK(t *testing.T) {
	qs := []Question{
		Question{ID: 1, Answer: OK},
		Question{ID: 2, Answer: OK},
		Question{ID: 3, Answer: Wrong},
		Question{ID: 4, Answer: NotAnsweredYet},
	}

	colorProfile := termenv.ColorProfile()
	config := Config{colorProfile: colorProfile}

	for i := 0; i < len(qs); i++ {
		setAnswerAsOK(&qs, &config)
		config.questionIndex++
	}

	for _, q := range qs {
		if q.Answer != OK {
			t.Errorf("%s should have been marked as OK.", q)
		}
	}
}
*/

/*
func Test_setAnswerAsWrong(t *testing.T) {
	qs := []Question{
		Question{ID: 1, Answer: OK},
		Question{ID: 2, Answer: OK},
		Question{ID: 3, Answer: Wrong},
		Question{ID: 4, Answer: NotAnsweredYet},
	}

	colorProfile := termenv.ColorProfile()
	config := Config{colorProfile: colorProfile}

	for i := 0; i < len(qs); i++ {
		setAnswerAsWrong(&qs, &config)
		config.questionIndex++
	}

	for _, q := range qs {
		if q.Answer != Wrong {
			t.Errorf("%s should have been marked as Wrong.", q)
		}
	}
}
*/
