package interview

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
	if err := db.Close(); err != nil {
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

/*
func Test_setAnswerAsNeutral(t *testing.T) {
	qs := []Question{
		Question{ID: 1, Answer: OK},
		Question{ID: 2, Answer: OK},
		Question{ID: 3, Answer: Wrong},
		Question{ID: 4, Answer: NotAnsweredYet},
	}

	colorProfile := termenv.ColorProfile()
	config := Config{colorProfile: colorProfile}

	for i := 0; i < len(qs); i++ {
		setAnswerAsNeutral(&qs, &config)
		config.questionIndex++
	}

	for _, q := range qs {
		if q.Answer != Neutral {
			t.Errorf("%s should have been marked as Neutral.", q)
		}
	}
}
*/

/*
func Test_answerAs(t *testing.T) {
	topics := make(map[string][]Question)
	linuxQuestions := []Question{
		Question{ID: 1, Q: "lx1", Level: AssociateOrProgrammer, Answer: NotAnsweredYet},
		Question{ID: 2, Q: "lx2", Level: AssociateOrProgrammer, Answer: NotAnsweredYet},
		Question{ID: 3, Q: "lx3", Level: AssociateOrProgrammer, Answer: NotAnsweredYet},
		Question{ID: 4, Q: "lx4", Level: SrProgrammer, Answer: NotAnsweredYet},
		Question{ID: 5, Q: "lx5", Level: SrProgrammer, Answer: NotAnsweredYet},
	}

	javaQuestions := []Question{
		Question{ID: 1, Q: "j1", Level: AssociateOrProgrammer, Answer: NotAnsweredYet},
		Question{ID: 2, Q: "j2", Level: SrProgrammer, Answer: Wrong},
	}

	randomQuestions := []Question{}

	topics["linux"] = linuxQuestions
	topics["java"] = javaQuestions
	topics["random"] = randomQuestions

	config := Config{}
	config.levels = [3]Level{
		AssociateOrProgrammer, ProgrammerAnalyst, SrProgrammer,
	}
	config.selectedTopic = "linux"
	config.interview = Interview{Interviewee: "Hello", Date: time.Now(), Topics: topics}
	config.individualLevelIndexes = []int{0, 0, 0}

	answerAs(&config, OK, green)
	config.individualLevelIndexes[int(AssociateOrProgrammer)-1]++

	config.topicQuestionsLevel = SrProgrammer
	config.levelIndex = int(SrProgrammer) - 1
	answerAs(&config, OK, green)

	config.individualLevelIndexes[int(SrProgrammer)-1]++

	if config.interview.Topics["linux"][0].Answer != OK {
		t.Errorf("got=[%s], want=[%s]", config.interview.Topics["linux"][0], OK)
	}

	if config.interview.Topics["linux"][3].Answer != OK {
		t.Errorf("got=[%s], want=[%s]", config.interview.Topics["linux"][0], OK)
	}

}
*/
