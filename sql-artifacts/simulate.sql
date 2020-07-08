/*
mysql> select * from question where topic_id = 5 and level_id = 2;
+----+-----------------------------------------------------------------------+----------------------------------------------------------------+----------+----------+
| id | question                                                              | answer                                                         | topic_id | level_id |
+----+-----------------------------------------------------------------------+----------------------------------------------------------------+----------+----------+
| 29 | How can you set an environmental variable globally?                   |                                                                |        5 |        2 |
| 30 | Where can you define environmental variable?                          |                                                                |        5 |        2 |
| 31 | How can you find how much memory Linux is using?                      |                                                                |        5 |        2 |
| 32 | What are symbolic links?                                              |                                                                |        5 |        2 |
| 33 | How do you change file permissions under Linux?                       |                                                                |        5 |        2 |
| 34 | How do you create a hidden file?                                      |                                                                |        5 |        2 |
| 37 | How can you check if there is network connectivity between two hosts? |                                                                |        5 |        2 |
| 38 | How do you know if a port is opened in a remote host?                 |                                                                |        5 |        2 |
| 39 | How do you know if a file is being used by a process?                 | lsof                                                           |        5 |        2 |
| 41 | How can you run a command every two minutes?                          | cron or a custom script ... while; true; do sleep 2m ...; done |        5 |        2 |
| 42 | What is a very well know way of login to a remote server?             | ssh                                                            |        5 |        2 |
| 43 | How can you transfer a file from a remote host to your machine?       | scp or rsync                                                   |        5 |        2 |
| 44 | What is the default port used by SSH?                                 | 22                                                             |        5 |        2 |
| 45 | How can you run a command in background?                              | nohup, &, disown ...                                           |        5 |        2 |
| 50 | How to check if a process is using a specific port?                   | netstat                                                        |        5 |        2 |
| 51 | How can you redirect stderr to a file?                                |                                                                |        5 |        2 |
+----+-----------------------------------------------------------------------+----------------------------------------------------------------+----------+----------+

mysql> select * from candidate;
+----+--------+
| id | name   |
+----+--------+
|  1 | Leo    |
|  2 | Brenda |
|  3 | Mariel |
+----+--------+

mysql> select * from level;
+----+------------------------+
| id | title                  |
+----+------------------------+
|  1 | Programmer             |
|  2 | Programmer Analyst     |
|  3 | Sr. Programmer Analyst |
+----+------------------------+

mysql> show columns from answer;
+--------------+---------------+------+-----+---------+----------------+
| Field        | Type          | Null | Key | Default | Extra          |
+--------------+---------------+------+-----+---------+----------------+
| id           | int           | NO   | PRI | NULL    | auto_increment |
| result       | int           | NO   |     | NULL    |                |
| comment      | varchar(1000) | YES  |     | NULL    |                |
| question_id  | int           | NO   | MUL | NULL    |                |
| candidate_id | int           | NO   | PRI | NULL    |                |
+--------------+---------------+------+-----+---------+----------------+


*/
-- Some values for Mariel (3)
insert into answer (result, comment, question_id, candidate_id) values(0, "an interesting way of answering ...", 37, 3);
insert into answer (result, comment, question_id, candidate_id) values(1, "Almost OK ...", 38, 3);
