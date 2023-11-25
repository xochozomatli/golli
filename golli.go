package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type AbstractRepository[T any] interface {
	Add(obj T)
	Get(id uint64) T
}

type UserRepository struct {
	Users []User
}

func (repo UserRepository) Add(user User) {
	repo.Users = append(repo.Users, user)
}

func (repo UserRepository) Get(id uint64) *User {
	for _, u := range repo.Users {
		if u.ID == id {
			return &u
		}
	}
	return nil
}

type LessonRepository struct {
	Lessons []Lesson
}

func (repo LessonRepository) Add(lesson Lesson) {
	repo.Lessons = append(repo.Lessons, lesson)
}
func (repo LessonRepository) Get(id uint64) *Lesson {
	for _, l := range repo.Lessons {
		if l.ID == id {
			return &l
		}
	}
	return nil
}

type TolliRepository struct {
	Tolli []Tolli
}

func (repo TolliRepository) Add(tolli Tolli) {
	repo.Tolli = append(repo.Tolli, tolli)
}

func (repo TolliRepository) Get(id uint64) *Tolli {
	for _, t := range repo.Tolli {
		if t.ID == id {
			return &t
		}
	}
	return nil
}

type User struct {
	ID         uint64
	Username   string `json:"username"`
	PassHash   uint64
	Email      string `json:"email"`
	Authored   []uint64
	SignupDate time.Time
}

type Lesson struct {
	ID        uint64
	Title     string `json:"title"`
	Content   string `json:"content"`
	AuthorID  uint32 `json:"author_id"`
	Language  string `json:"language"`
	TimeStamp time.Time
}

type Tolli struct {
	ID        uint64
	UserID    uint64 `json:"user_id"`
	Word      string `json:"word"`
	State     byte   `json:"state"`
	TimeStamp time.Time
}

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])

	if r.Method == "GET" {
		ptr, err := getUser(r)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		user := *ptr
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if r.Method == "POST" {
		ptr, err := createUser(r)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		user := *ptr
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
}

func getUser(r *http.Request) (*User, error) {
}

func createUser(r *http.Request) (*User, error) {
	var data map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	for k, v := range data {
		if v == nil {
			msg := fmt.Sprintf("nil value for key %s", k)
			return nil, errors.New(msg)
		}
	}
	user := User{}
	err = decoder.Decode(&user)
	if err != nil {
		return nil, errors.New("Couldn't decode request body into struct User")
	}
	user.ID = hash(data["email"].(string))
	user.PassHash = hash(data["password"].(string))
	user.SignupDate = time.Now().UTC()

	return &user, nil
}

func lessonsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ptr, err := createLesson(r)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		lesson := *ptr
		err = json.NewEncoder(w).Encode(lesson)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

}

func createLesson(r *http.Request) (*Lesson, error) {
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, errors.New("Couldn't decode request body into data map")
	}
	for k, v := range data {
		if v == nil {
			msg := fmt.Sprintf("nil value for key %s", k)
			return nil, errors.New(msg)
		}
	}
	lesson := Lesson{}
	err = json.NewDecoder(r.Body).Decode(lesson)
	if err != nil {
		return nil, errors.New("Couldn't decode request body into struct Lesson")
	}
	lesson.ID = hash(data["content"].(string))
	lesson.TimeStamp = time.Now().UTC()

	return &lesson, nil
}

func tolliHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ptr, err := createTolli(r)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		tolli := *ptr
		err = json.NewEncoder(w).Encode(tolli)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

}

func createTolli(r *http.Request) (*Tolli, error) {
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, errors.New("Couldn't decode request body into data map")
	}
	for k, v := range data {
		if v == nil {
			msg := fmt.Sprintf("nil value for key %s", k)
			return nil, errors.New(msg)
		}
	}
	tolli := Tolli{}
	err = json.NewDecoder(r.Body).Decode(tolli)
	if err != nil {
		return nil, errors.New("Couldn't decode request body into tolli struct")
	}
	tolli.TimeStamp = time.Now().UTC()
	body, _ := io.ReadAll(r.Body)
	tolli.ID = hash(string(body) + string(tolli.TimeStamp.Unix()))

	return &tolli, nil

}

func main() {
	var connStr = "postgres://postgres:secret@localhost/postgres?sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/users/", usersHandler)
	http.HandleFunc("/lessons/", lessonsHandler)
	http.HandleFunc("/tolli/", tolliHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
