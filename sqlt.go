package sqlt

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"fmt"
	"text/template"

	"github.com/jmoiron/sqlx"
)

type SQLT interface {
	Exec(ctx context.Context, name string, opts ...QueryOption) (sql.Result, error)
	ExecuteTemplate(name string, data interface{}) (string, []any, error)
	Iterate(ctx context.Context, name string, iter Iterator, opts ...QueryOption) error
	Query(ctx context.Context, name string, dest any, opts ...QueryOption) error
	QueryRow(ctx context.Context, name string, dest any, opts ...QueryOption) error
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

// Option for use when calling Must & New.
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

// QueryOption
type QueryOption func(*queryoptions)

type queryoptions struct {
	Input any
}

// Input is a QueryOption that sets the input data to be used when
// executing the template.
// Named bindvars (example: :my_field_name) in the sql template
// will be replaced with database specific placeholders (? or $1..$N)
// mapped to the asscociated input field-value (example: data.MyFieldName)
// to be sent as encoded parameters.
func Input(input any) QueryOption {
	return func(q *queryoptions) {
		q.Input = input
	}
}

// Exec executes the named template, and returns sql.Result.
// Use Exec when the sql template does not return data.
func (s *sqlt) Exec(ctx context.Context, name string, opts ...QueryOption) (sql.Result, error) {
	q := &queryoptions{}
	for _, opt := range opts {
		opt(q)
	}

	sql, args, err := s.ExecuteTemplate(name, q.Input)
	if err != nil {
		return nil, err
	}

	res, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		if s.debug {
			fmt.Println(sql, args)
		}
		return nil, err
	}

	return res, nil
}

// ExecuteTemplate parses the named template with input data,
// and returns the sql and args to be executed.
// Named bindvars (example: :my_field_name) in the sql template
// will be replaced with database specific placeholders (? or $1..$N)
// mapped to the asscociated input field-value (example: input.MyFieldName)
// to be sent as encoded parameters.
func (s *sqlt) ExecuteTemplate(name string, input any) (string, []any, error) {
	var b bytes.Buffer

	if err := s.tmpl.ExecuteTemplate(&b, name, input); err != nil {
		return "", nil, err
	}

	if input == nil {
		return b.String(), nil, nil
	}

	//sql, args, err := s.db.BindNamed(b.String(), data)
	sql, args, err := sqlx.Named(b.String(), input)
	if err != nil {
		if s.debug {
			fmt.Println(b.String())
		}
		return "", nil, err
	}

	sql, args, err = sqlx.In(sql, args...)
	if err != nil {
		if s.debug {
			fmt.Println(b.String())
		}
		return "", nil, err
	}

	sql = s.db.Rebind(sql)

	if s.debug {
		fmt.Println("")
		fmt.Println(sql)
		fmt.Println("")
		for i, v := range args {
			switch v.(type) {
			case nil:
				v = "NULL"
			case int, int8, int16, int32, int64, float32, float64, uint, uint8, uint16, uint32, uint64:
			default:
				v = fmt.Sprintf("'%v'", v)
			}
			fmt.Printf("$%d = %v\n", i+1, v)
		}
		fmt.Println("")
	}

	return sql, args, nil
}

type Iterator func(scan func(dest any) error) error

// Iterate executes the named template, and iterates over the results
// scanning each row into dest struct.
// Use Iterate when dealing with potentially unbounded rows.
func (s *sqlt) Iterate(ctx context.Context, name string, iter Iterator, opts ...QueryOption) error {
	q := &queryoptions{}
	for _, opt := range opts {
		opt(q)
	}

	sql, args, err := s.ExecuteTemplate(name, q.Input)
	if err != nil {
		return err
	}

	rows, err := s.db.Queryx(sql, args...)
	if err != nil {
		return err
	}

	for rows.Next() {
		scan := func(dest any) error {
			err := rows.StructScan(dest)
			return err
		}

		if err := iter(scan); err != nil {
			return err
		}
	}

	return nil
}

// Query executes the named template, and scans rows into dest slice.
// Use Query when many rows should be scanned into a slice in memory.
func (s *sqlt) Query(ctx context.Context, name string, dest any, opts ...QueryOption) error {
	q := &queryoptions{}
	for _, opt := range opts {
		opt(q)
	}

	sql, args, err := s.ExecuteTemplate(name, q.Input)
	if err != nil {
		return err
	}

	if err := s.db.Select(dest, sql, args...); err != nil {
		if s.debug {
			fmt.Println(sql, args)
		}
		return err
	}

	return nil
}

// QueryRow executes the named template, and scans row into dest struct.
// Use Query when a single row should be scanned into a struct in memory.
func (s *sqlt) QueryRow(ctx context.Context, name string, dest any, opts ...QueryOption) error {
	q := &queryoptions{}
	for _, opt := range opts {
		opt(q)
	}

	sql, args, err := s.ExecuteTemplate(name, q.Input)
	if err != nil {
		return err
	}

	if err := s.db.QueryRowxContext(ctx, sql, args...).StructScan(dest); err != nil {
		if s.debug {
			fmt.Println(sql, args)
		}
		return err
	}

	return nil
}
