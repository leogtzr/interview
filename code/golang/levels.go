package main

import (
	"fmt"

	"github.com/muesli/termenv"
)

func increaseLevel(lvlIndex *int, lvls [3]Level) {
	if (*lvlIndex + 1) < len(lvls) {
		*lvlIndex++
		printWithColorln(fmt.Sprintf("Level is now: %s", lvls[*lvlIndex]), yellow)
	} else {
		printWithColorln(fmt.Sprintf("Level cannot increased, currently at: %s", lvls[*lvlIndex]), red)
	}
}

func decreaseLevel(lvlIndex *int, lvls [3]Level) {
	if (*lvlIndex - 1) >= 0 {
		*lvlIndex--
		printWithColorln(fmt.Sprintf("Level is now: %s", lvls[*lvlIndex]), yellow)
	} else {
		printWithColorln(fmt.Sprintf("Level cannot be decreased, currently at: %s", lvls[*lvlIndex]), red)
	}
}

func toggleLevelChecking(lvlCheck *bool) {
	*lvlCheck = !(*lvlCheck)
	if *lvlCheck {
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
	if len(selectedTopic) == 0 {
		fmt.Println("Load a topic first.")
		return
	}

	if ignoreLevelChecking {
		if (questionIndex - 1) >= 0 {
			questionIndex--
		}
	} else {
		currentLevel := levels[levelIndex]
		index := individualLevelIndexes[int(currentLevel)-1]
		if (index - 1) >= 0 {
			index--
			individualLevelIndexes[int(currentLevel)-1] = index
		} else {
			printWithColorln("That was the last question", yellow)
		}
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
