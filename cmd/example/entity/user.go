package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `db:"id"`
	Email      string    `db:"email"`
	FamilyName string    `db:"family_name"`
	GivenName  string    `db:"given_name"`
	City       string    `db:"city"`
	Age        int       `db:"age"`

	CreateAuthor uuid.UUID  `db:"create_author"`
	CreateTime   time.Time  `db:"create_time"`
	UpdateAuthor *uuid.UUID `db:"update_author"`
	UpdateTime   *time.Time `db:"update_time"`
	DeleteAuthor *uuid.UUID `db:"delete_author"`
	DeleteTime   *time.Time `db:"delete_time"`
}

type UserListInput struct {
	Age    int
	Cities []string
	Limit  int
	Select string
	Where  string
}
