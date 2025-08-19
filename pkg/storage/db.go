package storage

import (
	"BlogAPI/pkg/jwtutils"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)



type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Printf("Connected to the DB!") //Just to be confident 
	return &PostgresStorage{db: db}, nil
}

func (p *PostgresStorage) Register(newUser User) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password),bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	var exists bool
	err = p.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", newUser.Email).Scan(&exists)
	if exists {
		return errors.New("email already exists")
	}
	if err != nil {
		return fmt.Errorf("error: %W", err)
	}
	
	_, err = p.db.Exec("INSERT into users (username, email, password) VALUES ($1,$2,$3)",newUser.Username,newUser.Email,hashedPass)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStorage) Login(email, password string) (string, error) {
	var (
		pass string 
		userID int
	)
	err := p.db.QueryRow("SELECT id,password FROM users WHERE email = $1",email).Scan(&userID,&pass)
	if err == sql.ErrNoRows {
		return "", errors.New("user not found")
	}

	isCorrectPassword := bcrypt.CompareHashAndPassword([]byte(pass),[]byte(password))
	if isCorrectPassword != nil {
		return "", errors.New("wrong password")
	}

	secret := os.Getenv("JWT_SECRET")

	token, err := jwtutils.GenerateToken(userID,secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (p *PostgresStorage) GetUserById(id int) (*User, error) {
	var user User 
	err := p.db.QueryRow(
		"SELECT id, username, email FROM users WHERE id = $1",
		id,
	).Scan(&user.ID,&user.Username,&user.Email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *PostgresStorage) UpdateProfile(id int,username, email, hashedPass string) error {
	_, err := p.db.Exec(`
		UPDATE users 
		SET username = COALESCE($1, username),
            email = COALESCE($2, email),
            password_hash = COALESCE($3, password_hash)
        WHERE id = $4`,username,email,hashedPass,id,
	)
	return err
}

func (p *PostgresStorage) CreatePost(info Post) error {
	_, err := p.db.Exec(`
		INSERT into posts (id,title,content,authorId,createdAt,updatedAt) VALUES ($1,$2,$3,$4,$5,$6) 
	`,info.ID,info.Title,info.Content,info.AuthorId,info.Createdat,info.Updatedat)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStorage) GetPosts(limit int) (*[]Post, error) {
	posts := []Post{}
	rows, err := p.db.Query(`
		SELECT * FROM posts LIMIT $1
	`, limit)

	if rows != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Post{}
		err := rows.Scan(&p.ID,&p.Title,&p.Content,&p.AuthorId,&p.Createdat,&p.Updatedat)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)

	}
	return &posts, nil
}

func (p *PostgresStorage) GetPostById(ID int) (*Post, error) {
	post := &Post{}
	err := p.db.QueryRow(`
		SELECT * FROM posts WHERE id = $1
	`, ID).Scan(&post.ID,&post.Title,&post.Content,&post.AuthorId,&post.Createdat,&post.Updatedat)
	return post, err
}

func (p *PostgresStorage) EditPost(postID int, authorID int, title *string, content *string) error {
	var dbAuthorID int 
	err := p.db.QueryRow(`
		SELECT authorId from posts WHERE id = $1
	`,postID).Scan(&dbAuthorID)

	if err != nil {
		return err
	}

	if dbAuthorID != authorID {
		return errors.New("you cant edit only your own posts")
	}

	 _, err = p.db.Exec(`
        UPDATE posts 
        SET title = COALESCE($1, title),
            content = COALESCE($2, content),
            updated_at = NOW()
        WHERE id = $3`,
        title, content, postID,
    )
    
    return err
}

func (p *PostgresStorage) DeletePost(ID int) error {
	_,err := p.db.Exec(`
		DELETE FROM posts WHERE id = $1
	`,ID)
	return err
}
