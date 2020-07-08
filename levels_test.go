package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/muesli/termenv"
)

func Test_getQuestionsFromLevel(t *testing.T) {
	type test struct {
		lvl       Level
		config    Config
		topic     string
		questions []Question
	}

	topics := make(map[string][]Question)
	linuxQuestions := []Question{
		Question{ID: 1, Q: "lx1", Level: AssociateOrProgrammer, Result: NotAnsweredYet},
		Question{ID: 2, Q: "lx2", Level: AssociateOrProgrammer, Result: NotAnsweredYet},
		Question{ID: 3, Q: "lx3", Level: ProgrammerAnalyst, Result: OK},
		Question{ID: 4, Q: "lx4", Level: SrProgrammer, Result: OK},
	}

	javaQuestions := []Question{
		Question{ID: 1, Q: "j1", Level: AssociateOrProgrammer, Result: NotAnsweredYet},
		Question{ID: 2, Q: "j2", Level: SrProgrammer, Result: Wrong},
	}

	randomQuestions := []Question{}

	topics["linux"] = linuxQuestions
	topics["java"] = javaQuestions
	topics["random"] = randomQuestions

	config := Config{}
	config.interview = Interview{Interviewee: "Hello", Date: time.Now(), Topics: topics}

	tests := []test{
		{
			lvl:   AssociateOrProgrammer,
			topic: "linux",
			questions: []Question{
				Question{ID: 1, Q: "lx1", Level: AssociateOrProgrammer, Result: NotAnsweredYet},
				Question{ID: 2, Q: "lx2", Level: AssociateOrProgrammer, Result: NotAnsweredYet},
			},
			config: config,
		},

		{
			lvl:   SrProgrammer,
			topic: "java",
			questions: []Question{
				Question{ID: 2, Q: "j2", Level: SrProgrammer, Result: Wrong},
			},
			config: config,
		},

		{
			lvl:       SrProgrammer,
			topic:     "random",
			questions: []Question{},
			config:    config,
		},
	}

	for _, tt := range tests {
		tt.config.selectedTopic = tt.topic
		if got := getQuestionsFromLevel(tt.lvl, &tt.config); !Equal(got, tt.questions) {
			t.Errorf("got=[%s], want=[%s]", got, tt.questions)
		}
	}
}

func Test_increaseLevel(t *testing.T) {

	colorProfile := termenv.ColorProfile()

	type test struct {
		config             Config
		lvls               [3]Level
		incraseTimes, want int
	}

	levels := [3]Level{
		AssociateOrProgrammer, ProgrammerAnalyst, SrProgrammer,
	}

	tests := []test{
		{config: Config{colorProfile: colorProfile, levelIndex: 0, levels: levels}, incraseTimes: 1, want: 1},
		{config: Config{colorProfile: colorProfile, levelIndex: 0, levels: levels}, incraseTimes: 10, want: len(levels) - 1},
	}

	for _, tt := range tests {
		for i := 0; i < tt.incraseTimes; i++ {
			increaseLevel(&tt.config)
		}
		if tt.config.levelIndex != tt.want {
			t.Errorf("want=[%d] after increasing %d times, current index is: %d", tt.want, tt.incraseTimes, tt.config.levelIndex)
		}
	}
}

func Test_decreaseLevel(t *testing.T) {

	colorProfile := termenv.ColorProfile()

	type test struct {
		config              Config
		lvls                [3]Level
		decreaseTimes, want int
	}

	interviewLevels := [3]Level{
		AssociateOrProgrammer, ProgrammerAnalyst, SrProgrammer,
	}

	tests := []test{
		{config: Config{colorProfile: colorProfile, levelIndex: len(interviewLevels), levels: interviewLevels}, decreaseTimes: 2, want: 1},
		{config: Config{colorProfile: colorProfile, levelIndex: len(interviewLevels), levels: interviewLevels}, decreaseTimes: 1, want: len(interviewLevels) - 1},
		{config: Config{colorProfile: colorProfile, levelIndex: len(interviewLevels), levels: interviewLevels}, decreaseTimes: 2, want: len(interviewLevels) - 2},
		{config: Config{colorProfile: colorProfile, levelIndex: 0, levels: interviewLevels}, decreaseTimes: 10, want: 0},
	}

	for _, tt := range tests {
		for i := 0; i < tt.decreaseTimes; i++ {
			decreaseLevel(&tt.config)
		}
		if tt.config.levelIndex != tt.want {
			t.Errorf("want=[%d] after decreasing %d times, current index is: %d", tt.want, tt.decreaseTimes, tt.config.levelIndex)
		}
	}
}

func Test_toggleLevelChecking(t *testing.T) {
	type test struct {
		config Config
		want   bool
	}

	tests := []test{
		{config: Config{ignoreLevelChecking: false}, want: true},
		{config: Config{ignoreLevelChecking: true}, want: false},
	}

	for _, tt := range tests {
		toggleLevelChecking(&tt.config)
		if tt.config.ignoreLevelChecking != tt.want {
			t.Errorf("the check should be %t after change. Currently is: %t", tt.want, tt.config.ignoreLevelChecking)
		}
	}
}

func Test_findLevel(t *testing.T) {
	questions := []Question{
		Question{Level: SrProgrammer, Q: "Q1"},
		Question{Level: SrProgrammer, Q: "Q2"},
		Question{Level: SrProgrammer, Q: "Q3"},
	}

	lvl := findLevel(&questions, ProgrammerAnalyst)

	if lvl != AssociateOrProgrammer {
		t.Errorf("got=[%s], want=[%s]", lvl, AssociateOrProgrammer)
	}

	questions[0].Level = ProgrammerAnalyst
	questions[1].Level = SrProgrammer
	questions[2].Level = ProgrammerAnalyst

	lvl = findLevel(&questions, SrProgrammer)

	if lvl != SrProgrammer {
		t.Errorf("got=[%s], want=[%s]", lvl, SrProgrammer)
	}
}

// not feeling very proud about this test but ... meh
func Test_gotoNextQuestion(t *testing.T) {
	topics := make(map[string][]Question)
	linuxQuestions := []Question{
		Question{ID: 1, Q: "lx1", Level: AssociateOrProgrammer, Result: NotAnsweredYet},
		Question{ID: 2, Q: "lx2", Level: AssociateOrProgrammer, Result: NotAnsweredYet},
		Question{ID: 3, Q: "lx3", Level: AssociateOrProgrammer, Result: OK},
		Question{ID: 4, Q: "lx4", Level: SrProgrammer, Result: OK},
	}

	javaQuestions := []Question{
		Question{ID: 1, Q: "j1", Level: AssociateOrProgrammer, Result: NotAnsweredYet},
		Question{ID: 2, Q: "j2", Level: SrProgrammer, Result: Wrong},
	}

	randomQuestions := []Question{}

	topics["linux"] = linuxQuestions
	topics["java"] = javaQuestions
	topics["random"] = randomQuestions

	config := Config{}
	config.levels = [3]Level{
		AssociateOrProgrammer, ProgrammerAnalyst, SrProgrammer,
	}
	config.interview = Interview{Interviewee: "Hello", Date: time.Now(), Topics: topics}
	config.selectedTopic = "linux"
	config.hasStarted = true
	config.questionIndex = 0
	config.individualLevelIndexes = []int{0, 0, 0}

	gotoNextQuestion(&config)
	gotoNextQuestion(&config)

	if config.individualLevelIndexes[int(AssociateOrProgrammer)-1] != 2 {
		t.Errorf("got=[%d], want=[2]", config.individualLevelIndexes[int(AssociateOrProgrammer)])
	}

	config.individualLevelIndexes[0] = 0
	config.ignoreLevelChecking = true
	gotoNextQuestion(&config)
	if config.questionIndex != 1 {
		t.Errorf("questionIndex should be 1")
	}
}

func Test_gotoPreviousQuestion(t *testing.T) {
	topics := make(map[string][]Question)
	linuxQuestions := []Question{
		Question{ID: 1, Q: "lx1", Level: AssociateOrProgrammer, Result: NotAnsweredYet},
		Question{ID: 2, Q: "lx2", Level: AssociateOrProgrammer, Result: NotAnsweredYet},
		Question{ID: 3, Q: "lx3", Level: SrProgrammer, Result: OK},
		Question{ID: 4, Q: "lx4", Level: SrProgrammer, Result: OK},
	}

	javaQuestions := []Question{
		Question{ID: 1, Q: "j1", Level: AssociateOrProgrammer, Result: NotAnsweredYet},
		Question{ID: 2, Q: "j2", Level: SrProgrammer, Result: Wrong},
	}

	randomQuestions := []Question{}

	topics["linux"] = linuxQuestions
	topics["java"] = javaQuestions
	topics["random"] = randomQuestions

	config := Config{}
	config.levels = [3]Level{
		AssociateOrProgrammer, ProgrammerAnalyst, SrProgrammer,
	}
	config.interview = Interview{Interviewee: "Hello", Date: time.Now(), Topics: topics}
	config.selectedTopic = ""
	config.hasStarted = true
	config.questionIndex = 0
	config.individualLevelIndexes = []int{0, 0, 0}
	config.ignoreLevelChecking = true
	config.questionIndex = 1

	gotoPreviousQuestion(&config)
	if config.questionIndex != 1 {
		t.Errorf("index should have been decreased.")
	}
	config.selectedTopic = "java"
	gotoPreviousQuestion(&config)
	if config.questionIndex != 0 {
		t.Errorf("index should be 0")
	}

	fmt.Println("************************************")
	config.ignoreLevelChecking = false
	config.questionIndex = 1
	config.levelIndex = 1
	config.individualLevelIndexes = []int{1, 2, 3}

	gotoPreviousQuestion(&config)
	if config.individualLevelIndexes[1] != 1 {
		t.Errorf("expected 1, got=[%d] instead", config.individualLevelIndexes[1])
	}

}
