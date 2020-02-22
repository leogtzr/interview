package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"
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
		got := sanitizeUserInput(tc.input)
		if got != tc.want {
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
		{input: "load", want: loadCmd},
		{input: "topics", want: topicsCmd},
		{input: "help", want: helpCmd},
		{input: "exf", want: exitInterviewFileCmd},
		{input: "+", want: increaseLevelCmd},
		{input: "-", want: decreaseLevelCmd},
		{input: "=", want: ignoreLevelCmd},
		{input: "lvl", want: showLevelCmd},
		{input: "stats", want: showStatsCmd},
		{input: "ap", want: setAssociateProgrammerLevelCmd},
		{input: "pa", want: setProgrammerAnalystLevelCmd},
		{input: "sr", want: setSRProgrammerLevelCmd},
	}

	for _, tc := range tests {
		got, _ := userInputToCmd(tc.input)
		if got != tc.want {
			t.Errorf("got=[%s], want=[%s]", got, tc.want)
		}
	}
}

func Test_questionHasValidFormat(t *testing.T) {
	rgx := regexp.MustCompile("^\\d+@.+@(\\d+)?$")
	type test struct {
		input string
		match bool
	}

	tests := []test{
		{input: "1@Cómo puedes sortear un archivo?@2", match: true},
		{input: "1@Cómo puedes sortear un archivo?@s", match: false},
		{input: "@Cómo puedes sortear un archivo?@s", match: false},
		{input: "1@x@2", match: true},
		{input: "1@@2", match: false},
	}
	for _, tc := range tests {
		match := isQuestionFormatValid(tc.input, rgx)
		if match != tc.match {
			t.Errorf("got=[%t], want=[%t]", match, tc.match)
		}
	}
}

func Test_toQuestion(t *testing.T) {
	type test struct {
		input    string
		question Question
	}

	tests := []test{
		{
			input:    "1@Cómo puedes sortear un archivo?@2",
			question: Question{ID: 1, Q: "Cómo puedes sortear un archivo?", Answer: NotAnsweredYet, Level: ProgrammerAnalyst},
		},
		{
			input:    "2@Cómo puedes obtener las ultimas 3 líneas de un archivo?@3",
			question: Question{ID: 2, Q: "Cómo puedes obtener las ultimas 3 líneas de un archivo?", Answer: NotAnsweredYet, Level: SrProgrammer},
		},
	}

	for _, tc := range tests {
		got := toQuestion(tc.input)
		if got != tc.question {
			t.Errorf("got=[%s], want=[%s]", got, tc.question)
		}
	}
}

func Test_retrieveTopics(t *testing.T) {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	randDirName := stringWithCharset(5, charset, seededRand)
	tmpPath := filepath.Join("/tmp", randDirName, "topics", "linux")

	err := os.MkdirAll(tmpPath, os.ModePerm)
	if err != nil {
		t.Errorf("Error creating directory (%s)", err)
	}

	topics := retrieveTopicsFromFileSystem(filepath.Join("/tmp", randDirName))
	if topics == nil || len(topics) != 1 {
		t.Errorf("Not able to get topics ... ")
	}

	topicExpectedName := "linux"
	if topics[0] != topicExpectedName {
		t.Errorf("got=[%s], want=[%s]", topics[0], topicExpectedName)
	}

	err = os.RemoveAll(tmpPath)
	if err != nil {
		t.Errorf("unexpedted error: [%s]", err)
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
		got := shortIntervieweeName(tt.name, minNumberOfCharsInIntervieweeName)
		if got != tt.want {
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
		{q: Question{ID: 1, Q: "hola", Answer: NotAnsweredYet, Level: ProgrammerAnalyst}, want: "Q1: hola [NotAnsweredYet] [ProgrammerAnalyst]"},
		{q: Question{ID: 1, Q: "hola", Answer: NotAnsweredYet, Level: SrProgrammer}, want: "Q1: hola [NotAnsweredYet] [SrProgrammer]"},
		{q: Question{ID: 1, Q: "hola", Answer: Neutral, Level: AssociateOrProgrammer}, want: "Q1: hola [Neutral] [AssociateOrProgrammer]"},
	}

	for _, tt := range tests {
		if tt.q.String() != tt.want {
			t.Errorf("got=[%s], want=[%s]", tt.q.String(), tt.want)
		}
	}
}

func Test_readIntervieweeName(t *testing.T) {
	name := "Leonardo"
	reader := strings.NewReader(name)

	got, ok := readIntervieweeName(reader)
	if !ok {
		t.Errorf("got=[%s], want=[%s]", got, name)
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

func Test_extractQuestionInfo(t *testing.T) {
	type test struct {
		record   string
		topic    string
		question Question
	}

	tests := []test{
		{record: "linux@1@Cómo puedes sortear un archivo?@2@1",
			topic: "linux", question: Question{ID: 1, Q: "Cómo puedes sortear un archivo?", Answer: NotAnsweredYet}},
	}

	for _, tt := range tests {
		gotTopic, gotQuestion := extractQuestionInfo(tt.record)
		if gotTopic != tt.topic {
			t.Errorf("got=[%s], want=[%s]", gotTopic, tt.topic)
		}
		if gotQuestion != tt.question {
			t.Errorf("got=[%s], want=[%s]", gotQuestion, tt.question)
		}
	}

}

func Test_extractDateFromInterviewHeaderRecord(t *testing.T) {
	type test struct {
		header     string
		want       time.Time
		shouldFail bool
	}

	want, _ := time.Parse(interviewFormatLayout, "2020-02-11 22:32:28")
	tests := []test{
		{header: "Leo Gtz@2020-02-11 22:32:28", want: want, shouldFail: false},
		{header: "abc", want: want, shouldFail: true},
	}

	for _, tt := range tests {
		got, err := extractDateFromInterviewHeaderRecord(tt.header)
		if err != nil && !tt.shouldFail {
			t.Errorf("It should have failed for -> [%s]", tt.header)
		}

		if got != tt.want && !tt.shouldFail {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}

func Test_extractNameFromInterviewHeaderRecord(t *testing.T) {
	type test struct {
		header     string
		want       string
		shouldFail bool
	}
	tests := []test{
		{header: "Leo Gtz@2020-02-11 22:32:28", want: "Leo Gtz", shouldFail: false},
		{header: "Leo Gtz", want: "", shouldFail: true},
	}
	for _, tt := range tests {
		got, err := extractNameFromInterviewHeaderRecord(tt.header)
		if !tt.shouldFail && err != nil {
			t.Errorf("it should have failed with: [%s]", tt.header)
		}

		if tt.shouldFail && err == nil {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}

func Test_setLevel(t *testing.T) {
	type test struct {
		level Level
		want  int
	}

	lvls := [3]Level{
		AssociateOrProgrammer, ProgrammerAnalyst, SrProgrammer,
	}

	tests := []test{
		{level: AssociateOrProgrammer, want: 0},
		{level: ProgrammerAnalyst, want: 1},
		{level: SrProgrammer, want: 2},
	}

	for _, tt := range tests {
		idx := -1
		setLevel(tt.level, &idx, lvls)
		if idx != tt.want {
			t.Errorf("got=[%d], want=[%d]", idx, tt.want)
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
		got := perc(tt.count, tt.total)
		if got != tt.want {
			t.Errorf("got=[%f], want=[%f]", got, tt.want)
		}
	}
}

func Test_countGeneral(t *testing.T) {
	type test struct {
		topics map[string][]Question
		want   map[Answer]int
	}

	tests := []test{
		{
			topics: map[string][]Question{
				"java": []Question{
					Question{ID: 1, Q: "A", Answer: NotAnsweredYet, Level: SrProgrammer},
					Question{ID: 2, Q: "A", Answer: NotAnsweredYet, Level: SrProgrammer},
					Question{ID: 3, Q: "A", Answer: Wrong, Level: SrProgrammer},
				},
			},
			want: map[Answer]int{
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

func Test_setAnswerAsNeutral(t *testing.T) {
	qs := []Question{
		Question{ID: 1, Answer: OK},
		Question{ID: 2, Answer: OK},
		Question{ID: 3, Answer: Wrong},
		Question{ID: 4, Answer: NotAnsweredYet},
	}

	for idx := range qs {
		setAnswerAsNeutral(&qs, idx)
	}

	for _, q := range qs {
		if q.Answer != Neutral {
			t.Errorf("%s should have been marked as Neutral.", q)
		}
	}
}

func Test_setAnswerAsOK(t *testing.T) {
	qs := []Question{
		Question{ID: 1, Answer: OK},
		Question{ID: 2, Answer: OK},
		Question{ID: 3, Answer: Wrong},
		Question{ID: 4, Answer: NotAnsweredYet},
	}

	for idx := range qs {
		setAnswerAsOK(&qs, idx)
	}

	for _, q := range qs {
		if q.Answer != OK {
			t.Errorf("%s should have been marked as OK.", q)
		}
	}
}

func Test_setAnswerAsWrong(t *testing.T) {
	qs := []Question{
		Question{ID: 1, Answer: OK},
		Question{ID: 2, Answer: OK},
		Question{ID: 3, Answer: Wrong},
		Question{ID: 4, Answer: NotAnsweredYet},
	}

	for idx := range qs {
		setAnswerAsWrong(&qs, idx)
	}

	for _, q := range qs {
		if q.Answer != Wrong {
			t.Errorf("%s should have been marked as Wrong.", q)
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
		got := extractTopicName(tt.options)
		if got != tt.want {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}

func Test_retrieveTopicsFromInterview(t *testing.T) {
	type test struct {
		topics map[string][]Question
		want   []string
	}
	tests := []test{
		{
			topics: map[string][]Question{
				"java":  []Question{},
				"linux": []Question{},
				"bash":  []Question{},
			},
			want: []string{"bash", "java", "linux"},
		},
	}

	for _, tt := range tests {
		got := retrieveTopicsFromInterview(&tt.topics)
		sort.Strings(got)
		if !EqualTopics(got, tt.want) {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}

}
