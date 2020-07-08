package interview

import (
	"fmt"

	"github.com/muesli/termenv"
)

func increaseLevel(config *Config) {
	if (config.levelIndex + 1) < len(config.levels) {
		config.levelIndex++
		printWithColorln(fmt.Sprintf("Level is now: %s", config.levels[config.levelIndex]), yellow, config)
	} else {
		printWithColorln(fmt.Sprintf("Level cannot increased, currently at: %s", config.levels[config.levelIndex]), red, config)
	}
}

func decreaseLevel(config *Config) {
	if (config.levelIndex - 1) >= 0 {
		config.levelIndex--
		printWithColorln(fmt.Sprintf("Level is now: %s", config.levels[config.levelIndex]), yellow, config)
	} else {
		printWithColorln(fmt.Sprintf("Level cannot be decreased, currently at: %s", config.levels[config.levelIndex]), red, config)
	}
}

func toggleLevelChecking(config *Config) {
	config.ignoreLevelChecking = !config.ignoreLevelChecking
	if config.ignoreLevelChecking {
		printWithColorln("Ignoring level", cyan, config)
	} else {
		printWithColorln("Using level", cyan, config)
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

func gotoNextQuestion(config *Config) {
	if len(config.selectedTopic) == 0 {
		fmt.Println("Load a topic first.")
		return
	}

	if !config.hasStarted {
		fmt.Println("run the start() command first.")
		return
	}

	if config.ignoreLevelChecking {
		if (config.questionIndex + 1) < len(config.interview.Topics[config.selectedTopic]) {
			config.questionIndex++
		} else {
			fmt.Println(termenv.String("No questions left ... ").Foreground(config.colorProfile.Color(yellow)))
		}
	} else {
		currentLevel := config.levels[config.levelIndex]
		currentLevelQuestions := getQuestionsFromLevel(currentLevel, config)
		index := config.individualLevelIndexes[int(currentLevel)-1]
		if (index + 1) < len(currentLevelQuestions) {
			index++
			config.individualLevelIndexes[int(currentLevel)-1] = index
		} else {
			printWithColorln("That was the last question", yellow, config)
		}
	}
}

func gotoPreviousQuestion(config *Config) {
	if len(config.selectedTopic) == 0 {
		fmt.Println("Load a topic first.")
		return
	}

	if config.ignoreLevelChecking {
		if (config.questionIndex - 1) >= 0 {
			config.questionIndex--
		}
	} else {
		currentLevel := config.levels[config.levelIndex]
		index := config.individualLevelIndexes[int(currentLevel)-1]
		if (index - 1) >= 0 {
			index--
			config.individualLevelIndexes[int(currentLevel)-1] = index
		} else {
			printWithColorln("That was the last question", yellow, config)
		}
	}
}

func getQuestionsFromLevel(lvl Level, config *Config) []Question {
	questions := make([]Question, 0)
	for _, q := range (config.interview.Topics)[config.selectedTopic] {
		if q.Level == lvl {
			questions = append(questions, q)
		}
	}
	return questions
}
