.DEFAULT_GOAL := install

BIN_FILE=interview

install:
	stringer -type=Command
	stringer -type=Answer
	stringer -type=Level
	go build -o "${BIN_FILE}"

clean:
	go clean
	rm --force interview
	rm --force cp.out

test:
	go test

check:
	go test

cover:
	go test -coverprofile cp.out
	go tool cover -html=cp.out

testdb:
	go test -v dao_test.go
	