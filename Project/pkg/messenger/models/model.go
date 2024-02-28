package models

import (
	"database/sql"
	"log"
	"os"
)

type Models struct {
	Users UsersModel
	//UserConversations UserConversationsModel
	//Messages MessagesModel
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
		//Restaurants: RestaurantModel{
		//	DB:       db,
		//	InfoLog:  infoLog,
		//	ErrorLog: errorLog,
		//},
	}
}
