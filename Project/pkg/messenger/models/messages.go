package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/KarenMirzayan/Project/pkg/messenger/validator"
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

func ValidateMessage(v *validator.Validator, message *Messages) {
	// Check if the content field is empty.
	v.Check(message.Content != "", "content", "must be provided")
	// Add additional validation rules as needed.
}

func (m MessagesModel) GetAll(query string, from, to int, filters Filters) ([]*Messages, Metadata, error) {
	// Construct the SQL query for retrieving messages with content containing the query string.
	// We use the ILIKE operator for case-insensitive search.
	// We also use the OFFSET and LIMIT clauses for pagination.
	sqlQuery := fmt.Sprintf(`
        SELECT count(*) OVER(), message_id, conversation_id, sender_id, content, timestamp
        FROM messages
        WHERE content ILIKE '%%%s%%'
        ORDER BY %s %s
        LIMIT $1 OFFSET $2`, query, filters.sortColumn(), filters.sortDirection())

	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query and retrieve the result set.
	rows, err := m.DB.QueryContext(ctx, sqlQuery, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			m.ErrorLog.Println(err)
		}
	}()

	// Declare variables to store total records and messages.
	totalRecords := 0
	var messages []*Messages

	// Iterate over the result set and scan each row into a Message struct.
	for rows.Next() {
		var message Messages
		if err := rows.Scan(&totalRecords, &message.MessageId, &message.ConversationId, &message.SenderId, &message.Content, &message.Timestamp); err != nil {
			return nil, Metadata{}, err
		}
		messages = append(messages, &message)
	}

	// Check for errors during iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate metadata based on total records and pagination parameters.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// Return the messages and metadata.
	return messages, metadata, nil
}
