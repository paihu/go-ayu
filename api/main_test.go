package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db               *sqlx.DB
	test_duration, _ = time.ParseDuration("1h")
	test_now         = time.Now()
	location         = time.FixedZone("Z", 0*60*60)
	postJson         = fmt.Sprintf("{\"id\":1,\"user\":\"hoge\",\"email\":\"hoge@hoge.com\",\"title\":\"hogehoge\",\"comment\":\"comment\",\"url\":\"https://localhost/hoge.zip\",\"kind\":\"other\",\"uploaded_at\":\"%s\",\"limit_time\":\"%s\"}\n", test_now.Format(time.RFC3339Nano), test_now.Add(test_duration).Format(time.RFC3339Nano))
	messages         []Message
	h                handler
)

func TestMain(m *testing.M) {
	db, _ = sqlx.Open("sqlite3", ":memory:")
	defer db.Close()
	db.MustExec("PRAGMA foreign_keys = ON")
	db.SetMaxOpenConns(1)
	if err := initDB(db); err != nil {
		panic(err)
	}
	messages = make([]Message, 5)
	for i := range messages {
		messages[i].PostId = 1
		messages[i].User = "hoge"
		messages[i].Email = "hoge@hoge.com"
		messages[i].Message = fmt.Sprintf("message%d", i)
	}
	db.MustExec("insert into post(user,email,title,delete_password,url,kind,comment,uploaded_at,limit_time) values(?,?,?,?,?,?,?,?,?)",
		"hoge", "hoge@hoge.com", "hogehoge", "password", "https://localhost/hoge.zip", "other", "comment",
		test_now, test_now.Add(test_duration))
	if _, err := db.NamedExec("insert into message(post_id, `user`, email, message) values(:post_id,:user,:email,:message)", messages); err != nil {
		fmt.Println("insert error", err)
	}
	messageService = Sqlite3MessageService{db}
	postService = Sqlite3PostService{db}
	h = handler{db, postService, messageService}
	os.Exit(m.Run())
}
