package main

import (
	"errors"
	"github.com/KarenMirzayan/Project/pkg/messenger/models"
	"github.com/KarenMirzayan/Project/pkg/messenger/validator"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

import "time"

func (app *application) createMessageHandler(w http.ResponseWriter, r *http.Request) { //TODO: get sender_id from url
	params := mux.Vars(r)
	conversationID := params["conversationId"]
	userID, err := strconv.Atoi(params["userId"])

	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Define a struct to hold JSON input data
	var input struct {
		SenderID int    `json:"sender_id"`
		Content  string `json:"content"`
	}

	// Read JSON input into the struct
	err = app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Extract conversation_id from the request
	if conversationID == "" {
		app.errorResponse(w, r, http.StatusBadRequest, "conversation_id is required")
		return
	}

	// Generate timestamp
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Create a new Message instance with the input data, conversation_id, and generated timestamp
	message := &models.Messages{
		ConversationId: conversationID,
		SenderId:       userID,
		Content:        input.Content,
		Timestamp:      timestamp,
	}

	// Insert the new message into the database
	err = app.models.Messages.Insert(message)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Respond with the JSON representation of the newly created message
	app.writeJSON(w, http.StatusCreated, envelope{"message": message}, nil)
}

func (app *application) getMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract parameters from the request URL
	params := mux.Vars(r)
	userID := params["userId"]
	conversationID := params["conversationId"]
	messageID := params["messageId"]

	// Query the database for the message using its IDs
	message, err := app.models.Messages.Get(conversationID, userID, messageID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Respond with the JSON representation of the message
	app.writeJSON(w, http.StatusOK, envelope{"message": message}, nil)
}

func (app *application) updateMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract parameters from the request URL
	params := mux.Vars(r)
	userID := params["userId"]
	conversationID := params["conversationId"]
	messageID := params["messageId"]

	// Query the database for the message using its IDs
	message, err := app.models.Messages.Get(conversationID, userID, messageID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Define struct to hold JSON input data
	var input struct {
		Content *string `json:"content"`
	}

	// Read JSON input into struct
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Update message fields if they're provided in the input
	if input.Content != nil {
		message.Content = *input.Content
	}

	// Validate the updated message
	v := validator.New()
	if models.ValidateMessage(v, message); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Update the message in the database
	err = app.models.Messages.Update(message)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Respond with the JSON representation of the updated message
	app.writeJSON(w, http.StatusOK, envelope{"message": message}, nil)
}

func (app *application) deleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract parameters from the request URL
	params := mux.Vars(r)
	userID := params["userId"]
	conversationID := params["conversationId"]
	messageID := params["messageId"]

	// Delete the message from the database
	err := app.models.Messages.Delete(conversationID, userID, messageID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Respond with success message
	app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
}

func (app *application) getMessagesList(w http.ResponseWriter, r *http.Request) {
	// Define a struct to hold input parameters and filters
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["userId"])
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid user ID")
		return
	}
	conversationID, err := strconv.Atoi(params["conversationId"])
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid conversation ID")
		return
	}
	var input struct {
		Query          string
		models.Filters // Embedding Filters struct for pagination and sorting
	}
	v := validator.New()
	qs := r.URL.Query()

	// Extract query string values and set defaults if not provided
	input.Query = app.readStrings(qs, "query", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readStrings(qs, "sort", "timestamp")
	input.Filters.SortSafeList = []string{
		"-timestamp", "timestamp",
	}

	// Validate input parameters and filters
	if models.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Get messages based on input parameters and filters
	messages, metadata, err := app.models.Messages.GetAll(userID, conversationID, input.Query, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Respond with JSON containing messages and metadata
	app.writeJSON(w, http.StatusOK, envelope{"messages": messages, "metadata": metadata}, nil)
}
