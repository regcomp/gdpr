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

// fastest code: https://dev.to/chigbeef_77/bool-int-but-stupid-in-go-3jb3
func Btoi(b bool) int {
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}
