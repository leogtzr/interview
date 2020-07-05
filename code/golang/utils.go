package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/muesli/termenv"
	"github.com/spf13/viper"
)

func sanitizeUserInput(input string) string {
	return strings.TrimSpace(input)
}

// Transforms user's input to a Command
func userInputToCmd(input string) (Command, []string) {
	if len(input) == 0 {
		return noCmd, []string{}
	}
	fullCommand := words(input)
	input = fullCommand[0]
	input = sanitizeUserInput(input)
	input = strings.ToLower(input)
	switch input {
	case "exit", "quit", ":q", "/q", "q":
		return exitCmd, []string{}
	case "topics", "tps", "t", "/t", ":t":
		return topicsCmd, []string{}
	case "help", ":h", "/h", "--h", "-h", "h":
		return helpCmd, []string{}
	case "use", "u", "/u", ":u", "-u", "--u", "set":
		if len(fullCommand) <= 1 {
			return noCmd, []string{}
		}
		return useCmd, fullCommand[1:]
	case "cls", "clear":
		return clearScreenCmd, []string{}
	case "pwd":
		return pwdCmd, []string{}
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
	case "count", "cnt", "c":
		return countCmd, []string{}
	case "cmt", "comment", "note", "nt":
		return createCommentCmd, []string{}
	case "cq":
		return createQuestionCmd, []string{}
	case "va":
		return viewCurrentQuestionAnwswerCmd, []string{}
	case "vas":
		return viewAnswers, []string{}
	}
	return noCmd, []string{}
}

func listTopics(db *sql.DB) error {
	topics, err := getTopics(db)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		fmt.Println(termenv.String(topic.Topic).Underline().Bold())
	}
	return nil
}

func printHelp() {
	usage := `
commands:

	exit|quit|:q|/q|q 			exits from this application.
	topics|tps|t|/t|:t 			list current available topics from the DB
	help|:h|/h|--h|-h 			shows this message.
	use|u|/u|:u|-u|--u|set 			sets an available topic.
	cls|clear 				clears the screen.
	pwd 					prints the current selected topic.
	start|begin 				starts the interview.
	print|print()|p|p() 			prints the current question.
	next|nxt|> 				moves to the next question.
	previous|prev|< 			moves to the previous question.
	view|v					prints the current available questions by level.
	va				view answer from current question
	no|n|mal|wrong|nop|bad|nel 		marks a question as wrong.
	ok|yes|si|right|y			marks a question as right / OK.
	hmm|meh|?				marks a question as neutral.
	finish|done|bye				finishes an interview.
	cq					create a question and save it to the database.
	+					increases the level of the interview, it could be from Programmer Analyst to a Sr Programmer Analyst as an example.
	- 					decreases the level of the interview.
	= 					ignore levels.
	lvl					prints the current interview level.
	stats					shows some stats and the current configuration for the interview.
	ap					sets the level of the interview to "Associate Programmer"
	pa					sets the level of the interview to "Programmer Analyst"
	sr					sets the level of the interview to "Sr Programmer Analyst"


	Any other command or sentence that is not listed here will be simply ignored.
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

func extractTopicName(options []string) string {
	topicName := options[0]
	topicName = strings.ToLower(topicName)
	return topicName
}

func setTopic(options []string, config *Config, db *sql.DB) error {
	topicName := extractTopicName(options)
	topics, err := getTopicsWithQuestions(db)
	if err != nil {
		return err
	}

	if topicExist(topicName, &topics) {
		config.selectedTopic = topicName
		questionsPerTopic, err := loadQuestionsFromTopic(config, db)
		if err != nil {
			return err
		}
		config.interview.Topics[config.selectedTopic] = questionsPerTopic
	} else {
		fmt.Println(
			termenv.String(fmt.Sprintf("topic '%s' not found or the topic selected doesn't have questions.", topicName)).Foreground(config.colorProfile.Color(red)))
	}
	return nil
}

func loadQuestionsFromTopic(config *Config, db *sql.DB) ([]Question, error) {
	// Clear previous questions ...
	questionsPerTopic, err := getQuestionsByTopic(config.selectedTopic, db)
	if err != nil {
		return []Question{}, err
	}

	levelFound := findLevel(&questionsPerTopic, AssociateOrProgrammer, ProgrammerAnalyst, SrProgrammer)
	fmt.Printf("Loaded -> '%d' questions, starting with: %s level.\n", len(questionsPerTopic), levelFound)

	levelQCounts := levelQuestionCounts(&questionsPerTopic)
	fmt.Printf("Associate Programmer = ")
	printWithColorf(config, "%d\n", green, levelQCounts[AssociateOrProgrammer])
	fmt.Printf("Programmer Analyst = ")
	printWithColorf(config, "%d\n", green, levelQCounts[ProgrammerAnalyst])
	fmt.Printf("Sr. Programmer  = ")
	printWithColorf(config, "%d\n", green, levelQCounts[SrProgrammer])

	return questionsPerTopic, nil
}

func levelQuestionCounts(qs *[]Question) map[Level]int {
	counts := make(map[Level]int)
	for _, q := range *qs {
		counts[q.Level]++
	}
	return counts
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
	return fmt.Sprintf(
		"/%s %s $ ",
		termenv.String(selectedTopic).Faint(), shortIntervieweeName(intervieweeName, minNumberOfCharsInIntervieweeName))
}

func (q Question) String() string {
	if int(q.Result) == 1 || int(q.Result) == 0 {
		return fmt.Sprintf("Q%d: %s [%s]", q.ID, q.Q, q.Level)
	}
	return fmt.Sprintf("Q%d: %s [%s] [%s]", q.ID, q.Q, q.Result, q.Level)
}

// StringNoResult ...
func (q Question) StringNoResult() string {
	return fmt.Sprintf("Q%d: %s [%s]", q.ID, q.Q, q.Level)
}

func printQuestion(questionIndex int, config *Config) {
	if !config.hasStarted {
		return
	}

	if config.ignoreLevelChecking && (len(config.interview.Topics[config.selectedTopic]) > 0) {
		printWithColorln(config.interview.Topics[config.selectedTopic][config.questionIndex].String(), gray, config)
		fmt.Println()
		return
	}
	currentLevel := config.levels[config.levelIndex]
	currentLevelQuestions := getQuestionsFromLevel(currentLevel, config)
	if len(currentLevelQuestions) == 0 {
		printWithColorln("There are no questions for this level.", yellow, config)
		fmt.Println()
		return
	}
	index := config.individualLevelIndexes[int(currentLevel)-1]
	fmt.Println(currentLevelQuestions[index])
	fmt.Println()
}

func viewQuestions(config *Config) {
	if len(config.interview.Topics[config.selectedTopic]) < 1 {
		printWithColorln("You need to select a topic first.", red, config)
		fmt.Println()
		return
	}
	for _, q := range config.interview.Topics[config.selectedTopic] {
		fmt.Println(q.StringNoResult())
	}
}

func viewQuestionsByLevel(config *Config) {
	if len(config.selectedTopic) == 0 {
		printWithColorln("You need to select a topic first.", red, config)
		return
	}
	currentLevel := config.levels[config.levelIndex]
	currentLevelQuestions := getQuestionsFromLevel(currentLevel, config)
	for _, q := range currentLevelQuestions {
		fmt.Println(q.StringNoResult())
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

func printWithColorln(msg, colorCode string, config *Config) {
	fmt.Println(termenv.String(msg).Foreground(config.colorProfile.Color(colorCode)))
}

func printWithColorf(config *Config, msg, colorCode string, a ...interface{}) {
	fmt.Printf(termenv.String(msg).Foreground(config.colorProfile.Color(colorCode)).String(), a...)
}

func resetStatus(config *Config) {
	config.interview = Interview{Topics: make(map[string][]Question)}
	//config.usingInterviewFile = false
	config.hasStarted = false
	config.questionIndex = 0
	config.selectedTopic = ""
	config.ps1 = "$ "
}

func showLevel(config *Config) {
	currentLevel := config.levels[config.levelIndex]
	printWithColorln(currentLevel.String(), cyan, config)
}

func setAnswerAsNeutral(config *Config, db *sql.DB) error {
	questions := config.interview.Topics[config.selectedTopic]
	q := questions[config.questionIndex]
	q.Result = Neutral

	if err := saveAnswer(&q, Neutral, config, db); err != nil {
		return err
	}

	currentLevel := config.levels[config.levelIndex]
	currentLevelQuestions := getQuestionsFromLevel(currentLevel, config)
	index := config.individualLevelIndexes[int(currentLevel)-1]
	id := currentLevelQuestions[index].ID
	qs := config.interview.Topics[config.selectedTopic]
	markQuestionAs(id, Neutral, &qs)

	printWithColorln(fmt.Sprintf("Answer has saved as '%s'", Neutral), magenta, config)
	return nil
}

func setAnswerAsOK(config *Config, db *sql.DB) error {
	questions := config.interview.Topics[config.selectedTopic]
	q := questions[config.questionIndex]
	q.Result = OK

	if err := saveAnswer(&q, OK, config, db); err != nil {
		return err
	}

	currentLevel := config.levels[config.levelIndex]
	currentLevelQuestions := getQuestionsFromLevel(currentLevel, config)
	index := config.individualLevelIndexes[int(currentLevel)-1]
	id := currentLevelQuestions[index].ID
	qs := config.interview.Topics[config.selectedTopic]
	markQuestionAs(id, OK, &qs)

	printWithColorln(fmt.Sprintf("Answer has saved as '%s'", OK), green, config)
	return nil
}

func setAnswerAsWrong(config *Config, db *sql.DB) error {
	questions := config.interview.Topics[config.selectedTopic]
	q := questions[config.questionIndex]
	q.Result = Wrong

	if err := saveAnswer(&q, Wrong, config, db); err != nil {
		return err
	}

	currentLevel := config.levels[config.levelIndex]
	currentLevelQuestions := getQuestionsFromLevel(currentLevel, config)
	index := config.individualLevelIndexes[int(currentLevel)-1]
	id := currentLevelQuestions[index].ID
	qs := config.interview.Topics[config.selectedTopic]
	markQuestionAs(id, Wrong, &qs)

	printWithColorln(fmt.Sprintf("Answer has saved as '%s'", Wrong), red, config)
	return nil
}

func answerAs(config *Config, ans Result, messageColorCode string, db *sql.DB) error {
	currentLevel := config.levels[config.levelIndex]
	currentLevelQuestions := getQuestionsFromLevel(currentLevel, config)
	index := config.individualLevelIndexes[int(currentLevel)-1]
	id := currentLevelQuestions[index].ID
	q := currentLevelQuestions[index]
	qs := config.interview.Topics[config.selectedTopic]
	markQuestionAs(id, ans, &qs)
	if err := saveAnswer(&q, ans, config, db); err != nil {
		return err
	}
	printWithColorln(fmt.Sprintf("Answer has saved as '%s'", ans), messageColorCode, config)
	return nil
}

func markQuestionAs(id int, ans Result, qs *[]Question) {
	for _, q := range *qs {
		if q.ID == id {
			(*qs)[id-1].Result = ans
			break
		}
	}
}

func showStats(config *Config, db *sql.DB) error {
	currentLevel := config.levels[config.levelIndex]

	if len(config.selectedTopic) == 0 {
		fmt.Printf("Level: ")
		printWithColorf(config, "%s\n", green, currentLevel)

		fmt.Printf("Ignoring level: ")
		printWithColorf(config, "%t\n", green, config.ignoreLevelChecking)

		fmt.Printf("Questions in bucket: ")
		printWithColorf(config, "%t\n", green, len(config.selectedTopic) != 0)
	} else {
		counts, err := getResultCounts(config.intervieweeID, db)
		if err != nil {
			return err
		}
		resultCounts := resultCount(&counts)
		notAnsweredCount := resultCounts[NotAnsweredYet]
		okCount := resultCounts[OK]
		wrongCount := resultCounts[Wrong]
		neutralCount := resultCounts[Neutral]
		total := notAnsweredCount + okCount + wrongCount + neutralCount

		if !config.ignoreLevelChecking {
			fmt.Printf("Level: ")
			printWithColorf(config, "%s\n", green, currentLevel)
		}

		fmt.Printf("Questions in bucket: ")
		printWithColorf(config, "%t\n", green, len(config.selectedTopic) != 0)

		fmt.Printf("Not Answered: ")
		printWithColorf(config, "%d (%.2f%%)\n", green, notAnsweredCount, perc(notAnsweredCount, total))

		fmt.Printf("OK: ")
		printWithColorf(config, "%d (%.2f%%)\n", green, okCount, perc(okCount, total))

		fmt.Printf("Wrong: ")
		printWithColorf(config, "%d (%.2f%%)\n", green, wrongCount, perc(wrongCount, total))

		fmt.Printf("Neutral: ")
		printWithColorf(config, "%d (%.2f%%)\n", green, neutralCount, perc(neutralCount, total))
	}
	return nil
}

func resultCount(counts *[]ResultCount) map[Result]int {
	rCounts := make(map[Result]int, 0)

	for _, v := range *counts {
		rCounts[Result(v.Result)] = v.Count
	}

	for _, v := range [4]Result{NotAnsweredYet, OK, Wrong, Neutral} {
		if _, ok := rCounts[v]; !ok {
			rCounts[v] = 0
		}
	}

	return rCounts
}

func perc(count, total int) float64 {
	if count == 0 {
		return 0.0
	}
	return (float64(count) * 100.0) / float64(total)
}

func setLevel(lvl Level, config *Config) {
	config.levelIndex = int(lvl) - 1
	currentLevel := config.levels[config.levelIndex]
	fmt.Printf("Current level is: ")
	printWithColorln(fmt.Sprintf("%s", currentLevel), green, config)
}

func showCounts(config *Config) {
	qs := config.interview.Topics[config.selectedTopic]
	levelQCounts := levelQuestionCounts(&qs)
	fmt.Printf("Associate Programmer = ")
	printWithColorf(config, "%d\n", green, levelQCounts[AssociateOrProgrammer])
	fmt.Printf("Programmer Analyst = ")
	printWithColorf(config, "%d\n", green, levelQCounts[ProgrammerAnalyst])
	fmt.Printf("Sr. Programmer  = ")
	printWithColorf(config, "%d\n", green, levelQCounts[SrProgrammer])
}

// NewConfig Creates a new Configuration object.
func NewConfig() Config {
	cfg := Config{}
	cfg.selectedTopic = ""
	cfg.ps1 = "$ "
	cfg.colorProfile = termenv.ColorProfile()
	cfg.interview = Interview{Topics: make(map[string][]Question)}
	cfg.topicQuestionsLevel = AssociateOrProgrammer
	cfg.levelIndex = 0
	cfg.ignoreLevelChecking = false
	cfg.individualLevelIndexes = []int{0, 0, 0}
	cfg.questionIndex = 0
	cfg.levels = [3]Level{
		AssociateOrProgrammer, ProgrammerAnalyst, SrProgrammer,
	}
	return cfg
}

func readConfig(filename, configPath string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(configPath)
	v.SetConfigType("env")
	err := v.ReadInConfig()
	return v, err
}

func readComment() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var b strings.Builder
	for scanner.Scan() {
		b.WriteString(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return b.String(), nil
}

func makeQuestion(config *Config, db *sql.DB) error {
	topics, err := getTopics(db)
	if err != nil {
		return err
	}

	for idx, topic := range topics {
		printWithColorf(config, "%d: %s\n", blue, idx, topic.Topic)
	}
	fmt.Println()
	fmt.Printf("Topic? ")

	reader := bufio.NewReader(os.Stdin)
	userInput, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	userInput = strings.TrimSpace(userInput)
	topicIndex, err := strconv.Atoi(userInput)
	if err != nil {
		return err
	}

	if topicIndex < 0 || topicIndex > len(topics) {
		return errors.New("invalid topic index")
	}

	printWithColorf(config, "\n1) Programmer\n2) Programmer Analyst\n3) Sr. Programmer Analyst ", blue)
	fmt.Println()
	fmt.Printf("Level? ")
	userInput, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	userInput = strings.TrimSpace(userInput)

	levelIndex, err := strconv.Atoi(userInput)
	if err != nil {
		return err
	}
	if levelIndex < 0 || levelIndex > 3 {
		return errors.New("invalid level index")
	}

	fmt.Printf("Question? ")
	userInput, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	question := strings.TrimSpace(userInput)

	fmt.Println()
	fmt.Printf("Answer? ")
	userInput, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	answer := strings.TrimSpace(userInput)

	q := Question{Q: question, Level: Level(levelIndex), Result: NotAnsweredYet}

	if err = saveQuestion(&q, topicIndex, answer, db); err != nil {
		return err
	}

	return nil
}

func viewAnswer(questionIndex int, config *Config) {
	if !config.hasStarted {
		return
	}
	if config.ignoreLevelChecking && (len(config.interview.Topics[config.selectedTopic]) > 0) {
		printWithColorln(config.interview.Topics[config.selectedTopic][config.questionIndex].Answer, gray, config)
		fmt.Println()
		return
	}
	currentLevel := config.levels[config.levelIndex]
	currentLevelQuestions := getQuestionsFromLevel(currentLevel, config)
	if len(currentLevelQuestions) == 0 {
		printWithColorln("There are no questions for this level.", yellow, config)
		fmt.Println()
		return
	}
	index := config.individualLevelIndexes[int(currentLevel)-1]
	fmt.Println(currentLevelQuestions[index].Answer)
	fmt.Println()
}

func listAnswers(config *Config, db *sql.DB) error {
	id := config.intervieweeID
	answers, err := getAnswersFromCandidate(id, db)
	if err != nil {
		return err
	}
	for _, ans := range answers {
		fmt.Println(ans)
	}
	return nil
}
