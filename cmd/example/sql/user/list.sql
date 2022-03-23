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
