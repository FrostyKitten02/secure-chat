package repo

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"secure-chat/repo/model"
)

// Search users by username excluding those already in my chats
func SearchUsersByUsernameExcludingChats(ctx context.Context, myID, query string) ([]model.User, error) {
	subQuery1, _, _ := psql.Select("user_1_id").From("chat").Where("user_2_id = ?", myID).ToSql()
	subQuery2, _, _ := psql.Select("user_2_id").From("chat").Where("user_1_id = ?", myID).ToSql()

	sql, args, err := psql.Select("id", "username", "email").
		From("user_tbl").
		Where(sq.And{
			sq.Like{"username": "%" + query + "%"},
			sq.NotEq{"id": myID},
			sq.Expr("id NOT IN ("+subQuery1+")", myID),
			sq.Expr("id NOT IN ("+subQuery2+")", myID),
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
