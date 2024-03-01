package models

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Users struct {
	UserId      string `json:"user_id"`
	First       string `json:"firstname"`
	Last        string `json:"lastname"`
	DateOfBirth string `json:"date_of_birth"`
	Login       string `json:"login"`
	Password    string `json:"password"`
}

//create table if not exists users (
//user_id serial primary key,
//firstname text not null,
//lastname text not null,
//date_of_birth timestamp(0) with time zone not null,
//login varchar(16) not null,
//password varchar(16) not null
//);

type UsersModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m UsersModel) Insert(users *Users) error {
	// Insert a new menu item into the database.
	query := `
		INSERT INTO users (firstname, lastname, date_of_birth, login, password) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING firstname, lastname, date_of_birth, login, password;
		`
	args := []interface{}{users.First, users.Last, users.DateOfBirth, users.Login, users.Password}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&users.First, &users.Last, &users.DateOfBirth, &users.Login, &users.Password)
}

func (m UsersModel) Get(id int) (*Users, error) {
	// Retrieve a specific menu item based on its ID.
	query := `
		SELECT user_id, firstname, lastname, date_of_birth, login, password
		FROM users
		WHERE user_id = $1;
		`
	var users Users
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&users.UserId, &users.First, &users.Last, &users.DateOfBirth, &users.Login, &users.Password)
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (m UsersModel) Update(users *Users) error {
	// Update a specific menu item in the database.
	query := `
		UPDATE users
		SET firstname = $1, lastname = $2, date_of_birth = $3, login = $4, password = $5
		WHERE user_id = $6
		RETURNING firstname, lastname, date_of_birth, login, password;
		`
	args := []interface{}{users.First, users.Last, users.DateOfBirth, users.Login, users.Password, users.UserId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&users.First, &users.Last, &users.DateOfBirth, &users.Login, &users.Password)
}

func (m UsersModel) Delete(id int) error {
	// Delete a specific menu item from the database.
	query := `
		DELETE FROM users
		WHERE users.user_id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
