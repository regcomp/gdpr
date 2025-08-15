package helpers

import (
	"net/http"
)

func RespondWithError(w http.ResponseWriter, err error, code int) {
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
