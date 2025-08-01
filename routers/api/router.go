package api

import (
	"github.com/go-chi/chi/v5"
	servicecontext "github.com/regcomp/gdpr/service_context"
)

func CreateApiRouter(stx *servicecontext.ServiceContext) *chi.Mux {
	mux := chi.NewRouter()

	return mux
}
