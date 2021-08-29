package main

import (
	"fmt"
	"time"
)

var (
	postService PostService
)

type Post struct {
	Id             int64     `json:"id,omitempty" db:"id"`
	User           string    `json:"user,omitempty" db:"user"`
	Email          string    `json:"email,omitempty" db:"email"`
	Title          string    `json:"title,omitempty" db:"title"`
	Comment        string    `json:"comment,omitempty" db:"comment"`
	Url            string    `json:"url,omitempty" db:"url"`
	DeletePassword string    `json:"delete_password,omitempty" db:"delete_password"`
	Kind           string    `json:"kind,omitempty" db:"kind"`
	Count          int       `json:"count,omitempty" db:"count,omitempty"`
	UploadedAt     time.Time `json:"uploaded_at,omitempty" db:"uploaded_at,omitempty"`
	LimitCount     int       `json:"limit_count,omitempty" db:"limit_count"`
	RequireMessage bool      `json:"require_message,omitempty" db:"require_message,omitempty"`
	LimitTime      time.Time `json:"limit_time,omitempty" db:"limit_time,omitempty"`
	Salt           string    `json:"-" db:"salt"`
}

type PostService interface {
	Delete(p Post) error
	Add(p Post) (Post, error)
	Get(p Post) (Post, error)
	CountIngrement(p Post) (Post, error)
	Gets(page int) ([]Post, error)
	GetCounts() int
}

func (p Post) Validate() error {
	if p.User == "" {
		return fmt.Errorf("requre user")
	}
	if p.Title == "" {
		return fmt.Errorf("requre title")
	}
	if p.Comment == "" {
		return fmt.Errorf("requre comment")
	}
	if p.DeletePassword == "" {
		return fmt.Errorf("requre delete_password")
	}
	if p.Kind == "" {
		return fmt.Errorf("requre kind")
	}
	if !validateEmail(p.Email) {
		return fmt.Errorf("email validation error")
	}
	if !validateUrl(p.Url) {
		return fmt.Errorf("url validation error")
	}
	return nil
}
func (p Post) Delete() error {
	return postService.Delete(p)
}

func (p Post) Add() (Post, error) {
	return postService.Add(p)
}

func (p Post) Get() (Post, error) {
	return postService.Get(p)
}
func (p Post) CountIngrement() (Post, error) {
	return postService.CountIngrement(p)
}
func GetPost(page int) ([]Post, error) {
	return postService.Gets(page)
}
