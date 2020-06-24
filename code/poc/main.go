package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Level ...
type Level struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// Question ...
type Question struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	TopicID  int    `json:"topic_id"`
	LevelID  int    `json:"level_id"`
}

// UserQuestionAnswer ...
type UserQuestionAnswer struct {
	ID         int    `json:"id"`
	Result     int    `json:"result"`
	Comment    string `json:"comment"`
	QuestionID int    `json:"question_id"`
}

func getLevels(db *sql.DB) ([]Level, error) {
	var levels []Level
	results, err := db.Query("SELECT * FROM level")
	if err != nil {
		return []Level{}, err
	}
	defer results.Close()

	for results.Next() {
		var lvl Level
		err = results.Scan(&lvl.ID, &lvl.Title)
		if err != nil {
			return []Level{}, err
		}
		levels = append(levels, lvl)
	}
	return levels, nil
}

func getTopics(db *sql.DB) ([]Topic, error) {
	var topics []Topic
	results, err := db.Query("SELECT * FROM topic")
	if err != nil {
		return []Topic{}, err
	}
	defer results.Close()

	for results.Next() {
		var topic Topic
		err = results.Scan(&topic.ID, &topic.Topic)
		if err != nil {
			return []Topic{}, err
		}
		topics = append(topics, topic)
	}
	return topics, nil
}

func main() {
	db, err := sql.Open("mysql", "root:lein23@/recruitment_interviews")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Hour * 3)
	defer db.Close()

	levels, err := getLevels(db)
	if err != nil {
		panic(err)
	}
	for _, lvl := range levels {
		fmt.Println(lvl)
	}

	topics, err := getTopics(db)
	if err != nil {
		panic(err)
	}
	for _, topic := range topics {
		fmt.Println(topic)
	}

}
