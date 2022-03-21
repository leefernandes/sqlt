package sqlt

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"text/template"

	"github.com/jmoiron/sqlx"
)

var (
	ErrBindNamed       = errors.New("error binding named vars")
	ErrCreate          = errors.New("error creating record")
	ErrExec            = errors.New("error executing query")
	ErrExecuteTemplate = errors.New("error executing template")
	ErrGet             = errors.New("error getting record")
	ErrSelect          = errors.New("error selecting records")
	ErrUpdate          = errors.New("error updating record")
)

type SQLT interface {
	Create(ctx context.Context, name string, data any) error
	Exec(ctx context.Context, name string, data any) (sql.Result, error)
	ExecuteTemplate(name string, data interface{}) (string, []any, error)
	Get(ctx context.Context, name string, data, dest any) error
	Select(ctx context.Context, name string, data, dest any) error
	Update(ctx context.Context, name string, data any) error
}

func Must(db *sqlx.DB, templates embed.FS, patterns []string, opts ...Option) SQLT {
	sqllib, err := New(db, templates, patterns, opts...)
	if err != nil {
		panic(err)
	}
	return sqllib
}

func New(db *sqlx.DB, templates embed.FS, patterns []string, opts ...Option) (SQLT, error) {
	tmpl, err := template.ParseFS(templates, patterns...)
	if err != nil {
		return nil, err
	}

	s := &sqlt{}

	for _, opt := range opts {
		opt(s)
	}

	s.db = db
	s.tmpl = tmpl

	return s, nil
}

// Option for use when calling sqlt.New.
type Option func(*sqlt)

// Debug enables debug logging.
func Debug() Option {
	return func(s *sqlt) {
		s.debug = true
	}
}

type sqlt struct {
	db    *sqlx.DB
	debug bool
	tmpl  *template.Template
}

// Create a record by parsing the named template
// with data and executing the query.
//
// Any named bindvar (example: :my_field_name) in the template
// will be replaced with database specific placeholders (? or $1..$N)
// along with the associated encoded parameters from data input.
func (s *sqlt) Create(ctx context.Context, name string, data any) error {
	sql, args, err := s.ExecuteTemplate(name, data)
	if err != nil {
		return fmt.Errorf("%v: %w", err, ErrCreate)
	}

	if err := s.db.QueryRowxContext(ctx, sql, args...).StructScan(data); err != nil {
		if s.debug {
			fmt.Println(sql, args)
		}
		return fmt.Errorf("%v: %w", err, ErrCreate)
	}

	return nil
}

// Get a record by parsing the named template
// with data and executing the query.
//
// Any named bindvar (example: :my_field_name) in the template
// will be replaced with database specific placeholders (? or $1..$N)
// along with the associated encoded parameters from data input.
func (s *sqlt) Get(ctx context.Context, name string, data, dest any) error {
	sql, args, err := s.ExecuteTemplate(name, data)
	if err != nil {
		return fmt.Errorf("%v: %w", err, ErrGet)
	}

	if err := s.db.GetContext(ctx, dest, sql, args...); err != nil {
		if s.debug {
			fmt.Println(sql, args)
		}
		return fmt.Errorf("%v: %w", err, ErrGet)
	}

	return nil
}

// Exec a query by parsing the named template
// with data and executing the query.
//
// Any named bindvar (example: :my_field_name) in the template
// will be replaced with database specific placeholders (? or $1..$N)
// along with the associated encoded parameters from data input.
func (s *sqlt) Exec(ctx context.Context, name string, data any) (sql.Result, error) {
	sql, args, err := s.ExecuteTemplate(name, data)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, ErrExec)
	}

	res, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		if s.debug {
			fmt.Println(sql, args)
		}
		return nil, fmt.Errorf("%v: %w", err, ErrExec)
	}

	return res, nil
}

// ExecuteTemplate parses the named template with data,
// returing the sql and args to be executed.
func (s *sqlt) ExecuteTemplate(name string, data any) (string, []any, error) {
	var b bytes.Buffer

	if err := s.tmpl.ExecuteTemplate(&b, name, data); err != nil {
		return "", nil, fmt.Errorf("%v: %w", err, ErrExecuteTemplate)
	}

	if data == nil {
		return b.String(), nil, nil
	}

	//sql, args, err := s.db.BindNamed(b.String(), data)
	sql, args, err := sqlx.Named(b.String(), data)
	if err != nil {
		if s.debug {
			fmt.Println(b.String())
		}
		return "", nil, fmt.Errorf("%v: %w", err, ErrBindNamed)
	}

	sql, args, err = sqlx.In(sql, args...)
	if err != nil {
		if s.debug {
			fmt.Println(b.String())
		}
		return "", nil, fmt.Errorf("%v: %w", err, ErrBindNamed)
	}

	sql = s.db.Rebind(sql)

	return sql, args, nil
}

// Select records by parsing the named template
// with data and executing the query.
//
// Any named bindvar (example: :my_field_name) in the template
// will be replaced with database specific placeholders (? or $1..$N)
// along with the associated encoded parameters from data input.
func (s *sqlt) Select(ctx context.Context, name string, data, dest any) error {
	sql, args, err := s.ExecuteTemplate(name, data)
	if err != nil {
		return fmt.Errorf("%v: %w", err, ErrSelect)
	}

	if err := s.db.Select(dest, sql, args...); err != nil {
		if s.debug {
			fmt.Println(sql, args)
		}
		return fmt.Errorf("%v: %w", err, ErrSelect)
	}

	return nil
}

// Update a record by parsing the named template
// with data and executing the query.
//
// Any named bindvar (example: :my_field_name) in the template
// will be replaced with database specific placeholders (? or $1..$N)
// along with the associated encoded parameters from data input.
func (s *sqlt) Update(ctx context.Context, name string, data any) error {
	sql, args, err := s.ExecuteTemplate(name, data)
	if err != nil {
		return fmt.Errorf("%v: %w", err, ErrUpdate)
	}

	err = s.db.QueryRowxContext(ctx, sql, args...).StructScan(data)
	if err != nil {
		if s.debug {
			fmt.Println(sql, args)
		}
		return fmt.Errorf("%w: %v", err, ErrUpdate)
	}

	return nil
}
