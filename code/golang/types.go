package main

// Command ...
type Command int

const (
	exitCmd   Command = iota
	topicsCmd Command = iota
	helpCmd   Command = iota
	useCmd    Command = iota
	noCmd     Command = iota
)
