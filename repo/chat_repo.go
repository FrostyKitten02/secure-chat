package repo

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"secure-chat/repo/model"
)

func FindChatsByUserID(ctx context.Context, userID string) ([]model.Chat, error) {
	sql, args, err := psql.
		Select("id", "user_1_id", "user_2_id").
		From("chat").
		Where(
			sq.Or{
				sq.Eq{"user_1_id": userID},
				sq.Eq{"user_2_id": userID},
			},
		).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := make([]model.Chat, 0)
	for rows.Next() {
		var c model.Chat
		err := rows.Scan(
			&c.ID,
			&c.User1ID,
			&c.User2ID,
		)
		if err != nil {
			return nil, err
		}
		chats = append(chats, c)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return chats, nil
}

func CreateChatIfNotExists(ctx context.Context, user1, user2 string) error {
	if user1 > user2 {
		user1, user2 = user2, user1
	}

	sql, args, err := psql.
		Insert("chat").
		Columns("id", "user_1_id", "user_2_id").
		Values(uuid.New(), user1, user2).
		Suffix("ON CONFLICT (user_1_id, user_2_id) DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}

	_, err = Pool.Exec(ctx, sql, args...)
	return err
}
