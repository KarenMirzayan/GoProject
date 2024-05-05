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

func (app *application) createChannelHandler(w http.ResponseWriter, r *http.Request) {
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
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	channel := &models.Channels{
		UserId: userID,
		Name:   input.Name,
	}

	if err := app.models.Channels.Insert(channel); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Write the created conversation as JSON response.
	app.writeJSON(w, http.StatusCreated, envelope{"channel": channel}, nil)
}

func (app *application) getChannelHandler(w http.ResponseWriter, r *http.Request) {
	// Extract userId and conversationId parameters from the request URL.
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["userId"])
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid user ID")
		return
	}
	channelID, err := strconv.Atoi(params["channelId"])
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid conversation ID")
		return
	}

	// Retrieve the conversation from the database.
	channel, err := app.models.Channels.Get(userID, channelID)
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
	app.writeJSON(w, http.StatusOK, envelope{"channel": channel}, nil)
}

func (app *application) deleteChannelHandler(w http.ResponseWriter, r *http.Request) {
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

	channelID, err := strconv.Atoi(params["channelId"])
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid conversation ID")
		return
	}

	// Delete the conversation from the database.
	err = app.models.Channels.Delete(userID, channelID)
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

func (app *application) updateChannelHandler(w http.ResponseWriter, r *http.Request) {
	// Extract parameters from the request URL
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["userId"])
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid user ID")
		return
	}
	channelID, err := strconv.Atoi(params["channelId"])
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid channel ID")
		return
	}

	// Query the database for the channel using its IDs
	channel, err := app.models.Channels.Get(userID, channelID)
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
		Name *string `json:"name"`
	}

	// Read JSON input into struct
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Update channel fields if they're provided in the input
	if input.Name != nil {
		channel.Name = *input.Name
	}

	// Validate the updated channel
	v := validator.New()
	if models.ValidateChannel(v, channel); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Update the channel in the database
	err = app.models.Channels.Update(channel)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Respond with the JSON representation of the updated channel
	app.writeJSON(w, http.StatusOK, envelope{"channel": channel}, nil)
}

func (app *application) getChannelsList(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from request parameters
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["userId"])
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get all channels for the given user ID
	channels, err := app.models.Channels.GetAll(userID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Respond with JSON containing channels
	app.writeJSON(w, http.StatusOK, envelope{"channels": channels}, nil)
}
