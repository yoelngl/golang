package model

import (
	"database/sql"
	"restfulapi/config"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,lte=50"`
	Password string `json:"password" validate:"required,gte=7"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,lte=50"`
	Email    string `json:"email" validate:"required,email,lte=50"`
	Password string `json:"password" validate:"required,gte=7"`
	Image    string
}

var message string

const (
	registerUser = `INSERT INTO users (uuid,name, email, password, image) VALUES (?,?,?,?,?)`

	checkUser = `SELECT email,password FROM users where email = ?`
)

func RegisterStore(user RegisterRequest) (string, error) {
	db, err := config.Connect()
	if err != nil {
		message = "Failed to connect the Database"
		return message, err
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		message = "Failed to Begin sql!"
		return message, err
	}

	stmt, err := db.Prepare(registerUser)
	if err != nil {
		tx.Rollback()
		message = "Failed to prepare sql!"
		return message, err
	}

	defer stmt.Close()

	uuid := uuid.New()
	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	_, err = stmt.Exec(uuid, user.Name, user.Email, password, user.Image)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			message = "Email has already exists!"
			return message, err
		}
		tx.Rollback()
		message = "Failed to register user!"
		return message, err
	}

	err = tx.Commit()
	if err != nil {
		message = "Failed to commit transaction"
		return message, err
	}

	message = "Successfully registered user!"
	return message, nil
}

func CheckAccount(user LoginRequest) (string, error) {
	db, err := config.Connect()
	if err != nil {
		message = "Failed to connect the Database"
		return message, err
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		message = "Failed to Begin sql!"
		return message, err
	}

	stmt, err := db.Prepare(checkUser)
	if err != nil {
		tx.Rollback()
		message = "Failed to prepare sql!"
		return message, err
	}

	defer stmt.Close()

	var email string
	var password string
	err = stmt.QueryRow(user.Email).Scan(&email, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			message = "No user found with this account!"
		} else {
			message = "Failed to check email!"
		}

		return message, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if err != nil {
		message = "Incorrect Password, please try again!"
		return message, err
	}

	message = "success"
	return message, nil
}


