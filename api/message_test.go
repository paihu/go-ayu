package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestGetMessages(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("GET", "/post/1/message", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/post/:id/message")
	c.SetParamNames("id")
	c.SetParamValues("1")
	if assert.NoError(t, h.getMessages(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var m []*Message
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &m)) {
			assert.Equal(t, len(messages), len(m))
		}
	}
}
