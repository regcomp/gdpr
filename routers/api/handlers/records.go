package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/regcomp/gdpr/config"
	"github.com/regcomp/gdpr/database"
)

func GetRecordsWithPagination(dbStore database.IDatabaseManager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get query params
		paginationLimitString := r.URL.Query().Get(config.QueryParamLimit)
		if paginationLimitString == "" {
			// TODO: FATAL
		}
		paginationLimit, err := strconv.Atoi(paginationLimitString)
		if err != nil {
			// TODO: FATAL
		}
		if paginationLimit < 1 || paginationLimit > 100 {
			// TODO: Bad request
		}

		queryStartString := r.URL.Query().Get(config.QueryParamAfter)
		queryStart, err := time.Parse(time.RFC3339, queryStartString)
		if err != nil {
			// TODO: FATAL
		}

		// make database call
		records, paginationInfo, err := dbStore.GetScheduledDeletionRecords(paginationLimit, queryStart)
		if err != nil {
			// TODO:
		}

		responseData := struct {
			Data       []database.RecordOfDeletionRequest `json:"data"`
			Pagination database.PaginationInfo            `json:"pagination"`
		}{
			Data:       records,
			Pagination: paginationInfo,
		}

		// create the json object
		data, err := json.Marshal(responseData)
		if err != nil {
			// TODO:
		}
		// respond with that json
		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")
		w.Write(data)
	})
}
