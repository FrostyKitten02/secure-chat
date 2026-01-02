package repo

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"secure-chat/repo/model"
)

// Search users by username excluding those already in my chats
func SearchUsersByUsernameExcludingChats(ctx context.Context, myID uuid.UUID, query string) ([]model.User, error) {
	sql, args, err := psql.
		Select("id", "username", "email").
		From("user_tbl").
		Where(sq.And{
			sq.Like{"username": "%" + query + "%"},
			sq.NotEq{"id": myID},
			sq.Expr("id NOT IN (SELECT user_1_id FROM chat WHERE user_2_id = ?::uuid)", myID),
			sq.Expr("id NOT IN (SELECT user_2_id FROM chat WHERE user_1_id = ?::uuid)", myID),
		}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
