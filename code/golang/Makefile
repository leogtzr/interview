.DEFAULT_GOAL := install

# INSTALL_SCRIPT=./install.sh
BIN_FILE=interview

install:
	stringer -type=Command
	stringer -type=Answer
	stringer -type=Level
	go test
	go build -o "${BIN_FILE}"

clean:
	go clean

test:
	go test

check:
	go test

cover:
	go test -coverprofile cp.out
	go tool cover -html=cp.out
	