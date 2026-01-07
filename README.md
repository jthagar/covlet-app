# Covlet — Cover Letter Templates Editor

Covlet is a desktop app for creating and managing cover‑letter templates using Go text/template. It includes a live editor, file browser, and a variable sidebar that detects {{ .Vars }} and lets you override them before rendering. You can preview and export the output as a PDF.

Covlet is built with the Fyne GUI toolkit and Go.


## Requirements
- Go 1.24 or newer
- Fyne v2 (handled via Go modules)
- A desktop environment supported by Fyne (Linux, macOS, Windows)


## Features
- Go text/template editing
- Variable sidebar: detects top‑level variables like {{ .Name }} and allows quick overrides
- Dual file trees for navigating your templates
- Render preview using config.yml + overrides
- Export to PDF (default save: `~/Downloads/covlet` on Linux)
- App home with organized `templates/` and `values/`


## Installation
Clone the repository and build:

```
git clone https://github.com/your-org-or-user/cover-letter-templates.git
cd cover-letter-templates
go build ./...
```

You can run the app directly:

```
go run .
```

Or build a binary and run it:

```
go build -o covlet
./covlet
```


## Running the Application
On first run Covlet prepares an application home directory. By default on Linux this is:

- `~/.local/share/covlet`
  - `templates/` — where your template files live
  - `values/` — where you can store YAML values files (future use)

You can change the app home directory by setting the environment variable before starting Covlet:

```
export COVLET_HOME=/path/to/your/covlet-home
go run .
```


## Templates
Covlet looks for templates in `<COVLET_HOME>/templates`. You can organize into subfolders such as `base/` and `partials/`. Supported file types: `.tpl`, `.tmpl`, `.txt`, `.md`, `.gohtml`, `.html`.

Minimal example (`templates/base/cover_letter.tpl`):

```
{{ .Name }}
{{ .Email }} | {{ .Phone }}

Dear Hiring Manager,

I am applying for {{ .RoleToApplyTo }} at {{ .CompanyToApplyTo }}.

Regards,
{{ .Name }}
```


## Configuration (config.yml)
Rendering uses data from a `config.yml` in your working directory plus any sidebar overrides.

Minimal example:

```
name: Jane Doe
email: jane@example.com
phone: +1 555 000 0000
company_to_apply_to: ACME Corp
role_to_apply_to: Senior Go Engineer
```

Note: The variable sidebar focuses on top‑level fields (e.g., `.Name`, `.Email`, `.CompanyToApplyTo`, `.RoleToApplyTo`). Nested data can be used in templates but is not auto‑surfaced yet.


## Rendering and PDF Export
1. Open a template in the editor.
2. Ensure `config.yml` exists with your data; add overrides in the sidebar if needed.
3. Click “Render” to preview.
4. In the preview window choose File → “Export as PDF…”. The file is saved as `<title>.pdf` to `~/Downloads/covlet` on Linux by default.


## Testing and CI
Unit tests cover core non-GUI logic. To run locally:

```
go test ./...
```

This repository includes a GitHub Actions workflow that automatically runs tests on pushes and pull requests.


## Troubleshooting
- If templates don’t appear in the file trees, ensure your `COVLET_HOME` is set as expected and that the `templates/` directory exists. The application will create it on first run if missing.
- If PDF export fails, confirm that `~/Downloads/covlet` exists or that Covlet has permission to create it. The app attempts to create it automatically.
- On first run, the theme uses a slightly reduced base font size for better density; you can switch to light/dark from the View menu.


## License
This project is licensed under the terms of the MIT License. See the LICENSE file for details.