package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/regcomp/gdpr/caching"
	"github.com/regcomp/gdpr/constants"
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

func VerifyAuthRetryIsRunning(requestStore caching.IRequestStore) func(http.Handler) http.Handler {
	return verifyServiceWorkerIsRunning(
		constants.WorkerAuthRetryPath,
		constants.WorkerAuthRetryScope,
		constants.HeaderAuthRetryWorkerRunning,
		requestStore,
	)
}

func verifyServiceWorkerIsRunning(
	workerPath, workerScope, workerHeader string,
	requestStore caching.IRequestStore,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.RT.UpdateRequestTrace(r, "VerifyServiceWorkerIsRunning")

			if r.Header.Get(workerHeader) != constants.ValueTrueString {
				log.Printf("HEADER=%s, VALUE=%s\n", workerHeader, r.Header.Get(workerHeader))
				requestID, err := requestStore.StoreRequest(r)
				if err != nil {
					log.Panicf("could not cache request, err=%s", err.Error())
				}

				// construct registration url
				registrationURL := fmt.Sprintf("%s?%s=%s&%s=%s&%s=%s",
					constants.EndpointRegisterServiceWorker,
					constants.QueryParamRequestID, requestID,
					constants.QueryParamWorkerPath, workerPath,
					constants.QueryParamWorkerScope, workerScope,
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
	return ScopeServiceWorkerAccess(constants.WorkerAuthRetryPath, constants.WorkerAuthRetryScope)
}

func ScopeServiceWorkerAccess(swPath, accessPath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.RT.UpdateRequestTrace(r, "ScopeServiceWorkerAccess")
			if r.URL.Path == swPath {
				w.Header().Add(constants.HeaderServiceWorkerAllowed, accessPath)
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
		cw := &logging.CustomWriter{ResponseWriter: w, Code: http.StatusOK}
		logging.RT.NewRequestTrace(r)
		next.ServeHTTP(cw, r)
		logging.RT.DumpRequestTrace(r)
		if logging.RT.DisplayResponses() == true {
			fmt.Printf("[RESPONSE HEADER %d]\n", cw.Code)
			cw.Header().Write(os.Stdout)
			fmt.Printf("[BODY]\n%s\n", cw.Body.String())
			fmt.Println("")
		}
	})
}
