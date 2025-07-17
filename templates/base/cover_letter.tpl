{{.Name}}
{{.Email}} | {{.Phone}} | {{.Website}} | {{.Github}}
{{.Address}}

{{ "2025-07-17" }}

Hiring Manager
{{.CompanyToApplyTo}}

Dear Hiring Manager,

I am writing to express my interest in the {{.RoleToApplyTo}} position at {{.CompanyToApplyTo}}. 
With a strong background in software development and a passion for building innovative solutions, I am confident I can contribute effectively to your team.

In my most recent role at {{ (index .Experience 0).Company }}, my responsibilities included:
{{- range (index .Experience 0).Responsibilities }}
- {{ . }}
{{- end }}

My technical skills include {{ range $i, $skill := .Skills }}{{ if $i }}, {{ end }}{{ $skill }}{{ end }}.

One of my key projects has been "{{ (index .Projects 0).Name }}", which is {{ (index .Projects 0).Description }}. You can find more details at {{ (index .Projects 0).Url }}.

I am very enthusiastic about the opportunity to bring my skills to {{.CompanyToApplyTo}}. Thank you for your time and consideration.

Sincerely,
{{.Name}}