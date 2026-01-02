package repo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"secure-chat/repo/model"
)

func FindActiveIdentityForUser(ctx context.Context, userID string) (*model.Identity, error) {
	sql, args, err := psql.
		Select("id", "pub_key", "enc_priv_key", "active", "user_id").
		From("identity_tbl").
		Where("user_id = ?", userID).
		Where("active = ?", true).
		ToSql()
	if err != nil {
		return nil, err
	}

	var identity model.Identity
	err = Pool.QueryRow(ctx, sql, args...).Scan(
		&identity.ID,
		&identity.PubKey,
		&identity.EncPrivKey,
		&identity.Active,
		&identity.UserID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // no active identity
		}
		return nil, err
	}

	return &identity, nil
}
