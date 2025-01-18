// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.819
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import t "github.com/a-h/templ"

type PageData struct {
	Title   string
	Content t.Component
}

func Layout(data PageData) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><meta name=\"description\" content=\"Goforms - A self-hosted form backend service built with Go\"><meta name=\"keywords\" content=\"forms, golang, backend, self-hosted\"><meta name=\"color-scheme\" content=\"light dark\"><title>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(data.Title)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/components/layout.templ`, Line: 19, Col: 21}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, " - Goforms</title><!-- Single CSS Entry Point --><link rel=\"stylesheet\" href=\"/static/css/main.css\"></head><body><nav><div class=\"container\"><a href=\"/\" class=\"logo\">Goforms</a><div class=\"nav-links\"><a href=\"/contact\">Contact Demo</a> <a href=\"https://github.com/jonesrussell/goforms\">GitHub</a> <button onclick=\"toggleTheme()\" class=\"btn btn-secondary btn-sm\">Toggle Theme</button></div></div></nav><main class=\"container\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = data.Content.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "</main><footer><div class=\"container\">Created by <a href=\"https://jonesrussell.github.io/me/\" target=\"_blank\" rel=\"noopener noreferrer\">Russell Jones</a></div></footer><!-- Theme Toggle Script --><script>\n\t\t\t// Use modern JavaScript features\n\t\t\tconst themeToggle = () => {\n\t\t\t\t// Get current theme or system preference\n\t\t\t\tconst getTheme = () => \n\t\t\t\t\tlocalStorage.getItem('theme') || \n\t\t\t\t\t(window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');\n\n\t\t\t\t// Toggle between light and dark\n\t\t\t\tconst theme = getTheme() === 'dark' ? 'light' : 'dark';\n\t\t\t\t\n\t\t\t\t// Update DOM and storage\n\t\t\t\tdocument.documentElement.dataset.theme = theme;\n\t\t\t\tlocalStorage.setItem('theme', theme);\n\t\t\t};\n\n\t\t\t// Initialize theme\n\t\t\tdocument.documentElement.dataset.theme = \n\t\t\t\tlocalStorage.getItem('theme') || \n\t\t\t\t(window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');\n\n\t\t\t// Expose toggle function\n\t\t\twindow.toggleTheme = themeToggle;\n\t\t</script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
