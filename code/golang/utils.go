package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/muesli/termenv"
)

func sanitizeUserInput(input string) string {
	return strings.TrimSpace(input)
}

// Transforms user's input to a Command
func userInputToCmd(input string) (Command, []string) {
	fullCommand := words(input)
	input = fullCommand[0]
	input = sanitizeUserInput(input)
	input = strings.ToLower(input)
	switch input {
	case "exit", "quit", ":q", "/q", "q":
		return exitCmd, []string{}
	case "topics", "tps", "t", "/t", ":t":
		return topicsCmd, []string{}
	case "help", ":h", "/h", "--h", "-h":
		return helpCmd, []string{}
	case "use", "u", "/u", ":u", "-u", "--u", "set":
		if len(fullCommand) <= 1 {
			return noCmd, []string{}
		}
		return useCmd, fullCommand[1:]
	case "cls", "clear":
		return clearScreenCommand, []string{}
	case "pwd":
		return pwdCommand, []string{}
	case "start", "begin":
		return startCmd, []string{}
	case "p", "print", "print()", "p()":
		return printCmd, []string{}
	case "next", "nxt", ">":
		return nextQuestionCmd, []string{}
	case "previous", "prev", "<":
		return previousQuestionCmd, []string{}
	case "view", "v":
		return viewCmd, []string{}
	case "y", "right", "ok", "yes", "si":
		return rightAnswerCmd, []string{}
	case "n", "no", "mal", "wrong", "nop", "bad", "nel":
		return wrongAnswerCmd, []string{}
	case "hmm", "meh", "?":
		return mehAnswerCmd, []string{}
	case "finish", "done", "bye":
		return finishCmd, []string{}
	}
	return noCmd, []string{}
}

func dirExists(dirPath string) bool {
	if _, err := os.Stat(dirPath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

func exists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}

func retrieveTopics(interviewsDir string) []string {
	topicsDir := filepath.Join(interviewsDir, "topics")
	topicsInDir := []string{}

	if !dirExists(topicsDir) {
		log.Fatalf("'%s' does not exist", topicsDir)
	}

	err := filepath.Walk(topicsDir, func(path string, info os.FileInfo, err error) error {
		path = filepath.Base(path)
		if path == "topics" {
			return nil
		}
		topicsInDir = append(topicsInDir, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return topicsInDir
}

func listTopics(interviewsDir string) {
	topics := retrieveTopics(interviewsDir)
	for _, topic := range topics {
		fmt.Println(termenv.String(topic).Underline().Bold())
	}
}

// TODO: complete help message.
func printHelp() {
	usage := `
commands:


	`

	fmt.Println(usage)
}

func words(input string) []string {
	return strings.Fields(input)
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func topicExist(topic string, topics *[]string) bool {
	r := false

	for _, t := range *topics {
		if t == topic {
			r = true
			break
		}
	}

	return r
}

func toQuestion(question string) Question {
	questionFields := strings.Split(question, "@")
	id, _ := strconv.ParseInt(questionFields[0], 10, 64)
	q := questionFields[1]
	nextID, _ := strconv.ParseInt(questionFields[2], 10, 64)
	if nextID == 0 {
		nextID = -1
	}
	return Question{ID: int(id), Q: q, NextQuestionID: int(nextID), Answer: NotAnsweredYet}
}

func loadTopics(topic, interviewsDir string, questions *[]Question) {
	// Clear previous questions ...
	questionsPerTopic = nil

	questionFilePath := filepath.Join(interviewsDir, "topics", topic, "questions")

	file, err := os.Open(questionFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		questionText := scanner.Text()
		if isQuestionFormatValid(questionText, rgxQuestions) {
			question := toQuestion(questionText)
			questionsPerTopic = append(questionsPerTopic, question)
		}
	}

	fmt.Printf("Loaded -> '%d' questions.\n", len(questionsPerTopic))

}

func setTopic(options []string) {
	topics := retrieveTopics(interviewTopicsDir)
	topicName := options[0]
	topicName = strings.ToLower(topicName)

	if topicExist(topicName, &topics) &&
		exists(filepath.Join(interviewTopicsDir, "topics", topicName, "questions")) {
		selectedTopic = topicName
		loadTopics(selectedTopic, interviewTopicsDir, &questionsPerTopic)
	} else {
		fmt.Println(
			termenv.String(fmt.Sprintf("topic '%s' not found or the topic selected doesn't have questions.", topicName)).Foreground(colorProfile.Color(red)))
	}
}

func shortIntervieweeName(name string, min int) string {
	if len(name) == 0 {
		return ""
	}
	if len(name) < min {
		return fmt.Sprintf("(%s)", name)
	}
	return fmt.Sprintf("(%s...)", name[0:min])
}

func ps1String(ps1, selectedTopic string) string {
	if selectedTopic == "" {
		return "$ "
	}
	return fmt.Sprintf("/%s %s $ ", termenv.String(selectedTopic).Faint(), shortIntervieweeName(intervieweeName, minNumberOfCharsInIntervieweeName))
}

// func printWorkingDirectory() {
// 	fmt.Println(termenv.String(selectedTopic).Bold())
// }

func isQuestionFormatValid(question string, rgx *regexp.Regexp) bool {
	return rgx.MatchString(question)
}

func (q Question) String() string {
	if q.NextQuestionID == -1 {
		return fmt.Sprintf("Q%d: %s", q.ID, q.Q)
	}
	return fmt.Sprintf("Q%d: %s (next: %d)", q.ID, q.Q, q.NextQuestionID)
}

func printQuestion(questionIndex int) {
	if hasStarted && (len(questionsPerTopic) > 0) {
		fmt.Println(questionsPerTopic[questionIndex])
	}
}

func gotoNextQuestion() {
	if len(selectedTopic) == 0 {
		fmt.Println("Load a topic first.")
		return
	}

	if !hasStarted {
		fmt.Println("run the start() command first.")
	}

	if (questionIndex + 1) < len(questionsPerTopic) {
		questionIndex++
	} else {
		fmt.Println(termenv.String("No questions left ... ").Foreground(colorProfile.Color(yellow)))
	}
}

func gotoPreviousQuestion() {
	if (questionIndex - 1) >= 0 {
		questionIndex--
	}
}

func viewStats() {
	if len(questionsPerTopic) < 1 {
		return
	}
	for _, q := range questionsPerTopic {
		fmt.Printf("[%s] -> [%s]\n", q, q.Answer)
	}
}

func readIntervieweeName() (string, bool) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Interviewee name: ")
	text, err := reader.ReadString('\n')
	if err != nil || (len(strings.TrimSpace(text)) == 0) {
		return "", false
	}
	return strings.TrimSpace(text), true
}

func markAnswerAs(ans Answer) {
	questionsPerTopic[questionIndex].Answer = ans
	printWithColorln(fmt.Sprintf("Answer has saved as '%s'", ans), green)
}

func printWithColorln(msg, colorCode string) {
	fmt.Println(termenv.String(msg).Foreground(colorProfile.Color(colorCode)))
}

func markAnswerAsOK() {
	questionsPerTopic[questionIndex].Answer = OK
	printWithColorln(fmt.Sprintf("Answer has saved as '%s'", OK), green)
}

func markAnswerAsWrong() {
	questionsPerTopic[questionIndex].Answer = Wrong
	printWithColorln(fmt.Sprintf("Answer has saved as '%s'", Wrong), red)
}

func markAnswerAsNeutral() {
	questionsPerTopic[questionIndex].Answer = Neutral
	printWithColorln(fmt.Sprintf("Answer has saved as '%s'", Neutral), magenta)
}

func saveInterview() {

}
