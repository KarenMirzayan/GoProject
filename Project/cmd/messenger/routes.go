package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// routes is our main application's router.
func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	// Convert the app.notFoundResponse helper to a http.Handler using the http.HandlerFunc()
	// adapter, and then set it as the custom error handler for 404 Not Found responses.
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// healthcheck
	//r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	v1 := r.PathPrefix("/api/v1").Subrouter()

	//Create conversation
	v1.HandleFunc("/users/{userId:[0-9]+}/conversations", app.createConversationHandler).Methods("POST")
	// Get a conversation
	v1.HandleFunc("/users/{userId:[0-9]+}/conversations/{conversationId:[0-9]+}", app.getConversationHandler).Methods("GET")
	// Delete a specific conversation
	v1.HandleFunc("/users/{userId:[0-9]+}/conversations/{conversationId:[0-9]+}", app.deleteConversationHandler).Methods("DELETE")
	// Get all conversations (with filtering)
	v1.HandleFunc("/users/{userId:[0-9]+}/conversations", app.getConversationsHandler).Methods("GET")

	//Create message in conversation
	v1.HandleFunc("/users/{userId:[0-9]+}/conversations/{conversationId:[0-9]+}/messages", app.createMessageHandler).Methods("POST")
	// Get a specific message
	v1.HandleFunc("/users/{userId:[0-9]+}/conversations/{conversationId:[0-9]+}/messages/{messageId:[0-9]+}", app.getMessageHandler).Methods("GET")
	//Update a specific message
	v1.HandleFunc("/users/{userId:[0-9]+}/conversations/{conversationId:[0-9]+}/messages/{messageId:[0-9]+}", app.updateMessageHandler).Methods("PUT")
	// Delete a specific message
	v1.HandleFunc("/users/{userId:[0-9]+}/conversations/{conversationId:[0-9]+}/messages/{messageId:[0-9]+}", app.deleteMessageHandler).Methods("DELETE")
	// Get all messages of conversation
	v1.HandleFunc("/users/{userId:[0-9]+}/conversations/{conversationId:[0-9]+}/messages", app.getMessagesList).Methods("GET")

	v1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	v1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	v1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")
	// Wrap the router with the panic recovery middleware and rate limit middleware.
	return app.authenticate(r)
}
