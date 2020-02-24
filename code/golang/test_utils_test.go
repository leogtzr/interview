package main

import (
	"testing"
)

func TestEqual(t *testing.T) {
	type test struct {
		a      []Question
		b      []Question
		result bool
	}

	tests := []test{
		{
			a: []Question{
				Question{ID: 1, Answer: NotAnsweredYet, Q: "A", Level: SrProgrammer},
			},
			b: []Question{
				Question{ID: 1, Answer: NotAnsweredYet, Q: "A", Level: SrProgrammer},
			}, result: true,
		},
		{
			a: []Question{
				Question{ID: 1, Answer: NotAnsweredYet, Q: "A", Level: SrProgrammer},
			},
			b: []Question{
				Question{ID: 2, Answer: NotAnsweredYet, Q: "A", Level: SrProgrammer},
			}, result: false,
		},
		{
			a: []Question{},
			b: []Question{}, result: true,
		},
		{
			a: []Question{
				Question{ID: 1},
			},
			b: []Question{}, result: false,
		},
	}

	for _, tt := range tests {
		got := Equal(tt.a, tt.b)
		if got != tt.result {
			t.Errorf("[%s] and [%s] should be equal", tt.a, tt.b)
		}
	}
}

func TestEqualTopics(t *testing.T) {
	type test struct {
		a      []string
		b      []string
		result bool
	}

	tests := []test{
		{a: []string{"java", "linux"}, b: []string{"linux", "java"}, result: false},
		{a: []string{"java", "linux"}, b: []string{"java", "linux"}, result: true},
		{a: []string{"java"}, b: []string{"java", "linux"}, result: false},
	}

	for _, tt := range tests {
		got := EqualTopics(tt.a, tt.b)
		if got != tt.result {
			t.Errorf("got=[%t], should be [%t]", got, tt.result)
		}
	}
}
