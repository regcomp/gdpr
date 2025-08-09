package helpers

import (
	"log"
	"net/http"
	"runtime/debug"
)

func RespondWithError(w http.ResponseWriter, err error, code int) {
	log.Printf("RespondWithError called: %v (code: %d)", err, code)
	log.Printf("Stack trace: %s", debug.Stack())
	http.Error(w, err.Error(), code)
}
