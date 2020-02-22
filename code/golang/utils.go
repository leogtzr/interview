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
	case "+":
		return increaseLevelCmd, []string{}
	case "-":
		return decreaseLevelCmd, []string{}
	case "=":
		return ignoreLevelCmd, []string{}
	case "lvl":
		return showLevelCmd, []string{}
	case "stats":
		return showStatsCmd, []string{}
	case "ap":
		return setAssociateProgrammerLevelCmd, []string{}
	case "pa":
		return setProgrammerAnalystLevelCmd, []string{}
	case "sr":
		return setSRProgrammerLevelCmd, []string{}
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

func retrieveTopicsFromInterview(topics *map[string][]Question) []string {
	tps := make([]string, 0)
	for t := range *topics {
		tps = append(tps, t)
	}
	return tps
}

func listTopicsFromInterviewFile(topics *map[string][]Question) {
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
	level, _ := strconv.ParseInt(questionFields[2], 10, 64)
	return Question{ID: int(id), Q: q, Answer: NotAnsweredYet, Level: Level(level)}
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

func setTopicFrom(options []string, topicsFromInterviewFile *map[string][]Question) {
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
		if strings.HasPrefix(questionText, "#") {
			fmt.Println("Filtering here ... ")
			continue
		}
		if len(strings.TrimSpace(questionText)) == 0 {
			fmt.Println("Filtering here 2 ")
			continue
		}
		if isQuestionFormatValid(questionText, rgxQuestions) {
			question := toQuestion(questionText)
			questionsPerTopic = append(questionsPerTopic, question)
		}
	}

	levelFound := findLevel(&questionsPerTopic, AssociateOrProgrammer, ProgrammerAnalyst, SrProgrammer)
	fmt.Printf("Loaded -> '%d' questions, starting with: %s level.\n", len(questionsPerTopic), levelFound)

	return questionsPerTopic
}

func shortIntervieweeName(name string, min int) string {
	if len(name) == 0 {
		return "(who?)"
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
	return fmt.Sprintf("Q%d: %s [%s] [%s]", q.ID, q.Q, q.Answer, q.Level)
}

func printQuestion(questionIndex int) {
	if !hasStarted {
		return
	}

	if ignoreLevelChecking && (len(interview.Topics[selectedTopic]) > 0) {
		fmt.Println(interview.Topics[selectedTopic][questionIndex])
	} else {
		currentLevel := levels[levelIndex]
		currentLevelQuestions := getQuestionsFromLevel(currentLevel, selectedTopic, &interview.Topics)
		index := individualLevelIndexes[int(currentLevel)-1]
		fmt.Println(currentLevelQuestions[index])
	}
}

func viewStats() {
	if len(interview.Topics[selectedTopic]) < 1 {
		printWithColorln("You need to select a topic first.", red)
		return
	}
	for _, q := range interview.Topics[selectedTopic] {
		fmt.Println(q)
	}
}

func viewStatsByLevel() {
	if len(selectedTopic) == 0 {
		printWithColorln("You need to select a topic first.", red)
		return
	}
	currentLevel := levels[levelIndex]
	currentLevelQuestions := getQuestionsFromLevel(currentLevel, selectedTopic, &interview.Topics)
	for _, q := range currentLevelQuestions {
		fmt.Println(q)
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

func printWithColorf(msg, colorCode string, a ...interface{}) {
	fmt.Printf(termenv.String(msg).Foreground(colorProfile.Color(colorCode)).String(), a...)
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

	interview.Topics = make(map[string][]Question)

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

	q := Question{ID: int(id), Q: question}
	x, _ := strconv.ParseInt(fields[4], 10, 64)
	q.Answer = Answer(int(x))

	return topic, q
}

func resetStatus() {
	interview = Interview{Topics: make(map[string][]Question)}
	usingInterviewFile = false
	hasStarted = false
	questionIndex = 0
	selectedTopic = ""
	ps1 = "$ "
}

func showLevel() {
	currentLevel := levels[levelIndex]
	printWithColorln(currentLevel.String(), cyan)
}

func setAnswerAsNeutral(questions *[]Question, idx int) {
	(*questions)[idx].Answer = Neutral
	printWithColorln(fmt.Sprintf("Answer has saved as '%s'", Neutral), magenta)
}

func setAnswerAsNeutralWithLevel() {
	currentLevel := levels[levelIndex]
	currentLevelQuestions := getQuestionsFromLevel(currentLevel, selectedTopic, &interview.Topics)
	index := individualLevelIndexes[int(currentLevel)-1]
	id := currentLevelQuestions[index].ID
	markQuestionAs(id, Neutral)
	printWithColorln(fmt.Sprintf("Answer has saved as '%s'", Neutral), magenta)

}

func setAnswerAsOK(questions *[]Question, idx int) {
	(*questions)[idx].Answer = OK
	printWithColorln(fmt.Sprintf("Answer has saved as '%s'", OK), green)
}

func setAnswerAsOkWithLevel() {
	currentLevel := levels[levelIndex]
	currentLevelQuestions := getQuestionsFromLevel(currentLevel, selectedTopic, &interview.Topics)
	index := individualLevelIndexes[int(currentLevel)-1]
	id := currentLevelQuestions[index].ID
	markQuestionAs(id, OK)
	printWithColorln(fmt.Sprintf("Answer has saved as '%s'", OK), green)
}

func setAnswerAsWrong(questions *[]Question, idx int) {
	(*questions)[idx].Answer = Wrong
	printWithColorln(fmt.Sprintf("Answer has saved as '%s'", Wrong), red)
}

func setAnswerAsWrongWithLevel() {
	currentLevel := levels[levelIndex]
	currentLevelQuestions := getQuestionsFromLevel(currentLevel, selectedTopic, &interview.Topics)
	index := individualLevelIndexes[int(currentLevel)-1]
	id := currentLevelQuestions[index].ID
	markQuestionAs(id, Wrong)
	printWithColorln(fmt.Sprintf("Answer has saved as '%s'", Wrong), red)
}

func markQuestionAs(id int, ans Answer) {
	for _, q := range interview.Topics[selectedTopic] {
		if q.ID == id {
			interview.Topics[selectedTopic][id-1].Answer = ans
		}
	}
}

func showStats() {
	currentLevel := levels[levelIndex]

	if len(selectedTopic) == 0 {
		fmt.Printf("Level: %s\nIgnoring level: %t\nQuestions in bucket: %t",
			currentLevel,
			ignoreLevelChecking,
			len(selectedTopic) != 0)
	} else {
		counts := countGeneral(&interview.Topics)
		notAnsweredCount := counts[NotAnsweredYet]
		okCount := counts[OK]
		wrongCount := counts[Wrong]
		neutralCount := counts[Neutral]
		total := notAnsweredCount + okCount + wrongCount + neutralCount

		fmt.Printf("Level: ")
		printWithColorf("%s\n", green, currentLevel)

		fmt.Printf("Ignoring level: ")
		printWithColorf("%t\n", green, ignoreLevelChecking)

		fmt.Printf("Questions in bucket: ")
		printWithColorf("%t\n", green, len(selectedTopic) != 0)

		fmt.Printf("Not Answered: ")
		printWithColorf("%d (%.2f%%)\n", green, notAnsweredCount, perc(notAnsweredCount, total))

		fmt.Printf("OK: ")
		printWithColorf("%d (%.2f%%)\n", green, okCount, perc(okCount, total))

		fmt.Printf("Wrong: ")
		printWithColorf("%d (%.2f%%)\n", green, wrongCount, perc(wrongCount, total))

		fmt.Printf("Neutral: ")
		printWithColorf("%d (%.2f%%)\n", green, neutralCount, perc(neutralCount, total))
	}
}

func count(questions *[]Question, ans Answer) int {
	c := 0
	for _, q := range *questions {
		if q.Answer == ans {
			c++
		}
	}
	return c
}

func perc(count, total int) float64 {
	return (float64(count) * 100.0) / float64(total)
}

func countGeneral(topics *map[string][]Question) map[Answer]int {
	counts := make(map[Answer]int, 0)

	// flat the questions ...
	questions := make([]Question, 0)
	for _, qs := range *topics {
		for _, q := range qs {
			questions = append(questions, q)
		}
	}

	counts[NotAnsweredYet] = count(&questions, NotAnsweredYet)
	counts[OK] = count(&questions, OK)
	counts[Wrong] = count(&questions, Wrong)
	counts[Neutral] = count(&questions, Neutral)

	return counts
}

func setLevel(lvl Level, index *int, lvls [3]Level) {
	*index = int(lvl) - 1
	currentLevel := lvls[*index]
	fmt.Printf("Current level is: ")
	printWithColorln(fmt.Sprintf("%s", currentLevel), green)
}
