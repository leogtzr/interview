package main

// Command ...
type Command int

const (
	exitCmd   Command = iota
	topicsCmd Command = iota
	noCmd     Command = iota
)
