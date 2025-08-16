package views

/*

  This package is a thin wrapper around the implementation logic of the templating engine.

*/

import (
	"context"
	"io"

	"github.com/regcomp/gdpr/internal/caching"
	"github.com/regcomp/gdpr/internal/views/templ/components"
	"github.com/regcomp/gdpr/internal/views/templ/pages"
)

func ServeLogin(w io.Writer, ctx context.Context) error {
	return pages.Login().Render(ctx, w)
}

func ServeDashboard(w io.Writer, ctx context.Context,
	accessToken, refreshToken, sessionID string,
) error {
	return pages.Dashboard(accessToken, refreshToken, sessionID).Render(ctx, w)
}

func ServeRegisterServiceWorker(w io.Writer, ctx context.Context,
	cachedRequest *caching.CachedRequest, swPath, swScope string,
) error {
	return pages.RegisterServiceWorker(swPath, swScope, *cachedRequest).Render(ctx, w)
}

func WriteRecordsManagement(w io.Writer, ctx context.Context) error {
	return components.RecordsManagement().Render(ctx, w)
}
