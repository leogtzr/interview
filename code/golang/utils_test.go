package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/muesli/termenv"
)

func Test_sanitizeUserInput(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		{"	alv", "alv"},
	}

	for _, tc := range tests {
		if got := sanitizeUserInput(tc.input); got != tc.want {
			t.Errorf("got=[%s], want=[%s]", got, tc.want)
		}
	}
}

func Test_userInputToCmd(t *testing.T) {
	type test struct {
		input string
		want  Command
	}

	tests := []test{
		{input: ":_", want: noCmd},
		{input: ":q", want: exitCmd},
		{input: "exit", want: exitCmd},
		{input: "use golang", want: useCmd},
		{input: "cls", want: clearScreenCmd},
		{input: "pwd", want: pwdCmd},
		{input: "start", want: startCmd},
		{input: "print", want: printCmd},
		{input: ">", want: nextQuestionCmd},
		{input: "<", want: previousQuestionCmd},
		{input: "view", want: viewCmd},
		{input: "ok", want: rightAnswerCmd},
		{input: "no", want: wrongAnswerCmd},
		{input: "meh", want: mehAnswerCmd},
		{input: "finish", want: finishCmd},
		{input: "topics", want: topicsCmd},
		{input: "help", want: helpCmd},
		{input: "+", want: increaseLevelCmd},
		{input: "-", want: decreaseLevelCmd},
		{input: "=", want: ignoreLevelCmd},
		{input: "lvl", want: showLevelCmd},
		{input: "stats", want: showStatsCmd},
		{input: "cmt", want: createCommentCmd},
		{input: "cq", want: createQuestionCmd},
		{input: "ap", want: setAssociateProgrammerLevelCmd},
		{input: "pa", want: setProgrammerAnalystLevelCmd},
		{input: "sr", want: setSRProgrammerLevelCmd},
		{input: "use", want: noCmd},
		{input: "c", want: countCmd},
		{input: "", want: noCmd},
	}

	for _, tc := range tests {
		if got, _ := userInputToCmd(tc.input); got != tc.want {
			t.Errorf("got=[%s], want=[%s]", got, tc.want)
		}
	}
}

func Test_topicExist(t *testing.T) {
	topics := []string{"linux", "sql", "java", "go", "c", "c++"}
	if !topicExist("linux", &topics) {
		t.Errorf("Should be in the list of topics ... ")
	}

	if !topicExist("java", &topics) {
		t.Errorf("Should be in the list of topics ... ")
	}

	if topicExist("ashkdh", &topics) {
		t.Errorf("Should NOT be in the list of topics ... ")
	}
}

func Test_shortIntervieweeName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Leonardo Gutierrez", "(Leonardo G...)"},
		{"Leonardo", "(Leonardo)"},
		{"Leo", "(Leo)"},
		{"", "(who?)"},
	}
	for _, tt := range tests {
		if got := shortIntervieweeName(tt.name, minNumberOfCharsInIntervieweeName); got != tt.want {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}

func TestQuestion_String(t *testing.T) {
	type test struct {
		q    Question
		want string
	}

	tests := []test{
		{q: Question{ID: 1, Q: "hola", Result: NotAnsweredYet, Level: ProgrammerAnalyst}, want: "Q1: hola [NotAnsweredYet] [ProgrammerAnalyst]"},
		{q: Question{ID: 1, Q: "hola", Result: NotAnsweredYet, Level: SrProgrammer}, want: "Q1: hola [NotAnsweredYet] [SrProgrammer]"},
		{q: Question{ID: 1, Q: "hola", Result: Neutral, Level: AssociateOrProgrammer}, want: "Q1: hola [Neutral] [AssociateOrProgrammer]"},
	}

	for _, tt := range tests {
		if tt.q.String() != tt.want {
			t.Errorf("got=[%s], want=[%s]", tt.q.String(), tt.want)
		}
	}
}

func Test_readIntervieweeName(t *testing.T) {
	type test struct {
		input      string
		output     string
		hasContent bool
	}

	tests := []test{
		{input: "Leonardo", output: "Leonardo", hasContent: true},
		{input: "", output: "", hasContent: false},
	}

	for _, tt := range tests {
		got, ok := readIntervieweeName(strings.NewReader(tt.input))
		if ok != tt.hasContent {
			t.Errorf("got=[%t], want=[%t]", ok, tt.hasContent)
		}
		if got != tt.output {
			t.Errorf("got=[%s], want=[%s]", got, tt.output)
		}
	}
}

func Test_ps1String(t *testing.T) {
	type args struct {
		ps1             string
		selectedTopic   string
		intervieweeName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "asd",
			args: args{ps1: "$ ", selectedTopic: "linux", intervieweeName: "leo"},
			want: "2f1b5b326d6c696e75781b5b306d20286c656f29202420",
		},
		{
			name: "asd2",
			args: args{ps1: "$ ", selectedTopic: "", intervieweeName: "leo"},
			want: "2420",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmt.Sprintf("%x", ps1String(tt.args.ps1, tt.args.selectedTopic, tt.args.intervieweeName)); got != tt.want {
				t.Errorf("ps1String() = [%v], want [%v]", got, tt.want)
			}
		})
	}
}

func Test_setLevel(t *testing.T) {

	colorProfile := termenv.ColorProfile()

	type test struct {
		level  Level
		config Config
		want   int
	}

	lvls := [3]Level{
		AssociateOrProgrammer, ProgrammerAnalyst, SrProgrammer,
	}

	tests := []test{
		{level: AssociateOrProgrammer, config: Config{levelIndex: -1, colorProfile: colorProfile, levels: lvls}, want: 0},
		{level: ProgrammerAnalyst, config: Config{levelIndex: -1, colorProfile: colorProfile, levels: lvls}, want: 1},
		{level: SrProgrammer, config: Config{levelIndex: -1, colorProfile: colorProfile, levels: lvls}, want: 2},
	}

	for _, tt := range tests {
		setLevel(tt.level, &tt.config)
		if tt.config.levelIndex != tt.want {
			t.Errorf("got=[%d], want=[%d]", tt.config.levelIndex, tt.want)
		}
	}

}

func Test_perc(t *testing.T) {
	type test struct {
		count, total int
		want         float64
	}

	tests := []test{
		{count: 30, total: 100, want: 30.0},
		{count: 10, total: 100, want: 10.0},
	}

	for _, tt := range tests {
		if got := perc(tt.count, tt.total); got != tt.want {
			t.Errorf("got=[%f], want=[%f]", got, tt.want)
		}
	}
}

func Test_countGeneral(t *testing.T) {
	type test struct {
		topics map[string][]Question
		want   map[Result]int
	}

	tests := []test{
		{
			topics: map[string][]Question{
				"java": []Question{
					Question{ID: 1, Q: "A", Result: NotAnsweredYet, Level: SrProgrammer},
					Question{ID: 2, Q: "A", Result: NotAnsweredYet, Level: SrProgrammer},
					Question{ID: 3, Q: "A", Result: Wrong, Level: SrProgrammer},
				},
			},
			want: map[Result]int{
				NotAnsweredYet: 2,
				Wrong:          1,
				Neutral:        0,
			},
		},
	}

	for _, tt := range tests {
		got := countGeneral(&tt.topics)
		if got[NotAnsweredYet] != tt.want[NotAnsweredYet] {
			t.Errorf("got=[%d], want=[%d]", got[NotAnsweredYet], tt.want[NotAnsweredYet])
		}
		if got[Wrong] != tt.want[Wrong] {
			t.Errorf("got=[%d], want=[%d]", got[Wrong], tt.want[Wrong])
		}
		if got[OK] != tt.want[OK] {
			t.Errorf("got=[%d], want=[%d]", got[OK], tt.want[OK])
		}
		if got[Neutral] != tt.want[Neutral] {
			t.Errorf("got=[%d], want=[%d]", got[Neutral], tt.want[Neutral])
		}
	}

}

func Test_extractTopicName(t *testing.T) {
	type test struct {
		options []string
		want    string
	}

	tests := []test{
		{options: []string{"a", "b", "c"}, want: "a"},
	}

	for _, tt := range tests {
		if got := extractTopicName(tt.options); got != tt.want {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}

func Test_resetStatus(t *testing.T) {
	config := Config{}
	topics := make(map[string][]Question)
	javaQuestions := []Question{
		Question{ID: 1, Q: "j1", Level: AssociateOrProgrammer, Result: NotAnsweredYet},
		Question{ID: 2, Q: "j2", Level: SrProgrammer, Result: Wrong},
	}
	topics["java"] = javaQuestions

	config.interview = Interview{Interviewee: "Hello", Date: time.Now(), Topics: topics}
	// config.usingInterviewFile = true
	config.hasStarted = true
	config.questionIndex = 23
	config.selectedTopic = "java"
	config.ps1 = "hello"

	resetStatus(&config)

	if len(config.interview.Topics) != 0 {
		t.Error("topics should be empty")
	}
	// if config.usingInterviewFile {
	// 	t.Error("flag should been changed")
	// }
	if config.hasStarted {
		t.Error("flag should been changed")
	}
	if config.questionIndex != 0 {
		t.Error("questionIndex should be 0 after reset.")
	}
	if config.selectedTopic != "" {
		t.Error("topic should be empty after reset")
	}
	if config.ps1 != "$ " {
		t.Error("ps1 should be '$ ' after reset")
	}
}

func Test_levelQuestionCounts(t *testing.T) {
	type test struct {
		qs   []Question
		want map[Level]int
	}

	tests := []test{
		{qs: []Question{
			Question{Level: AssociateOrProgrammer},
			Question{Level: AssociateOrProgrammer},
			Question{Level: AssociateOrProgrammer},
			Question{Level: SrProgrammer},
		}, want: map[Level]int{
			AssociateOrProgrammer: 3,
			SrProgrammer:          1,
			ProgrammerAnalyst:     0,
		}},

		{qs: []Question{
			Question{Level: AssociateOrProgrammer},
			Question{Level: AssociateOrProgrammer},
			Question{Level: AssociateOrProgrammer},
		}, want: map[Level]int{
			AssociateOrProgrammer: 3,
			SrProgrammer:          0,
			ProgrammerAnalyst:     0,
		}},

		{qs: []Question{}, want: map[Level]int{}},
	}

	lvls := []Level{
		AssociateOrProgrammer, ProgrammerAnalyst, SrProgrammer,
	}

	for _, tt := range tests {
		got := levelQuestionCounts(&tt.qs)
		for _, lvl := range lvls {
			if got[lvl] != tt.want[lvl] {
				t.Errorf("got=[%d], want=[%d] for %s", got[lvl], tt.want[lvl], lvl)
			}
		}
	}
}

func Test_markQuestionAs(t *testing.T) {
	type test struct {
		id  int
		ans Result
		qs  []Question
	}

	tests := []test{
		{id: 1, ans: Wrong, qs: []Question{
			Question{ID: 0, Result: NotAnsweredYet},
			Question{ID: 1, Result: NotAnsweredYet},
			Question{ID: 2, Result: NotAnsweredYet},
			Question{ID: 3, Result: NotAnsweredYet},
			Question{ID: 4, Result: NotAnsweredYet},
		}},
	}

	for _, tt := range tests {
		markQuestionAs(tt.id, tt.ans, &tt.qs)
		if tt.qs[tt.id-1].Result != tt.ans {
			t.Errorf("want=[%s], got=[%s]", tt.ans, tt.qs[tt.id].Answer)
		}
	}
}

func TestNewConfig(t *testing.T) {
	const expectedNumberOfIndividualLevelIndexes = 3
	const expectedSelectedTopic = ""
	const expectedPs1 = "$ "
	const expectedNumberOfInitialTopics = 0
	const expectedQuestionLevel = AssociateOrProgrammer
	const expectedInitialLevelIndex = 0
	const expectedIgnoringLevelCheck = false
	expectedIndividualLevelIndexes := []int{0, 0, 0}
	const expectedUsingInterviewFile = false
	expectedLevels := [3]Level{
		AssociateOrProgrammer, ProgrammerAnalyst, SrProgrammer,
	}

	config := NewConfig()
	if len(config.individualLevelIndexes) != expectedNumberOfIndividualLevelIndexes {
		t.Errorf("got=[%d], want=[%d]", len(config.individualLevelIndexes), expectedNumberOfIndividualLevelIndexes)
	}
	if config.selectedTopic != expectedSelectedTopic {
		t.Errorf("got=[%s], want=[%s]", config.selectedTopic, expectedSelectedTopic)
	}
	if config.ps1 != expectedPs1 {
		t.Errorf("got=[%s], want=[%s]", config.ps1, expectedPs1)
	}
	if len(config.interview.Topics) != expectedNumberOfInitialTopics {
		t.Errorf("got=[%d], want=[%d]", len(config.interview.Topics), expectedNumberOfInitialTopics)
	}
	if config.topicQuestionsLevel != expectedQuestionLevel {
		t.Errorf("got=[%s], want=[%s]", config.topicQuestionsLevel, expectedQuestionLevel)
	}
	if config.ignoreLevelChecking != expectedIgnoringLevelCheck {
		t.Errorf("got=[%t], want=[%t]", config.ignoreLevelChecking, expectedIgnoringLevelCheck)
	}
	if config.levelIndex != expectedInitialLevelIndex {
		t.Errorf("got=[%d], want=[%d]", config.levelIndex, expectedInitialLevelIndex)
	}
	if !EqualNumbers(expectedIndividualLevelIndexes, config.individualLevelIndexes) {
		got := strings.Trim(strings.Replace(fmt.Sprint(config.individualLevelIndexes), " ", ",", -1), "[]")
		want := strings.Trim(strings.Replace(fmt.Sprint(expectedIndividualLevelIndexes), " ", ",", -1), "[]")
		t.Errorf("got=[%s], want=[%s]", got, want)
	}

	if config.levels != expectedLevels {
		t.Errorf("got=[%s], want=[%s]", config.levels, expectedLevels)
	}
}
