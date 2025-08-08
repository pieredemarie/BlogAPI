package storage

import (
	"database/sql"
	"errors"
	"log"

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
	var pass string
	err := p.db.QueryRow("SELECT (password) FROM users WHERE email = $1",email).Scan(&pass)
	if err == sql.ErrNoRows {
		return "", errors.New("user not found")
	}

	isCorrectPassword := bcrypt.CompareHashAndPassword([]byte(pass),[]byte(password))
	if isCorrectPassword != nil {
		return "", errors.New("wrong password")
	}

	// here wil be JWT logic but i don't want it right now
	// just because i'm lazy bruh
	// maaan i got tired after writing ten code lines
	// it pissed me off every time i open vs code
	//TODO: 
	///сделать JWT в отдельном пакете и сюда присобачить
	//желательно сегодня.
	return "dummy-token", nil
}
