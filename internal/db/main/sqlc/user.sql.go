// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user.sql

package sqlcdb

import (
	"context"

	uuid "github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const CountUsers = `-- name: CountUsers :one
SELECT COUNT(*) FROM users WHERE deleted_at IS NULL
`

func (q *Queries) CountUsers(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, CountUsers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const CreateUser = `-- name: CreateUser :exec
INSERT INTO users (id, first_name, last_name, email, phone, password) 
VALUES ($1, $2, $3, $4, $5, $6)
`

type CreateUserParams struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Password  string
}

func (q *Queries) CreateUser(ctx context.Context, arg *CreateUserParams) error {
	_, err := q.db.Exec(ctx, CreateUser,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Phone,
		arg.Password,
	)
	return err
}

const DeleteSessionID = `-- name: DeleteSessionID :exec
UPDATE users SET session_id = NULL, session_expired_at = NULL WHERE id = $1
`

func (q *Queries) DeleteSessionID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, DeleteSessionID, id)
	return err
}

const DeleteUserByID = `-- name: DeleteUserByID :exec
UPDATE users SET deleted_at = NOW() WHERE id = $1
`

func (q *Queries) DeleteUserByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, DeleteUserByID, id)
	return err
}

const GetUserByEmail = `-- name: GetUserByEmail :one
SELECT id, first_name, last_name, email, phone, password, session_id, email_confirmed_at, phone_confirmed_at, created_at, updated_at, session_expired_at, deleted_at FROM users WHERE email = $1 AND deleted_at IS NULL
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := q.db.QueryRow(ctx, GetUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.Password,
		&i.SessionID,
		&i.EmailConfirmedAt,
		&i.PhoneConfirmedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SessionExpiredAt,
		&i.DeletedAt,
	)
	return &i, err
}

const GetUserByID = `-- name: GetUserByID :one
SELECT id, first_name, last_name, email, phone, password, session_id, email_confirmed_at, phone_confirmed_at, created_at, updated_at, session_expired_at, deleted_at FROM users WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	row := q.db.QueryRow(ctx, GetUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.Password,
		&i.SessionID,
		&i.EmailConfirmedAt,
		&i.PhoneConfirmedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SessionExpiredAt,
		&i.DeletedAt,
	)
	return &i, err
}

const GetUserByPhone = `-- name: GetUserByPhone :one
SELECT id, first_name, last_name, email, phone, password, session_id, email_confirmed_at, phone_confirmed_at, created_at, updated_at, session_expired_at, deleted_at FROM users WHERE phone = $1 AND deleted_at IS NULL
`

func (q *Queries) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
	row := q.db.QueryRow(ctx, GetUserByPhone, phone)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.Password,
		&i.SessionID,
		&i.EmailConfirmedAt,
		&i.PhoneConfirmedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SessionExpiredAt,
		&i.DeletedAt,
	)
	return &i, err
}

const ListUsers = `-- name: ListUsers :many
SELECT id, first_name, last_name, email, phone, password, session_id, email_confirmed_at, phone_confirmed_at, created_at, updated_at, session_expired_at, deleted_at
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type ListUsersParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) ListUsers(ctx context.Context, arg *ListUsersParams) ([]*User, error) {
	rows, err := q.db.Query(ctx, ListUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.Password,
			&i.SessionID,
			&i.EmailConfirmedAt,
			&i.PhoneConfirmedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.SessionExpiredAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const UpdateSessionID = `-- name: UpdateSessionID :exec
UPDATE users SET session_id = $1, session_expired_at = $2 WHERE id = $3
`

type UpdateSessionIDParams struct {
	SessionID        *string
	SessionExpiredAt pgtype.Timestamp
	ID               uuid.UUID
}

func (q *Queries) UpdateSessionID(ctx context.Context, arg *UpdateSessionIDParams) error {
	_, err := q.db.Exec(ctx, UpdateSessionID, arg.SessionID, arg.SessionExpiredAt, arg.ID)
	return err
}

const UpdateUser = `-- name: UpdateUser :exec
UPDATE users
SET first_name = $2, last_name = $3, email = $4, phone = $5, password = $6
WHERE id = $1
`

type UpdateUserParams struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Password  string
}

func (q *Queries) UpdateUser(ctx context.Context, arg *UpdateUserParams) error {
	_, err := q.db.Exec(ctx, UpdateUser,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Phone,
		arg.Password,
	)
	return err
}