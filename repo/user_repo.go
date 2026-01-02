package repo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"secure-chat/repo/model"
)

func FindUserByUsername(ctx context.Context, username string) (*model.User, error) {
	sql, args, err := psql.
		Select("id", "username", "email", "password", "created_at").
		From("user_tbl").
		Where("username = ?", username).
		ToSql()
	if err != nil {
		return nil, err
	}

	var u model.User
	err = Pool.QueryRow(ctx, sql, args...).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // not found
		}
		return nil, err
	}

	return &u, nil
}

func FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	sql, args, err := psql.
		Select("id", "username", "email", "password", "created_at").
		From("user_tbl").
		Where("email = ?", email).
		ToSql()
	if err != nil {
		return nil, err
	}

	var u model.User
	err = Pool.QueryRow(ctx, sql, args...).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func SaveUser(ctx context.Context, user model.User, identity model.Identity) error {
	tx, poolErr := Pool.Begin(ctx)
	if poolErr != nil {
		return poolErr
	}
	defer tx.Rollback(ctx)

	// --- insert user ---
	userSQL, userArgs, usrPsqlErr := psql.
		Insert("user_tbl").
		Columns("id", "username", "email", "password").
		Values(
			user.ID,
			user.Username,
			user.Email,
			user.Password,
		).
		ToSql()
	if usrPsqlErr != nil {
		return usrPsqlErr
	}

	_, saveUsrErr := tx.Exec(ctx, userSQL, userArgs...)
	if saveUsrErr != nil {
		return saveUsrErr
	}

	// --- insert identity ---
	identitySQL, identityArgs, identityPsqlErr := psql.
		Insert("identity_tbl").
		Columns("id", "pub_key", "enc_priv_key", "active", "user_id").
		Values(
			identity.ID,
			identity.PubKey,
			identity.EncPrivKey,
			identity.Active,
			user.ID,
		).
		ToSql()
	if identityPsqlErr != nil {
		return identityPsqlErr
	}

	_, saveIdentityErr := tx.Exec(ctx, identitySQL, identityArgs...)
	if saveIdentityErr != nil {
		return saveIdentityErr
	}

	return tx.Commit(ctx)
}

func FindUserByID(ctx context.Context, id string) (*model.User, error) {
	sql, args, err := psql.
		Select("id", "username", "email", "password", "created_at").
		From("user_tbl").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, err
	}

	var u model.User
	err = Pool.QueryRow(ctx, sql, args...).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // user not found
		}
		return nil, err
	}

	return &u, nil
}
