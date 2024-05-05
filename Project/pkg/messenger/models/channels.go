package models

import (
	"context"
	"database/sql"
	"github.com/KarenMirzayan/Project/pkg/messenger/validator"
	"log"
	"time"
)

type Channels struct {
	ChannelId int    `json:"channel_id"`
	UserId    int    `json:"user_id"`
	Name      string `json:"name"`
}

type ChannelsModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (c ChannelsModel) Insert(channels *Channels) error {
	// Insert a new user item into the database.
	query := `
		INSERT INTO channels (user_id, name) 
		VALUES ($1, $2)
		RETURNING channel_id, user_id, name;
		`
	args := []interface{}{channels.UserId, channels.Name}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&channels.ChannelId, &channels.UserId,
		&channels.Name)
}

func (c ChannelsModel) Get(userId, channelId int) (*Channels, error) {
	query := `
		SELECT channel_id, user_id, name
		FROM channels
		WHERE channel_id = $1 AND user_id = $2;
		`
	var channels Channels
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := c.DB.QueryRowContext(ctx, query, channelId, userId)
	err := row.Scan(&channels.ChannelId, &channels.UserId, &channels.Name)
	if err != nil {
		return nil, err
	}
	return &channels, nil
}

func (c ChannelsModel) Delete(userId, channelId int) error {
	// Delete a specific user item from the database.
	query := `
		DELETE FROM channels
		WHERE channel_id = $1 AND user_id = $2;
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.DB.ExecContext(ctx, query, channelId, userId)
	return err
}

func (c ChannelsModel) Update(channel *Channels) error {
	// Update a specific channel in the database.
	query := `
		UPDATE channels
		SET name = $1
		WHERE channel_id = $2
		RETURNING channel_id, user_id, name
	`
	args := []interface{}{channel.Name, channel.ChannelId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := c.DB.QueryRowContext(ctx, query, args...)
	err := row.Scan(&channel.ChannelId, &channel.UserId, &channel.Name)
	if err != nil {
		return err
	}
	return nil
}

func ValidateChannel(v *validator.Validator, channel *Channels) {
	v.Check(channel.Name != "", "name", "must be provided")
}

func (c ChannelsModel) GetAll(userId int) ([]*Channels, error) {
	// Construct the SQL query for retrieving all channels of a specific user.
	sqlQuery := `
        SELECT channel_id, user_id, name
        FROM channels
        WHERE user_id = $1
    `

	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query and retrieve the result set.
	rows, err := c.DB.QueryContext(ctx, sqlQuery, userId)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			c.ErrorLog.Println(err)
		}
	}()

	// Declare a variable to store channels.
	var channels []*Channels

	// Iterate over the result set and scan each row into a Channels struct.
	for rows.Next() {
		var channel Channels
		if err := rows.Scan(&channel.ChannelId, &channel.UserId, &channel.Name); err != nil {
			return nil, err
		}
		channels = append(channels, &channel)
	}

	// Check for errors during iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return the channels.
	return channels, nil
}
