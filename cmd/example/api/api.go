package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/leefernandes/sqlt"
	"github.com/leefernandes/sqlt/cmd/example/domain"
)

var (
	ErrCreateUserSchema = errors.New("error creating user schema")
	ErrCreateUser       = errors.New("error creating user")
	ErrGetUser          = errors.New("error getting user")
	ErrListUsers        = errors.New("error listing users")
	ErrUpdateUser       = errors.New("error updating user")
)

func New(db *sqlx.DB, lib sqlt.SQLT) *API {
	return &API{
		db:  db,
		lib: lib,
	}
}

type API struct {
	db  *sqlx.DB
	lib sqlt.SQLT
}

func Context(duration time.Duration, author *uuid.UUID) (api_context, context.CancelFunc) {
	if duration == 0 {
		duration = 1 * time.Second
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), duration)

	return api_context{
		Context: ctx,
		author:  author,
	}, cancelFunc
}

type api_context struct {
	context.Context
	author *uuid.UUID
}

func (a api_context) Author() *uuid.UUID {
	return a.author
}

func (api *API) CreateUserSchema(ctx context.Context) error {
	_, err := api.lib.Exec(ctx, "user/schema")
	if err != nil {
		return fmt.Errorf("%v: %w", err, ErrCreateUserSchema)
	}

	return nil
}

func (api *API) CreateUser(ctx api_context, user *domain.User) error {
	user.CreateAuthor = *ctx.Author()

	if err := api.lib.QueryRow(ctx, "user/create", user, sqlt.Input(user)); err != nil {
		return fmt.Errorf("%v: %w", err, ErrCreateUser)
	}

	return nil
}

func (api *API) GetUser(ctx api_context, id uuid.UUID) (*domain.User, error) {
	user := &domain.User{
		ID: id,
	}

	if err := api.lib.QueryRow(ctx, "user/get", user, sqlt.Input(user)); err != nil {
		return nil, fmt.Errorf("%v: %w", err, ErrGetUser)
	}

	return user, nil
}

func (api *API) ListUsers(ctx api_context, query domain.UserListQuery) ([]domain.User, error) {
	var users []domain.User

	err := api.lib.Query(ctx, "user/list", &users, sqlt.Input(query))

	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, ErrListUsers)
	}

	return users, nil
}

func (api *API) UserJob(ctx api_context) error {
	err := api.lib.Iterate(ctx, "user/list", func(scan func(dest any) error) error {
		user := &domain.User{}
		scan(user)
		fmt.Println("UserJob completed for User:", user.ID)
		return nil
	})

	return err
}

func (api *API) UpdateUser(ctx api_context, user *domain.User) error {
	user.UpdateAuthor = ctx.Author()

	err := api.lib.QueryRow(ctx, "user/update", user, sqlt.Input(user))
	if err != nil {
		return fmt.Errorf("%v: %w", err, ErrUpdateUser)
	}

	return nil
}