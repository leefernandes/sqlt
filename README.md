# SQLT - sql/template

A lean experimental library for rendering SQL from embedded text/template files that supports ? or $1...$N parameterized placeholders, and scans results into structs.

See [cmd/example/main.go](cmd/example/main.go) for a functional CRUD+ api example.

## EXAMPLE

Given a sql/template ...

```sql
{{ define "user/list" -}}

-- SELECT
{{ with .Select -}} select {{ . }} {{ else -}} select * {{ end }}
-- FROM
from "iam"."user"
-- WHERE
{{ with .Where -}} where {{ . }} {{ end }}
-- LIMIT
{{ with .Limit -}}
{{ if and (gt . 0) (lt . 50) -}} limit {{ . }} {{ else -}} limit 50 {{ end }}
{{- end }}

{{- end }}
```

... and some input data ...

```go
userListInput := UserListInput{
  Where:  "city in (:cities) and age > :age",
  Limit:  10,
  Age:    98,
  Cities: []string{"Tampa", "São Paulo", "Rio de Janeiro"},
}
```

... we should generate the following SQL ...

```sql
-- SELECT
select *
-- FROM
from "iam"."user"
-- WHERE
where city in ($1, $2, $3) and age > $4
-- LIMIT
limit 10
```

... where the encoded placeholder parameters equal the following ...

```sh
$1 = 'Tampa'
$2 = 'São Paulo'
$3 = 'Rio de Janeiro'
$4 = 98
```

... and are easily scanned into in-memory structures, or [iterated](https://github.com/leefernandes/sqlt/blob/main/cmd/example/api/api.go#L111) for large per row jobs.

```go
//go:embed sql/**/*
var templates embed.FS

// name of our template to render.
templateName := "user/list"

// input data to execute against our template.
userListInput := UserListQuery{
  Where:  "city in (:cities) and age > :age",
  Limit:  10,
  Age:    98,
  Cities: []string{"Tampa", "São Paulo", "Rio de Janeiro"},
}

// users slice to scan rows into.
users := []User{}

err := api.lib.Query(ctx, templateName, &users, sqlt.Input(userListQuery))
```

SQLT is comprised of 5 lightweight methods named similarly to Go's standard templating & sql library names, `Exec` `ExecuteTemplate` `Iterate` `Query` `QueryRow`

We use embed.FS to load templates, standard go text/template for parsing templates, and [sqlx](https://github.com/jmoiron/sqlx) to interface with databases.

```go
type SQLT interface {
	Exec(ctx context.Context, name string, opts ...QueryOption) (sql.Result, error)
	ExecuteTemplate(name string, data interface{}) (string, []any, error)
	Iterate(ctx context.Context, name string, iter Iterator, opts ...QueryOption) error
	Query(ctx context.Context, name string, dest any, opts ...QueryOption) error
	QueryRow(ctx context.Context, name string, dest any, opts ...QueryOption) error
}
```

Inspired by: [github.com/Davmuz/gqt](https://github.com/Davmuz/gqt)
