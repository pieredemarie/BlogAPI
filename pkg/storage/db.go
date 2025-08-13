package storage

import (
	"BlogAPI/pkg/jwtutils"
	"database/sql"
	"errors"
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
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password),bcrypt.DefaultCost)

	var exists bool
	err := p.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", newUser.Email).Scan(&exists)
	if exists {
		return errors.New("email already exists")
	}
	
	_, err1 := p.db.Exec("INSERT into users (username, email, password) VALUES ($1,$2,$3)",newUser.Username,newUser.Email,hashedPass)
	if err != nil {
		return err1
	}
	return nil
}

func (p *PostgresStorage) Login(email, password string) (string, error) {
	var (
		pass string 
		userID int
	)
	err := p.db.QueryRow("SELECT (id,password) FROM users WHERE email = $1",email).Scan(&userID,&pass)
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