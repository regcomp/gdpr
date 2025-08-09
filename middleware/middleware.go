package middleware

import (
	"fmt"
	"net/http"

	"github.com/regcomp/gdpr/caching"
	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/helpers"
	"github.com/regcomp/gdpr/logging"
)

func RequestLogging(requestLogger logging.ILogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/favicon.ico" {
				requestLogger.Log(r)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func VerifyAuthRetryIsRunning(requestStash *caching.RequestStash) func(http.Handler) http.Handler {
	return verifyServiceWorkerIsRunning(
		config.WorkerAuthRetryPath,
		config.WorkerAuthRetryScope,
		config.HeaderAuthRetryWorkerRunning,
		requestStash,
	)
}

func verifyServiceWorkerIsRunning(
	workerPath, workerScope, workerHeader string,
	requestStash *caching.RequestStash,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.RT.UpdateRequestTrace(r, "VerifyServiceWorkerIsRunning")

			if r.Header.Get(workerHeader) != config.ValueTrue {
				// log.Printf("HEADER=%s, VALUE=%s\n", workerHeader, r.Header.Get(workerHeader))
				requestID, err := requestStash.StashRequest(r)
				if err != nil {
					helpers.RespondWithError(w, err, http.StatusInternalServerError)
					return
				}

				// construct registration url
				registrationURL := fmt.Sprintf("%s?%s=%s&%s=%s&%s=%s",
					config.PathBaseRegisterServiceWorker,
					config.QueryParamRequestId, requestID,
					config.QueryParamWorkerPath, workerPath,
					config.QueryParamWorkerScope, workerScope,
				)

				// redirect to the url
				http.Redirect(w, r, registrationURL, http.StatusTemporaryRedirect)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func ScopeAuthRetryAccess() func(http.Handler) http.Handler {
	return scopeServiceWorkerAccess(config.WorkerAuthRetryPath, config.WorkerAuthRetryScope)
}

func scopeServiceWorkerAccess(swPath, accessPath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.RT.UpdateRequestTrace(r, "ScopeServiceWorkerAccess")
			if r.URL.Path == swPath {
				w.Header().Add(config.HeaderServiceWorkerAllowed, accessPath)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func SetHSTSPolicy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "SetHSTSPolicy")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		next.ServeHTTP(w, r)
	})
}

func TraceRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cw := logging.CreateCustomWriter(w)
		logging.RT.NewRequestTrace(cw, r)
		next.ServeHTTP(cw, r)
		logging.RT.DumpRequestTrace(r)
	})
}
