package models

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Conversations struct {
	ConversationId int `json:"conversation_id"`
	UserId         int `json:"user_id"`
	FriendId       int `json:"friend_id"`
}

type ConversationsModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m ConversationsModel) GetAll() ([]*Conversations, error) {
	// Retrieve all conversations from the database.
	query := `
        SELECT conversation_id, user_id, friend_id
        FROM user_conversations;
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var conversations []*Conversations
	for rows.Next() {
		var conversation Conversations
		if err := rows.Scan(&conversation.ConversationId, &conversation.UserId, &conversation.FriendId); err != nil {
			return nil, err
		}
		conversations = append(conversations, &conversation)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return conversations, nil
}

func (m ConversationsModel) Insert(conversations *Conversations) error {
	// Insert a new user item into the database.
	query := `
		INSERT INTO user_conversations (user_id, friend_id) 
		VALUES ($1, $2)
		RETURNING conversation_id, user_id, friend_id;
		`
	args := []interface{}{conversations.UserId, conversations.FriendId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&conversations.ConversationId, &conversations.UserId,
		&conversations.FriendId)
}

func (m ConversationsModel) Get(userId, conversationId int) (*Conversations, error) {
	query := `
		SELECT conversation_id, user_id, friend_id
		FROM user_conversations
		WHERE conversation_id = $1 AND (user_id = $2 OR friend_id = $2);
		`
	var conversations Conversations
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, conversationId, userId)
	err := row.Scan(&conversations.ConversationId, &conversations.UserId, &conversations.FriendId)
	if err != nil {
		return nil, err
	}
	return &conversations, nil
}

func (m ConversationsModel) Delete(userCheck, userId, conversationId int) error {
	// Delete a specific user item from the database.
	query := `
		DELETE FROM user_conversations
		WHERE conversation_id = $1 AND (user_id = $3 OR friend_id = $3) AND (user_id = $2 OR friend_id = $2);
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, conversationId, userId, userCheck)
	return err
}

func (m ConversationsModel) GetByUserIDWithPagination(userID int, filters Filters) ([]*Conversations, Metadata, error) {
	// Retrieve conversations specific to the user from the database with pagination
	query := `
        SELECT conversation_id, user_id, friend_id
        FROM user_conversations
        WHERE user_id = $1 OR friend_id = $1
        ORDER BY ` + filters.sortColumn() + ` ` + filters.sortDirection() + `
        LIMIT $2 OFFSET $3;
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	var conversations []*Conversations
	for rows.Next() {
		var conversation Conversations
		if err := rows.Scan(&conversation.ConversationId, &conversation.UserId, &conversation.FriendId); err != nil {
			return nil, Metadata{}, err
		}
		conversations = append(conversations, &conversation)
	}
	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Retrieve total number of records for the user
	var totalRecords int
	err = m.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_conversations WHERE user_id = $1 OR friend_id = $1;",
		userID).Scan(&totalRecords)
	if err != nil {
		return nil, Metadata{}, err
	}

	// Calculate metadata
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return conversations, metadata, nil
}
