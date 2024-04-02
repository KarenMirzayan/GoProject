package models

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Messages struct {
	MessageId      string `json:"message_id"`
	ConversationId string `json:"conversation_id"`
	SenderId       string `json:"sender_id"`
	Content        string `json:"content"`
	Timestamp      string `json:"timestamp"`
}

type MessagesModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m MessagesModel) Insert(messages *Messages) error {
	// Insert a new menu item into the database.
	query := `
		INSERT INTO messages (conversation_id, sender_id, content, timestamp) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING message_id, conversation_id, sender_id,content, timestamp;
		`
	args := []interface{}{messages.ConversationId, messages.SenderId, messages.Content, messages.Timestamp}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&messages.ConversationId, &messages.SenderId, &messages.Content, &messages.Timestamp)
}

func (m MessagesModel) Get(id int) (*Messages, error) {
	// Retrieve a specific menu item based on its ID.
	query := `
		SELECT message_id, conversation_id, sender_id, content, timestamp
		FROM messages
		WHERE message_id = $1;
		`
	var messages Messages
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&messages.MessageId, &messages.ConversationId, &messages.SenderId, &messages.Content, &messages.Timestamp)
	if err != nil {
		return nil, err
	}
	return &messages, nil
}

func (m MessagesModel) Update(messages *Messages) error {
	// Update a specific menu item in the database.
	query := `
		UPDATE messages
		SET conversation_id = $1, sender_id = $2, content = $3, timestamp = $4
		WHERE message_id = $6
		RETURNING conversation_id, sender_id, content, timestamp
		`
	args := []interface{}{messages.MessageId, messages.ConversationId, messages.SenderId, messages.Content, messages.Timestamp}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&messages.ConversationId, &messages.SenderId, &messages.Content, &messages.Timestamp)
}

func (m MessagesModel) Delete(id int) error {
	// Delete a specific menu item from the database.
	query := `
		DELETE FROM messages
		WHERE messages.message_id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
