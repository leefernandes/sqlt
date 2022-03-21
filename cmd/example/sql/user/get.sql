{{ define "user/get" -}}

select *
from "iam"."user"
where ID = :id
limit 1

{{- end }}
