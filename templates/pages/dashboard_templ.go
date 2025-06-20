// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.898
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/regcomp/gdpr/templates/components"

func Dashboard(accessToken, refreshToken, sessionID string) templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"en\"><head>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.HeadContents("Dashboard").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<style>\n        :root {\n          --navbar-height: 40px;\n        }\n\n        .navbar {\n          height: var(--navbar-height);\n          min-height: var(--navbar-height);\n          padding: 0.25rem 1rem;\n        }\n\n        .main-container {\n          margin-top: var(--navbar-height);\n        }\n\n        .sidebar-height {\n          height: calc(100vh - var(--navbar-height));\n        }\n\n        .sidebar-nav {\n          margin-bottom: 0.5rem;\n        }\n      </style></head><body><!-- Top Bar - Full Width --><nav class=\"navbar navbar-expand-lg py-0 border-bottom fixed-top\"><div class=\"container-fluid\"><!-- Logo/Brand --><div class=\"navbar-brand mb-0 h5\">GDPR Compliance Utility</div></div></nav><!-- Main Layout Container - Below Fixed Top Bar --><div class=\"container-fluid p-0 main-container\"><div class=\"row g-0\"><!-- Sidebar --><div class=\"col-md-3 col-lg-2 sidebar border-end\"><div class=\"d-flex flex-column sidebar-height\"><!-- Navigation --><nav class=\"flex-grow-1 p-3\"><ul class=\"nav nav-pills flex-column\"><li class=\"nav-item sidebar-nav\"><a href=\"#\" class=\"nav-link active\">Dashboard</a></li><li class=\"nav-item sidebar-nav\"><a href=\"#\" class=\"nav-link\">Analytics</a></li><li class=\"nav-item sidebar-nav\"><a href=\"#\" class=\"nav-link\">Users</a></li><li class=\"nav-item sidebar-nav\"><a href=\"#\" class=\"nav-link\">Projects</a></li><li class=\"nav-item sidebar-nav\"><a href=\"#\" class=\"nav-link\">Reports</a></li><li class=\"nav-item\"><a href=\"#\" class=\"nav-link\">Settings</a></li></ul></nav><!-- Sidebar Footer --><div class=\"p-3 border-top border-secondary\"><div class=\"d-flex align-items-center\"><div class=\"bg-primary rounded-circle d-flex align-items-center justify-content-center me-2\" style=\"width: 32px; height: 32px;\"><small class=\"fw-bold\">ID</small></div><div class=\"flex-grow-1\"><div class=\"small\">Identifier</div><div class=\"text-muted small\">Admin</div></div></div></div></div></div><!-- Main Content Area --><div class=\"col-md-9 col-lg-10\"><main class=\"p-4\"><div class=\"mb-4\"><h2 class=\"mb-1\">Dashboard</h2><div>access token: ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(accessToken)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/dashboard.templ`, Line: 104, Col: 40}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "</div><div>refresh token: ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(refreshToken)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/dashboard.templ`, Line: 105, Col: 42}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "</div><div>session id: ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(sessionID)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/pages/dashboard.templ`, Line: 106, Col: 36}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "</div></div></main></div></div></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
