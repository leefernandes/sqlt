package main

import (
	"embed"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/leefernandes/sqlt"
	"github.com/leefernandes/sqlt/cmd/example/api"
	"github.com/leefernandes/sqlt/cmd/example/entity"
)

//go:embed sql/**/*
var templates embed.FS

func main() {
	db := sqlx.MustConnect("postgres", "dbname=postgres sslmode=disable")

	lib := sqlt.Must(db, templates, []string{"sql/**/*"}, sqlt.Debug())

	me := uuid.New()

	ctx, cancel := api.Context(5*time.Second, &me)
	defer cancel()

	dbapi := api.New(db, lib)

	if err := dbapi.CreateUserSchema(ctx); err != nil {
		panic(err)
	}

	//
	// CreateUser
	//
	user := &entity.User{
		City:  "São Paulo",
		Email: fmt.Sprintf("notan@email+%d.lol", time.Now().Unix()),
	}

	err := dbapi.CreateUser(ctx, user)
	if err != nil {
		panic(err)
	}

	spew.Dump("CreateUser:", user)

	//
	// CreateUser
	//
	user2 := &entity.User{
		City:  "Tampa",
		Email: fmt.Sprintf("stillnotan@email+%d.lol", time.Now().Unix()),
	}

	err = dbapi.CreateUser(ctx, user2)
	if err != nil {
		panic(err)
	}

	spew.Dump("CreateUser:", user2)

	//
	// GetUser
	//
	user, err = dbapi.GetUser(ctx, user.ID)
	if err != nil {
		panic(err)
	}

	spew.Dump("GetUser:", user)

	//
	// UpdateUser
	//
	user.Age = 99

	err = dbapi.UpdateUser(ctx, user)
	if err != nil {
		panic(err)
	}

	spew.Dump("UpdateUser:", user)

	//
	// ListUsers
	//
	userListInput := entity.UserListInput{
		Where:  "city in (:cities) and age > :age",
		Limit:  10,
		Age:    98,
		Cities: []string{"Tampa", "São Paulo", "Rio de Janeiro"},
	}

	users, err := dbapi.ListUsers(ctx, userListInput)
	if err != nil {
		panic(err)
	}

	spew.Dump("ListUsers:", users)

	//
	// UserJob
	//
	if err = dbapi.UserJob(ctx); err != nil {
		panic(err)
	}

}
