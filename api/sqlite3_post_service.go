package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Sqlite3PostService struct {
	DB *sqlx.DB
}

func (s Sqlite3PostService) Delete(p Post) error {
	_, err := s.DB.Exec("delete from post where post_id = ? ", p.Id)
	return err
}
func (s Sqlite3PostService) Add(p Post) (Post, error) {
	query := `
INSERT INTO post(
user,
email,
title,
delete_password,
url,
kind,
uploaded_at,
limit_count,
require_message,
limit_time,
comment
)
VALUES(
:user,
:email,
:title,
:delete_password,
:url,
:kind,
:uploaded_at,
:limit_count,
:require_message,
:limit_time,
:comment
)`
	r, err := s.DB.NamedExec(query, p)
	if err != nil {
		return Post{}, err
	}
	p.Id, err = r.LastInsertId()
	if err != nil {
		return Post{}, err
	}
	return p, nil
}
func (s Sqlite3PostService) Get(p Post) (r Post, err error) {
	err = s.DB.Get(&r, "select * from post where id = ?", p.Id)
	if err != nil {
		return Post{}, err
	}
	return
}
func (s Sqlite3PostService) CountIngrement(p Post) (r Post, err error) {
	tx := s.DB.MustBegin()
	defer tx.Rollback()
	err = tx.Get(&r, "select * from post where id = ?", p.Id)
	if err != nil {
		fmt.Println(err)
		return Post{}, err
	}
	_, err = tx.Exec("update post set count= ? where id = ?", r.Count+1, r.Id)
	if err != nil {
		fmt.Println(err)
		return Post{}, err
	}
	err = tx.Commit()
	return
}
func (s Sqlite3PostService) Gets(page int) (posts []Post, err error) {
	err = s.DB.Select(&posts, "select * from post order by uploaded_at limit 30 offset ?", page*30)
	fmt.Println(posts)
	return
}

func (s Sqlite3PostService) GetCounts() (count int) {
	err := s.DB.Get(&count, "select count(*) from post")
	if err != nil {
		return 0
	}
	return
}
