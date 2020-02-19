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
