package main

import (
	"database/sql"
)

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

func saveAnswer(question *Question, result Answer, config *Config, db *sql.DB) error {

	intervieweeID := config.intervieweeID
	exists, err := existsAnswer(intervieweeID, question.ID, db)
	if err != nil {
		panic(err)
	}
	if exists {
		if err := updateAnswer(question, result, config, db); err != nil {
			return err
		}
	} else {
		stmt, err := db.Query(`insert into answer (result, question_id, candidate_id) values(?, ?, ?)`,
			result, question.ID, intervieweeID)
		if err != nil {
			return err
		}
		defer stmt.Close()
	}

	return nil
}

func updateAnswer(q *Question, result Answer, config *Config, db *sql.DB) error {
	intervieweeID := config.intervieweeID
	comment := config.comment
	if _, err :=
		db.Exec(`update answer set result = ?, comment = ? where question_id = ? and candidate_id = ?`,
			result, comment, q.ID, intervieweeID); err != nil {
		return err
	}
	return nil
}

func getQuestionsByTopic(topic string, db *sql.DB) ([]Question, error) {
	questionsPerTopic := make([]Question, 0)

	results, err :=
		db.Query(
			`select q.id, question, q.level_id from question q, topic t where t.topic = ? and t.id = q.topic_id`,
			topic)
	if err != nil {
		return []Question{}, err
	}
	defer results.Close()

	for results.Next() {
		var question Question
		err = results.Scan(&question.ID, &question.Q, &question.Level)
		if err != nil {
			return []Question{}, err
		}
		questionsPerTopic = append(questionsPerTopic, question)
	}

	return questionsPerTopic, nil
}

func getQuestionsByTopicWithLevel(topic string, level Level, db *sql.DB) ([]Question, error) {
	questionsPerTopic := make([]Question, 0)

	results, err :=
		db.Query(
			`select q.id, question, q.level_id from question q, topic t where t.topic = ? and t.id = q.topic_id and level_id = ?`,
			topic, level)
	if err != nil {
		return []Question{}, err
	}
	defer results.Close()

	for results.Next() {
		var question Question
		err = results.Scan(&question.ID, &question.Q, &question.Level)
		if err != nil {
			return []Question{}, err
		}
		questionsPerTopic = append(questionsPerTopic, question)
	}

	return questionsPerTopic, nil
}

func getTopicsWithQuestions(db *sql.DB) ([]string, error) {
	var topics []string
	results, err := db.Query("select distinct(t.topic) from topic t inner join question q on t.id = q.topic_id")
	if err != nil {
		return []string{}, err
	}
	defer results.Close()

	for results.Next() {
		var topic string
		err = results.Scan(&topic)
		if err != nil {
			return []string{}, err
		}
		topics = append(topics, topic)
	}

	return topics, nil
}

func saveIntervieweeName(interviewee string, db *sql.DB) (int, error) {
	stmt, err := db.Exec("insert into candidate(name, date) values(?, now())", interviewee)
	if err != nil {
		return -1, err
	}
	id, err := stmt.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

func existsAnswer(candidateID, questionID int, db *sql.DB) (bool, error) {
	results, err :=
		db.Query(`select count(*) from answer where candidate_id = ? and question_id = ?`, candidateID, questionID)
	if err != nil {
		return false, err
	}
	defer results.Close()

	var count int
	if results.Next() {
		err = results.Scan(&count)
		if err != nil {
			return false, err
		}
	}

	return count > 0, err
}

func saveQuestion(q *Question, topicID int, answer string, db *sql.DB) error {
	_, err := db.Exec(`insert into question (question, answer, topic_id, level_id) values(?, ?, ?, ?)`, q.Q, answer, topicID, q.Level)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}
