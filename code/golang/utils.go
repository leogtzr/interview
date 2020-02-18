package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	case "load":
		return loadCmd, fullCommand[1:]
	case "exf":
		return exitInterviewFile, []string{}
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

func retrieveTopicsFromFileSystem(interviewsDir string) []string {
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

func retrieveTopicsFromInterview(topics *map[string]Questions) []string {
	tps := make([]string, 0)
	for t := range *topics {
		tps = append(tps, t)
	}
	return tps
}

func listTopicsFromInterviewFile(topics *map[string]Questions) {
	if usingInterviewFile {
		topics := retrieveTopicsFromInterview(&interview.Topics)
		for _, topic := range topics {
			fmt.Println(termenv.String(topic).Underline().Bold())
		}
	}
}

func listTopics(interviewsDir string) {
	topics := retrieveTopicsFromFileSystem(interviewsDir)
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
	if (nextID == 0) || (nextID == id) {
		nextID = -1
	}
	return Question{ID: int(id), Q: q, Answer: NotAnsweredYet}
}

func extractTopicName(options []string) string {
	topicName := options[0]
	topicName = strings.ToLower(topicName)
	return topicName
}

func setTopicFromFileSystem(options []string) {
	topicName := extractTopicName(options)
	topics := retrieveTopicsFromFileSystem(interviewTopicsDir)

	if topicExist(topicName, &topics) &&
		exists(filepath.Join(interviewTopicsDir, "topics", topicName, "questions")) {
		selectedTopic = topicName
		questionsPerTopic := loadQuestionsFromTopic(selectedTopic, interviewTopicsDir)
		interview.Topics[selectedTopic] = questionsPerTopic
	} else {
		fmt.Println(
			termenv.String(fmt.Sprintf("topic '%s' not found or the topic selected doesn't have questions.", topicName)).Foreground(colorProfile.Color(red)))
	}
}

func setTopicFrom(options []string, topicsFromInterviewFile *map[string]Questions) {
	topicName := extractTopicName(options)
	topics := retrieveTopicsFromInterview(topicsFromInterviewFile)
	if topicExist(topicName, &topics) {
		selectedTopic = topicName
		return
	}

	fmt.Println(
		termenv.String(fmt.Sprintf("topic '%s' not found or the topic selected doesn't have questions.", topicName)).Foreground(colorProfile.Color(red)))
}

func loadQuestionsFromTopic(topic, interviewsDir string) []Question {
	// Clear previous questions ...
	questionsPerTopic := make([]Question, 0)

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

	return questionsPerTopic
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

func ps1String(ps1, selectedTopic, intervieweeName string) string {
	if selectedTopic == "" {
		return "$ "
	}
	return fmt.Sprintf("/%s %s $ ", termenv.String(selectedTopic).Faint(), shortIntervieweeName(intervieweeName, minNumberOfCharsInIntervieweeName))
}

func isQuestionFormatValid(question string, rgx *regexp.Regexp) bool {
	return rgx.MatchString(question)
}

func (q Question) String() string {
	return fmt.Sprintf("Q%d: %s [%s]", q.ID, q.Q, q.Answer)
}

func printQuestion(questionIndex int) {
	if hasStarted && (len(interview.Topics[selectedTopic]) > 0) {
		fmt.Println(interview.Topics[selectedTopic][questionIndex])
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

	if (questionIndex + 1) < len(interview.Topics[selectedTopic]) {
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
	if len(interview.Topics[selectedTopic]) < 1 {
		printWithColorln("You need to select a topic first.", red)
		return
	}
	for _, q := range interview.Topics[selectedTopic] {
		fmt.Printf("[%s]\n", q)
	}
}

func readIntervieweeName(stdin io.Reader) (string, bool) {
	reader := bufio.NewScanner(stdin)
	reader.Scan()
	text := reader.Text()
	if len(strings.TrimSpace(text)) == 0 {
		return "", false
	}
	return strings.TrimSpace(text), true
}

func printWithColorln(msg, colorCode string) {
	fmt.Println(termenv.String(msg).Foreground(colorProfile.Color(colorCode)))
}

func saveInterview() error {
	intervieweeName := interview.Interviewee
	savedDir := filepath.Join(interviewTopicsDir, "saved")
	if !dirExists(savedDir) {
		return fmt.Errorf("[%s] does not exist", savedDir)
	}

	savedInterviewName := filepath.Join(savedDir, intervieweeName)
	if dirExists(savedInterviewName) {
		printWithColorln(fmt.Sprintf("[%s] already exists, we will generate another name.", savedInterviewName), red)
		seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
		randDirName := stringWithCharset(2, charset, seededRand)
		savedInterviewName = fmt.Sprintf("%s-%s", savedInterviewName, randDirName)
	}
	err := os.MkdirAll(savedInterviewName, os.ModePerm)
	if err != nil {
		return err
	}
	return saveData(filepath.Join(savedInterviewName, "interview"), interview)
}

func saveData(savedInterviewNamePath string, interview Interview) error {
	file, err := os.Create(savedInterviewNamePath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	fmt.Fprintf(w, "%s@%s\n", interview.Interviewee, interview.Date.Format(interviewFormatLayout))

	for topicName, questions := range interview.Topics {
		for _, q := range questions {
			if q.Answer != NotAnsweredYet {
				fmt.Fprintf(w, "%s@%d@%s@%d\n", topicName, q.ID, q.Q, int(q.Answer))
			}
		}
	}
	return w.Flush()
}

func loadInterview(options []string) (Interview, error) {
	interviewName := strings.Join(options, " ")
	interviewFile := filepath.Join(interviewTopicsDir, "saved", interviewName, "interview")
	if !dirExists(interviewFile) {
		return Interview{}, fmt.Errorf("'%s' does not exist", interviewFile)
	}
	file, err := os.Open(interviewFile)
	if err != nil {
		return Interview{}, err
	}
	defer file.Close()

	interview := Interview{}

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		header := scanner.Text()
		intervieweeName, err := extractNameFromInterviewHeaderRecord(header)
		if err != nil {
			return Interview{}, err
		}
		interview.Interviewee = intervieweeName

		interviewDate, err := extractDateFromInterviewHeaderRecord(header)
		if err != nil {
			return Interview{}, err
		}
		interview.Date = interviewDate
	}

	interview.Topics = make(map[string]Questions)

	// Load questions:
	for scanner.Scan() {
		questionFileRecord := scanner.Text()
		topic, question := extractQuestionInfo(questionFileRecord)
		interview.Topics[topic] = append(interview.Topics[topic], question)
	}

	return interview, nil
}

func extractNameFromInterviewHeaderRecord(header string) (string, error) {
	fields := strings.Split(strings.TrimSpace(header), "@")
	if len(fields) != 2 {
		return "", fmt.Errorf("'%s' wrong header format", header)
	}
	return fields[0], nil
}

func extractDateFromInterviewHeaderRecord(header string) (time.Time, error) {
	fields := strings.Split(strings.TrimSpace(header), "@")
	if len(fields) != 2 {
		return time.Time{}, fmt.Errorf("'%s' wrong header format", header)
	}
	interviewDate, err := time.Parse(interviewFormatLayout, fields[1])
	return interviewDate, err
}

func extractQuestionInfo(questionFileRecord string) (string, Question) {
	fields := strings.Split(questionFileRecord, "@")
	topic := fields[0]
	id, _ := strconv.ParseInt(fields[1], 10, 64)
	question := fields[2]

	q := Question{}
	q.ID = int(id)
	q.Q = question
	x, _ := strconv.ParseInt(fields[4], 10, 64)
	q.Answer = Answer(int(x))

	return topic, q
}

func resetStatus() {
	interview = Interview{Topics: make(map[string]Questions)}
	usingInterviewFile = false
	hasStarted = false
	questionIndex = 0
	selectedTopic = ""
	ps1 = "$ "
}
