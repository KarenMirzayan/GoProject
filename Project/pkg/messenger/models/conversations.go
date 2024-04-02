package models

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Conversations struct {
	ConversationId string `json:"conversation_id"`
	UserId         string `json:"user_id"`
	FriendId       string `json:"friend_id"`
}

type ConversationsModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m ConversationsModel) Insert(conversations *Conversations) error {
	// Insert a new user item into the database.
	query := `
		INSERT INTO user_conversations (conversation_id, user_id, friend_id) 
		VALUES ($1, $2, $3)
		RETURNING conversation_id, user_id, friend_id;
		`
	args := []interface{}{conversations.ConversationId, conversations.UserId, conversations.FriendId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&conversations.ConversationId, &conversations.UserId,
		&conversations.FriendId)
}

func (m ConversationsModel) Get(id int) (*Conversations, error) {
	// Retrieve a specific user item based on its ID.
	query := `
		SELECT conversation_id, user_id, friend_id
		FROM user_conversations
		WHERE conversation_id = $1;
		`
	var conversations Conversations
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&conversations.ConversationId, &conversations.UserId, &conversations.FriendId)
	if err != nil {
		return nil, err
	}
	return &conversations, nil
}

func (m ConversationsModel) Update(conversations *Conversations) error {
	// Update a specific user item in the database.
	query := `
		UPDATE user_conversations
		SET conversation_id = $1, user_id = $2, friend_id = $3
		WHERE user_id = $4
		RETURNING conversation_id, user_id, friend_id;
		`
	args := []interface{}{conversations.ConversationId, conversations.UserId, conversations.FriendId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&conversations.ConversationId, &conversations.UserId, &conversations.FriendId)
}

func (m ConversationsModel) Delete(id int) error {
	// Delete a specific user item from the database.
	query := `
		DELETE FROM user_conversations
		WHERE conversation_id = $1;
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
