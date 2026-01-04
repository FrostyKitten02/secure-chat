package repo

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"secure-chat/repo/model"
	"time"
)

func FindDirectMessagesBetweenUsers(
	ctx context.Context,
	user1 string,
	user2 string,
) ([]model.DirectMessage, error) {

	sql, args, err := psql.
		Select(
			"id",
			"sender_id",
			"receiver_id",
			"cipher_text",
			"nonce",
			"sender_identity_id",
			"receiver_identity_id",
			"created_at",
		).
		From("direct_message").
		Where(
			sq.Or{
				sq.And{
					sq.Eq{"sender_id": user1},
					sq.Eq{"receiver_id": user2},
				},
				sq.And{
					sq.Eq{"sender_id": user2},
					sq.Eq{"receiver_id": user1},
				},
			},
		).
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]model.DirectMessage, 0)

	for rows.Next() {
		var m model.DirectMessage
		err := rows.Scan(
			&m.ID,
			&m.SenderID,
			&m.ReceiverID,
			&m.CipherText,
			&m.Nonce,
			&m.SenderIdentityID,
			&m.ReceiverIdentityID,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return messages, nil
}

// CreateDirectMessage inserts a new encrypted direct message into the database.
func CreateDirectMessage(
	ctx context.Context,
	senderID string,
	receiverID string,
	cipherText []byte,
	nonce []byte,
	senderIdentityID string,
	receiverIdentityID string,
) error {

	sql, args, err := psql.
		Insert("direct_message").
		Columns(
			"id",
			"sender_id",
			"receiver_id",
			"cipher_text",
			"nonce",
			"sender_identity_id",
			"receiver_identity_id",
			"created_at",
		).
		Values(
			uuid.New(),
			senderID,
			receiverID,
			cipherText,
			nonce,
			senderIdentityID,
			receiverIdentityID,
			time.Now(),
		).
		ToSql()
	if err != nil {
		return err
	}

	_, err = Pool.Exec(ctx, sql, args...)
	return err
}
