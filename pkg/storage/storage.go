package storage

import "time"

type Post struct {
	ID int `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	AuthorId int `json:"authorID"`
	Createdat time.Time `json:"createdAt"`
	Updatedat time.Time `json:"updatedAt"`
}

type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type AuthStorage interface {
	Login(email, password string) (string, error)
	Register(data User) error

}

//
type PostStorage interface {
	GetPosts(limit int) ([]Post, error)
	GetPostById(ID int) (Post, error)
	CreatePost(info Post) (error)
	EditPost(ID int, newData Post) (error)
	DeletePost(ID int) (error)
}
