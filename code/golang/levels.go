package main

import (
	"fmt"

	"github.com/muesli/termenv"
)

func increaseLevel() {
	if (levelIndex + 1) < len(levels) {
		levelIndex++
		printWithColorln(fmt.Sprintf("Level is now: %s", levels[levelIndex]), yellow)
	} else {
		printWithColorln(fmt.Sprintf("Level cannot increased, currently at: %s", levels[levelIndex]), red)
	}
}

func decreaseLevel() {
	if (levelIndex - 1) >= 0 {
		levelIndex--
		printWithColorln(fmt.Sprintf("Level is now: %s", levels[levelIndex]), yellow)
	} else {
		printWithColorln(fmt.Sprintf("Level cannot be decreased, currently at: %s", levels[levelIndex]), red)
	}
}

func ignoreLevel() {
	ignoreLevelChecking = !ignoreLevelChecking
	if ignoreLevelChecking {
		printWithColorln("Ignoring level", cyan)
	} else {
		printWithColorln("Using level", cyan)
	}
}

func findLevel(questions *[]Question, levels ...Level) Level {
	foundLevel := AssociateOrProgrammer
	found := false
	for _, lvl := range levels {
		if found {
			break
		}
		for _, q := range *questions {
			if q.Level == lvl {
				found = true
				foundLevel = q.Level
				break
			}
		}
	}
	return foundLevel
}

func gotoNextQuestion() {
	if len(selectedTopic) == 0 {
		fmt.Println("Load a topic first.")
		return
	}

	if !hasStarted {
		fmt.Println("run the start() command first.")
	}

	if ignoreLevelChecking {
		if (questionIndex + 1) < len(interview.Topics[selectedTopic]) {
			questionIndex++
		} else {
			fmt.Println(termenv.String("No questions left ... ").Foreground(colorProfile.Color(yellow)))
		}
	} else {
		// TODO: put this in a method once we know it really works ...
		currentLevel := levels[levelIndex]
		currentLevelQuestions := getQuestionsFromLevel(currentLevel, selectedTopic, &interview.Topics)
		index := individualLevelIndexes[int(currentLevel)-1]
		if (index + 1) < len(currentLevelQuestions) {
			index++
			individualLevelIndexes[int(currentLevel)-1] = index
		} else {
			printWithColorln("That was the last question", yellow)
		}
	}
}

func gotoPreviousQuestion() {
	if (questionIndex - 1) >= 0 {
		questionIndex--
	}
}

func getQuestionsFromLevel(lvl Level, topic string, topics *map[string][]Question) []Question {
	questions := make([]Question, 0)
	for _, q := range (*topics)[topic] {
		if q.Level == lvl {
			questions = append(questions, q)
		}
	}
	return questions
}
