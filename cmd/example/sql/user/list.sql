{{ define "user/list" -}}

-- SELECT
{{ if .Select -}} select {{ .Select }} {{ else -}} select * {{ end }}
-- FROM
from "iam"."user"
-- WHERE
{{ if .Where -}} where {{ .Where }} {{ end }}
-- LIMIT
{{ if and (gt .Limit 0) (lt .Limit 50) -}} limit {{ .Limit }} {{ else -}} limit 50 {{ end }}

{{- end }}
