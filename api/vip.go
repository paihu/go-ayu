package main

type VIP struct {
	Id       int64  `json:"id,omitempty" db:"id"`
	User     string `json:"user" db:"user"`
	Password string `json:"password" db:"password"`
}
