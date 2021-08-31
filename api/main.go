package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type handler struct {
	db             *sqlx.DB
	postService    PostService
	messageService MessageService
}

func initDB(db *sqlx.DB) error {
	post_sql, err := ioutil.ReadFile("post_init.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(post_sql))
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	message_sql, err := ioutil.ReadFile("message_init.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(message_sql))
	if err != nil {
		return err
	}

	return nil
}

func main() {

	api := echo.New()
	db, _ := sqlx.Open("sqlite3", ":memory:?_busy_timeout=5000")
	defer db.Close()
	db.SetMaxOpenConns(1)
	db.MustExec("PRAGMA foreign_keys = ON")
	if err := initDB(db); err != nil {
		panic(err)
	}

	postService = Sqlite3PostService{DB: db}
	messageService = Sqlite3MessageService{DB: db}
	h := handler{db, postService, messageService}

	api.Use(middleware.Logger())
	api.Use(middleware.Recover())
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowCredentials, echo.HeaderAuthorization},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	apis := api.Group("/api")
	api.GET("/", h.hello)
	apis.GET("/post/page/:page", h.getPosts)
	apis.POST("/post", h.addPost)
	apis.GET("/post/:id", h.getPost)
	apis.DELETE("/post/:id", h.deletePost)
	apis.GET("/post/counts", h.getPostCounts)

	apis.GET("/post/:id/message", h.getMessages)
	apis.POST("/message", h.addMessage)

	api.Logger.Fatal(api.Start(":8080"))
}

func (h *handler) hello(c echo.Context) error {
	return c.String(http.StatusOK, "hello, world!")
}

type getPostProps struct {
	Id int64 `param:"id" json:"id"`
}
type getPostsProps struct {
	Page int `param:"page" json:"page"`
}

type getMessagesProps struct {
	Id int64 `param:"id" json:"id"`
}
type deletePostProps struct {
	Id int64 `param:"id" json:"id"`
}

func (h *handler) addMessage(c echo.Context) error {
	var props Message
	if err := c.Bind(&props); err != nil {
		fmt.Println(err)
		return c.String(http.StatusNotFound, "parameter error")
	}
	r, err := h.messageService.Add(props)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	p, err := h.postService.Get(Post{Id: r.PostId})
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	if p.LimitCount == 0 || p.LimitCount > p.Count {
		hash := sha512.Sum512([]byte(fmt.Sprintf("%s%d", p.Salt, r.Id)))
		c.SetCookie(&http.Cookie{
			Name:  fmt.Sprintf("post.%d", r.PostId),
			Value: base64.StdEncoding.EncodeToString(hash[:]),
		})
		h.postService.CountIngrement(p)
	}
	return c.JSON(http.StatusCreated, r)
}

func (h *handler) getMessages(c echo.Context) error {
	var props getPostProps
	if err := c.Bind(&props); err != nil {
		return c.String(http.StatusNotFound, "parameter error")
	}
	p := Post{Id: props.Id}

	r, err := h.messageService.Gets(p)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, r)

}
func (h *handler) getPost(c echo.Context) error {
	var props getPostProps
	if err := c.Bind(&props); err != nil {
		fmt.Println(err)
		return c.String(http.StatusNotFound, "parameter error")
	}
	p := Post{Id: props.Id}
	r, err := h.postService.Get(p)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	r.DeletePassword = ""
	if r.LimitTime.Sub(r.UploadedAt) > 0 && time.Now().Sub(r.LimitTime) > 0 {
		r.Url = ""
		fmt.Println("time limit")
	}
	if r.RequireMessage {
		cookie, err := c.Cookie(fmt.Sprintf("post.%d", r.Id))
		if err != nil {
			r.Url = ""
			fmt.Println("cookie not found")
			return c.JSON(http.StatusOK, r)
		}
		messages, err := h.messageService.Gets(p)
		if err != nil {
			r.Url = ""
			fmt.Println("message not found")
			return c.JSON(http.StatusOK, r)
		}
		isOk := false
		for v := range messages {
			hash := sha512.Sum512([]byte(fmt.Sprintf("%s%d", r.Salt, messages[v].Id)))
			fmt.Println(cookie.Value, messages[v].Id)
			if cookie.Value == base64.StdEncoding.EncodeToString(hash[:]) {
				isOk = true
				fmt.Println(cookie.Value, messages[v].Id)
				break
			}
		}
		if !isOk {
			r.Url = ""
			fmt.Println("cookie not match")
		}
		return c.JSON(http.StatusOK, r)
	}
	if r.LimitCount != 0 && r.LimitCount <= r.Count {
		r.Url = ""
		fmt.Println("max count exceed")
	}
	if r, err := h.postService.CountIngrement(r); err != nil {
		fmt.Println(r, err)
	}
	return c.JSON(http.StatusOK, r)
}

func (h *handler) getPosts(c echo.Context) error {
	var props getPostsProps
	if err := c.Bind(&props); err != nil {
		props.Page = 0
	}
	posts, err := h.postService.Gets(props.Page)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	for i := range posts {
		posts[i].Url = ""
		posts[i].DeletePassword = ""
	}
	return c.JSON(http.StatusOK, posts)
}

func (h *handler) addPost(c echo.Context) error {
	var p Post
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "parameter error"})
	}
	now := time.Now()
	if p.UploadedAt.IsZero() {
		p.UploadedAt = now
	}
	if p.LimitTime.IsZero() {
		p.LimitTime = p.UploadedAt
	}
	if err := p.Validate(); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	r, err := h.postService.Add(p)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, r)
}

func (h *handler) getPostCounts(c echo.Context) error {
	count := h.postService.GetCounts()
	return c.JSON(http.StatusOK, map[string]int{"count": count})
}

func (h *handler) deletePost(c echo.Context) error {
	var props deletePostProps
	if err := c.Bind(&props); err != nil {
		fmt.Println(err)
		return c.String(http.StatusNotFound, "parameter error")
	}
	p := Post{Id: props.Id}
	p, err := h.postService.Get(p)
	if err != nil {
		return c.String(http.StatusNoContent, "post not found")
	}
	password := c.Request().Header.Get(echo.HeaderAuthorization)
	if password != p.DeletePassword {
		return c.String(http.StatusUnauthorized, "password mismatch")
	}
	if err := h.postService.Delete(p); err != nil {
		return c.String(http.StatusNoContent, "post not found")
	}

	return c.String(http.StatusOK, "deleted")

}
