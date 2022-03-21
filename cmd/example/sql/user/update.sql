{{ define "user/update" -}}

-- UPDATE
update "iam"."user"
-- SET
set 
{{- if .Email }}
  email = :email,
{{- end }}

{{- if .FamilyName }}
  family_name = :family_name,
{{- end }}

{{- if .GivenName }}
  given_name = :given_name,
{{- end }}

{{- if .City }}
  city = :city,
{{- end }}

{{- if .Age }}
  age = :age,
{{- end }}

  update_author = :update_author,
  update_time = now()
-- WHERE
where id = :id
returning *;

{{- end }}
