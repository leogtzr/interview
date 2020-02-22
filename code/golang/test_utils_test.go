package main

import "testing"

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
