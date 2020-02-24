package main

import (
	"testing"
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
	type test struct {
		lvlIndex           int
		lvls               [3]Level
		incraseTimes, want int
	}

	levels := [3]Level{
		AssociateOrProgrammer, ProgrammerAnalyst, SrProgrammer,
	}

	tests := []test{
		{lvlIndex: 0, lvls: levels, incraseTimes: 1, want: 1},
		{lvlIndex: 0, lvls: levels, incraseTimes: 10, want: len(levels) - 1},
	}

	for _, tt := range tests {
		for i := 0; i < tt.incraseTimes; i++ {
			increaseLevel(&tt.lvlIndex, tt.lvls)
		}
		if tt.lvlIndex != tt.want {
			t.Errorf("want=[%d] after increasing %d times, current index is: %d", tt.want, tt.incraseTimes, tt.lvlIndex)
		}
	}
}

func Test_decreaseLevel(t *testing.T) {
	type test struct {
		lvlIndex            int
		lvls                [3]Level
		decreaseTimes, want int
	}

	interviewLevels := [3]Level{
		AssociateOrProgrammer, ProgrammerAnalyst, SrProgrammer,
	}

	tests := []test{
		{lvlIndex: len(interviewLevels), lvls: interviewLevels, decreaseTimes: 2, want: 1},
		{lvlIndex: len(interviewLevels), lvls: interviewLevels, decreaseTimes: 1, want: len(interviewLevels) - 1},
		{lvlIndex: len(interviewLevels), lvls: interviewLevels, decreaseTimes: 2, want: len(interviewLevels) - 2},
		{lvlIndex: 0, lvls: levels, decreaseTimes: 10, want: 0},
	}

	for _, tt := range tests {
		for i := 0; i < tt.decreaseTimes; i++ {
			decreaseLevel(&tt.lvlIndex, tt.lvls)
		}
		if tt.lvlIndex != tt.want {
			t.Errorf("want=[%d] after decreasing %d times, current index is: %d", tt.want, tt.decreaseTimes, tt.lvlIndex)
		}
	}
}

func Test_toggleLevelChecking(t *testing.T) {
	type test struct {
		lvlCheck bool
		want     bool
	}

	tests := []test{
		{lvlCheck: false, want: true},
		{lvlCheck: true, want: false},
	}

	for _, tt := range tests {
		toggleLevelChecking(&tt.lvlCheck)
		if tt.lvlCheck != tt.want {
			t.Errorf("the check should be %t after change. Currently is: %t", tt.want, tt.lvlCheck)
		}
	}
}
