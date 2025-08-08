package helpers

import (
	"net/http"
)

func RespondWithError(w http.ResponseWriter, err error, code int) {
	http.Error(w,
		err.Error(),
		code,
	)
}
