package main

import "database/sql"

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

func saveAnswer(question *Question, result Answer, intervieweeID int, db *sql.DB) error {
	stmt, err := db.Query(`insert into answer (result, question_id, candidate_id) values(?, ?, ?)`,
		result, question.ID, intervieweeID)
	if err != nil {
		return err
	}
	defer stmt.Close()
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
