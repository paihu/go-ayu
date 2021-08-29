package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestGetPost(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("GET", "/post/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/post/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	if assert.NoError(t, h.getPost(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, postJson, rec.Body.String())
	}
}

func TestAddPostWithInvalidUrl(t *testing.T) {
	e := echo.New()
	p := Post{
		User:           "hoge",
		Email:          "hoge@hoge.com",
		Title:          "title",
		DeletePassword: "fuga",
		Url:            "https://localhost/hoge.zip",
		Kind:           "other",
		UploadedAt:     time.Now().Add(test_duration),
	}
	j, err := json.Marshal(p)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest("POST", "/post", bytes.NewReader(j))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	req.Header.Set("Content-type", "application/json")
	id := struct {
		Id int64 `db:"seq"`
	}{}
	_ = db.Get(&id, `select seq from sqlite_sequence where name="post"`)
	if assert.NoError(t, h.addPost(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "{\"error\":\"url validation error\"}\n", rec.Body.String())
	}
}
func TestAddPostWithInvalidEmail(t *testing.T) {
	e := echo.New()
	p := Post{
		User:           "hoge",
		Email:          "hoge.hoge.com",
		Title:          "title",
		DeletePassword: "fuga",
		Url:            "https://localhost.com/hoge.zip",
		Kind:           "other",
		UploadedAt:     time.Now().Add(test_duration),
	}
	j, err := json.Marshal(p)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest("POST", "/post", bytes.NewReader(j))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	req.Header.Set("Content-type", "application/json")
	id := struct {
		Id int64 `db:"seq"`
	}{}
	_ = db.Get(&id, `select seq from sqlite_sequence where name="post"`)
	if assert.NoError(t, h.addPost(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "{\"error\":\"email validation error\"}\n", rec.Body.String())
	}
}
func TestAddPost(t *testing.T) {
	e := echo.New()
	p := Post{
		User:           "hoge",
		Email:          "hoge@hoge.com",
		Title:          "title",
		DeletePassword: "fuga",
		Url:            "https://localhost.com/hoge.zip",
		Kind:           "other",
		UploadedAt:     time.Now().Add(test_duration),
	}
	j, err := json.Marshal(p)
	assert.Equal(t, err, nil)

	req := httptest.NewRequest("POST", "/post", bytes.NewReader(j))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	req.Header.Set("Content-type", "application/json")
	if assert.NoError(t, h.addPost(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		id := struct {
			Id int64 `db:"seq"`
		}{}
		_ = db.Get(&id, `select seq from sqlite_sequence where name="post"`)
		p.Id = id.Id
		p.LimitTime = p.UploadedAt
		j, _ = json.Marshal(p)
		assert.Equal(t, string(j)+"\n", rec.Body.String())
	}

	// check inster data
	j, _ = json.Marshal(p)
	req = httptest.NewRequest("GET", fmt.Sprintf("/post/%d", p.Id), nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/post/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", p.Id))
	if assert.NoError(t, h.getPost(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(j)+"\n", rec.Body.String())
	}
}
