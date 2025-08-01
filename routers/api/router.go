package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/routers/api/handlers"
	servicecontext "github.com/regcomp/gdpr/service_context"
)

func CreateApiRouter(stx *servicecontext.ServiceContext) *chi.Mux {
	mux := chi.NewRouter()

	mux.Get(config.EndpointRecords, handlers.GetRecordsWithPagination(stx.DatabaseStore))

	return mux
}
