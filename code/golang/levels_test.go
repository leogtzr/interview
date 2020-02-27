package main

import (
	"testing"

	"github.com/muesli/termenv"
)

func Test_getQuestionsFromLevel(t *testing.T) {
	type test struct {
		lvl       Level
		topic     string
		questions []Question
	}

	topics := make(map[string][]Question)
	linuxQuestions := []Question{
		Question{ID: 1, Q: "lx1", Level: AssociateOrProgrammer, Answer: NotAnsweredYet},
		Question{ID: 2, Q: "lx2", Level: AssociateOrProgrammer, Answer: NotAnsweredYet},
		Question{ID: 3, Q: "lx3", Level: ProgrammerAnalyst, Answer: OK},
		Question{ID: 4, Q: "lx4", Level: SrProgrammer, Answer: OK},
	}

	javaQuestions := []Question{
		Question{ID: 1, Q: "j1", Level: AssociateOrProgrammer, Answer: NotAnsweredYet},
		Question{ID: 2, Q: "j2", Level: SrProgrammer, Answer: Wrong},
	}

	randomQuestions := []Question{}

	topics["linux"] = linuxQuestions
	topics["java"] = javaQuestions
	topics["random"] = randomQuestions

	tests := []test{
		{
			lvl:   AssociateOrProgrammer,
			topic: "linux",
			questions: []Question{
				Question{ID: 1, Q: "lx1", Level: AssociateOrProgrammer, Answer: NotAnsweredYet},
				Question{ID: 2, Q: "lx2", Level: AssociateOrProgrammer, Answer: NotAnsweredYet},
			},
		},

		{
			lvl:   SrProgrammer,
			topic: "java",
			questions: []Question{
				Question{ID: 2, Q: "j2", Level: SrProgrammer, Answer: Wrong},
			},
		},

		{
			lvl:       SrProgrammer,
			topic:     "random",
			questions: []Question{},
		},
	}

	for _, tt := range tests {
		got := getQuestionsFromLevel(tt.lvl, tt.topic, &topics)
		if !Equal(got, tt.questions) {
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

func Test_gotoNextQuestion(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotoNextQuestion(tt.args.config)
		})
	}
}
