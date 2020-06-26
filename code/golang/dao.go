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
