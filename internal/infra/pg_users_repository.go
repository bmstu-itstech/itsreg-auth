package infra

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"

	"github.com/itsreg-auth/internal/domain/auth"
)

type pgUserRepository struct {
	db *sqlx.DB
}

func NewPgUserRepository(db *sqlx.DB) auth.UsersRepository {
	return &pgUserRepository{
		db: db,
	}
}

func (r *pgUserRepository) Save(ctx context.Context, u *auth.User) error {
	row := mapUserToRow(u)
	res, err := pgutils.Exec(
		ctx, r.db,
		`INSERT INTO 
			users (uuid, email, passhash, created_at, updated_at)
         VALUES 
			($1, $2, $3, $4, $5)`,
		row.UUID, row.Email, row.Passhash, row.CreatedAt, row.UpdatedAt)
	if pgutils.IsUniqueViolationError(err) {
		return auth.ErrUserAlreadyExists
	} else if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return errors.New("no affected rows")
	}

	return nil
}

func (r *pgUserRepository) User(ctx context.Context, uuid string) (*auth.User, error) {
	var row userRow
	err := pgutils.Get(
		ctx, r.db, &row,
		`SELECT
			uuid, email, passhash, created_at, updated_at
	     FROM 
			users
		 WHERE 
			uuid = $1`,
		uuid,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, auth.UserNotFound{UserUUID: uuid}
	} else if err != nil {
		return nil, err
	}

	return mapUserFromRow(row)
}

func (r *pgUserRepository) UserByEmail(ctx context.Context, email string) (*auth.User, error) {
	var row userRow
	err := pgutils.Get(
		ctx, r.db, &row,
		`SELECT 
			uuid, email, passhash, created_at, updated_at
         FROM 
			users
         WHERE 
			email = $1`,
		email,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, auth.UserEmailNotFound{Email: email}
	} else if err != nil {
		return nil, err
	}

	return mapUserFromRow(row)
}

func (r *pgUserRepository) Update(
	ctx context.Context,
	uuid string,
	updateFn func(context.Context, *auth.User) error,
) error {
	u, err := r.User(ctx, uuid)
	if err != nil {
		return err
	}

	err = updateFn(ctx, u)
	if err != nil {
		return err
	}

	row := mapUserToRow(u)
	res, err := pgutils.Exec(
		ctx, r.db,
		`UPDATE
			users
         SET 
			email = $2,
            passhash = $3,
			updated_at = $4
         WHERE 
			uuid = $1`,
		row.UUID, row.Email, row.Passhash, row.UpdatedAt,
	)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return auth.UserNotFound{UserUUID: u.UUID}
	}

	return nil
}

func (r *pgUserRepository) Delete(ctx context.Context, uuid string) error {
	res, err := pgutils.Exec(
		ctx, r.db,
		`DELETE FROM
			users
		WHERE
			uuid = $1`,
		uuid,
	)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return nil
	}

	if aff == 0 {
		return auth.UserNotFound{UserUUID: uuid}
	}

	return nil
}

type userRow struct {
	UUID      string    `db:"uuid"`
	Email     string    `db:"email"`
	Passhash  []byte    `db:"passhash"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func mapUserFromRow(row userRow) (*auth.User, error) {
	return auth.NewUserFromDB(
		row.UUID,
		row.Email,
		row.Passhash,
		row.CreatedAt.Local(),
		row.UpdatedAt.Local(),
	)
}

func mapUserToRow(u *auth.User) userRow {
	return userRow{
		UUID:      u.UUID,
		Email:     u.Email,
		Passhash:  u.Passhash,
		CreatedAt: u.CreatedAt.UTC(),
		UpdatedAt: u.UpdatedAt.UTC(),
	}
}
