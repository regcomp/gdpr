package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/regcomp/gdpr/internal/config"
	"github.com/regcomp/gdpr/internal/routers/api/handlers"
	"github.com/regcomp/gdpr/internal/servicecontext"
)

func CreateApiRouter(stx *servicecontext.ServiceContext) *chi.Mux {
	mux := chi.NewRouter()

	mux.Get(config.EndpointApiRecords, handlers.GetRecordsWithPagination(stx.DatabaseManager))

	return mux
}
