package models

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

var (
	// ErrRecordNotFound is returned when a record doesn't exist in database.
	ErrRecordNotFound = errors.New("record not found")

	// ErrEditConflict is returned when a there is a data race, and we have an edit conflict.
	ErrEditConflict = errors.New("edit conflict")
)

type Models struct {
	Users         UsersModel
	Conversations ConversationsModel
	Messages      MessagesModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Users: UsersModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Conversations: ConversationsModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Messages: MessagesModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
	}
}
