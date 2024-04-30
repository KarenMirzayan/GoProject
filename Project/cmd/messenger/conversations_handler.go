package main

import (
	"encoding/json"
	"errors"
	"github.com/KarenMirzayan/Project/pkg/messenger/validator"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/KarenMirzayan/Project/pkg/messenger/models"
)

func (app *application) getConversationsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract userID from URL parameters
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["userId"])
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid user ID")
		return
	}
	v := validator.New()
	qs := r.URL.Query()
	// Extract filters from query parameters
	filters := models.Filters{
		Page:         app.readInt(qs, "page", 1, v),
		PageSize:     app.readInt(qs, "page_size", 10, v),
		Sort:         app.readStrings(qs, "sort", "conversation_id"),
		SortSafeList: []string{"conversation_id", "-conversation_id"},
	}

	// Retrieve conversations from the database with pagination
	conversations, metadata, err := app.models.Conversations.GetByUserIDWithPagination(userID, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Write the conversations and metadata as JSON response
	app.writeJSON(w, http.StatusOK, envelope{"conversations": conversations, "metadata": metadata}, nil)
}

func (app *application) createConversationHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := app.contextGetUser(r)
	userID, err := strconv.Atoi(params["userId"])
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if int(user.ID) != userID {
		app.errorResponse(w, r, http.StatusUnauthorized, "Wrong token")
	}

	var input struct {
		FriendId int `json:"friend_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	conversation := &models.Conversations{
		UserId:   userID,
		FriendId: input.FriendId,
	}

	if err := app.models.Conversations.Insert(conversation); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Write the created conversation as JSON response.
	app.writeJSON(w, http.StatusCreated, envelope{"conversation": conversation}, nil)
}

func (app *application) getConversationHandler(w http.ResponseWriter, r *http.Request) {
	// Extract userId and conversationId parameters from the request URL.
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

	// Retrieve the conversation from the database.
	conversation, err := app.models.Conversations.Get(userID, conversationID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Write the conversation as JSON response.
	app.writeJSON(w, http.StatusOK, envelope{"conversation": conversation}, nil)
}

func (app *application) deleteConversationHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := app.contextGetUser(r)
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
	// Delete the conversation from the database.
	err = app.models.Conversations.Delete(int(user.ID), userID, conversationID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Write success message as JSON response.
	app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
}
