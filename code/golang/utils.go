package main

import "strings"

func sanitizeUserInput(input string) string {
	return strings.TrimSpace(input)
}

// Transforms user's input to a Command
func userInputToCmd(input string) Command {
	input = sanitizeUserInput(input)
	input = strings.ToLower(input)
	switch input {
	case "exit", "quit", ":q", "/q":
		return exitCmd
	case "topics", "tps", "t", "/t":
		return topicsCmd
	}
	return noCmd
}
