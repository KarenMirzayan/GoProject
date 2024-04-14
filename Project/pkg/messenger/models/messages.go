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
	SenderId       int    `json:"sender_id"`
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
		VALUES ($1, $2, $3, $4)
		RETURNING message_id, conversation_id, sender_id, content, timestamp;
		`
	args := []interface{}{messages.ConversationId, messages.SenderId, messages.Content, messages.Timestamp}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&messages.MessageId, &messages.ConversationId, &messages.SenderId, &messages.Content, &messages.Timestamp)
}

func (m MessagesModel) Get(conversationID, senderID, messageID string) (*Messages, error) {
	query := `
		SELECT m.message_id, m.conversation_id, m.sender_id, m.content, m.timestamp
		FROM messages m
		INNER JOIN conversations c ON m.conversation_id = c.conversation_id
		WHERE m.conversation_id = $1 AND m.message_id = $2 AND (c.user_id = $3 OR c.friend_id = $3);
	`
	var messages Messages
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, conversationID, senderID, messageID)
	err := row.Scan(&messages.MessageId, &messages.ConversationId, &messages.SenderId, &messages.Content, &messages.Timestamp)
	if err != nil {
		return nil, err
	}
	return &messages, nil
}

func (m MessagesModel) Update(messages *Messages) error {
	// Update a specific menu item in the database.
	query := `
		UPDATE messages m
		SET content = $1
		FROM conversations c
		WHERE m.conversation_id = c.conversation_id 
		AND m.conversation_id = $2 
		AND m.message_id = $3
		AND (c.user_id = $4 OR c.friend_id = $4)
		RETURNING m.message_id, m.conversation_id, m.sender_id, m.content, m.timestamp
	`
	args := []interface{}{messages.Content, messages.ConversationId, messages.SenderId, messages.MessageId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, args...)
	err := row.Scan(&messages.MessageId, &messages.ConversationId, &messages.SenderId, &messages.Content, &messages.Timestamp)
	if err != nil {
		return err
	}
	return nil
}

func (m MessagesModel) Delete(conversationID, senderID, messageID string) error {
	// Delete a specific menu item from the database.
	query := `
		DELETE FROM messages
		WHERE conversation_id = $1 AND sender_id = $2 AND (c.user_id = $3 OR c.friend_id = $3);
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, conversationID, senderID, messageID)
	return err
}

func ValidateMessage(v *validator.Validator, message *Messages) {
	// Check if the content field is empty.
	v.Check(message.Content != "", "content", "must be provided")
	// Add additional validation rules as needed.
}

func (m MessagesModel) GetAll(userId, conversationId int, query string, filters Filters) ([]*Messages, Metadata, error) {
	// Construct the SQL query for retrieving messages with content containing the query string.
	// We use the ILIKE operator for case-insensitive search.
	// We also use the OFFSET and LIMIT clauses for pagination.
	sqlQuery := fmt.Sprintf(`
        SELECT count(*) OVER(), m.message_id, m.conversation_id, m.sender_id, m.content, m.timestamp
        FROM messages m
        INNER JOIN user_conversations c ON m.conversation_id = c.conversation_id
        WHERE m.content ILIKE '%%%s%%' AND m.conversation_id = $1
        AND (c.user_id = $2 OR c.friend_id = $2)
        ORDER BY %s %s
        LIMIT $3 OFFSET $4`, query, filters.sortColumn(), filters.sortDirection())

	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query and retrieve the result set.
	rows, err := m.DB.QueryContext(ctx, sqlQuery, conversationId, userId, filters.limit(), filters.offset())
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
