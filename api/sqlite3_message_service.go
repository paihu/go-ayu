package main

import "github.com/jmoiron/sqlx"

type Sqlite3MessageService struct {
	DB *sqlx.DB
}

func (s Sqlite3MessageService) Add(m Message) (Message, error) {

	r, err := s.DB.NamedExec("insert into message(post_id,user,email,message) values(:post_id,:user,:email,:message)", m)
	if err != nil {
		return Message{}, err
	}
	m.Id, err = r.LastInsertId()
	if err != nil {
		return Message{}, err
	}
	return m, nil
}

func (s Sqlite3MessageService) Gets(p Post) (m []Message, err error) {
	err = s.DB.Select(&m, "select * from message where post_id = ? order by inserted_at asc", p.Id)
	return
}

func (s Sqlite3MessageService) Get(m Message) (Message, error) {
	err := s.DB.Get(&m, "select * from message where id = ?", m.Id)
	return m, err
}
