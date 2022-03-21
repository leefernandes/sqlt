{{ define "user/create" -}}

insert into "iam"."user" (
  email,
  family_name,
  given_name,
  city,
  age,
  create_author,
  create_time
) 
values (
  :email,
  :family_name,
  :given_name,
  :city,
  :age,
  :create_author,
  now()
)
returning *;

{{- end }}
