package main

import (
	"fmt"
	"os"
	"testing"
)

var x int = 23

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	// Do something here.
	fmt.Printf("\033[1;36m%s\033[0m", "> Setup completed\n")
}

func teardown() {
	// Do something here.

	fmt.Printf("\033[1;36m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}

func TestCreate_One(t *testing.T) {
	x++
	t.Logf("1) The value is: %d\n", x)
}

func TestCreate_Two(t *testing.T) {
	x++
	t.Logf("2) The value is: %d\n", x)
}
