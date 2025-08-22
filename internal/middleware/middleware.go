package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/regcomp/gdpr/internal/caching"
	"github.com/regcomp/gdpr/internal/config"
	"github.com/regcomp/gdpr/internal/logging"
	"github.com/regcomp/gdpr/pkg/helpers"
	pkglogging "github.com/regcomp/gdpr/pkg/logging"
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

func VerifyAuthRetryIsRunning(rm *caching.RequestManager) func(http.Handler) http.Handler {
	return verifyServiceWorkerIsRunning(
		config.WorkerAuthRetryPath,
		config.WorkerAuthRetryScope,
		config.HeaderAuthRetryWorkerRunning,
		rm,
		[]string{"/static", "/favicon.ico", "/healthz"},
	)
}

func verifyServiceWorkerIsRunning(
	workerPath, workerScope, workerHeader string,
	rm *caching.RequestManager,
	exempt []string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pkglogging.RT.UpdateRequestTrace(r, "VerifyServiceWorkerIsRunning")

			for _, prefix := range exempt {
				if strings.HasPrefix(r.URL.Path, prefix) {
					next.ServeHTTP(w, r)
					return
				}
			}

			if r.Header.Get(workerHeader) != config.ValueTrue {
				requestID, err := rm.StashRequest(r)
				if err != nil {
					helpers.RespondWithError(w, err, http.StatusInternalServerError)
					return
				}

				registrationURL := fmt.Sprintf("%s?%s=%s&%s=%s&%s=%s",
					config.PathBaseRegisterServiceWorker,
					config.QueryParamRequestId, requestID,
					config.QueryParamWorkerPath, workerPath,
					config.QueryParamWorkerScope, workerScope,
				)

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

func scopeServiceWorkerAccess(sourcePath, accessPath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pkglogging.RT.UpdateRequestTrace(r, "ScopeServiceWorkerAccess")
			if r.URL.Path == sourcePath {
				w.Header().Add(config.HeaderServiceWorkerAllowed, accessPath)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func SetHSTSPolicy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pkglogging.RT.UpdateRequestTrace(r, "SetHSTSPolicy")
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		next.ServeHTTP(w, r)
	})
}
