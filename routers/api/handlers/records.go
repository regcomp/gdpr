package handlers

import (
	"net/http"

	"github.com/regcomp/gdpr/database"
)

func GetRecordsWithPagination(dbStore database.DatabaseStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get query params

		// make database call

		// create the json object

		// respond with that json
	})
}
