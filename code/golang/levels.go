package main

import "fmt"

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
