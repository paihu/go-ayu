package main

import (
	"fmt"
	"time"
)

var (
	messageService MessageService
)

type Message struct {
	Id         int64     `json:"id,omitempty" db:"id"`
	PostId     int64     `json:"post_id" db:"post_id"`
	User       string    `json:"user" db:"user"`
	Email      string    `json:"email" db:"email"`
	Message    string    `json:"message" db:"message"`
	InsertedAt time.Time `json:"inserted_at,omitempty" db:"inserted_at"`
}

type MessageService interface {
	Add(m Message) (Message, error)
	Get(m Message) (Message, error)
	Gets(p Post) ([]Message, error)
}

func (m Message) Validate() error {
	if m.PostId == 0 {
		return fmt.Errorf("requre post_id")
	}
	if m.User == "" {
		return fmt.Errorf("requre user")
	}
	if m.Message == "" {
		return fmt.Errorf("requre message")
	}
	if !validateEmail(m.Email) {
		return fmt.Errorf("email validation error")
	}
	return nil
}
func (m Message) Add() (Message, error) {
	return messageService.Add(m)
}

func (m Message) Get() (Message, error) {
	return messageService.Get(m)
}

func GetMessages(post_id int64) ([]Message, error) {
	p := Post{Id: post_id}

	p, err := postService.Get(p)
	if err != nil {
		return []Message{}, err
	}
	return messageService.Gets(p)
}
