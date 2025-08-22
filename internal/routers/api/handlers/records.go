package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/regcomp/gdpr/internal/config"
	db "github.com/regcomp/gdpr/internal/database"
	"github.com/regcomp/gdpr/pkg/helpers"
	"github.com/regcomp/gdpr/pkg/logging"
)

const (
	paginationMin = 1
	paginationMax = 100
)

func GetRecordsWithPagination(dbm *db.DatabaseManager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.RT.UpdateRequestTrace(r, "GetRecordsWithPagination")
		// Get query params
		paginationLimitString := r.URL.Query().Get(config.QueryParamLimit)
		if paginationLimitString == "" {
			helpers.RespondWithError(w, fmt.Errorf("missing pagination limit"), http.StatusBadRequest)
			return
		}
		paginationLimit, err := strconv.Atoi(paginationLimitString)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusBadRequest)
			return
		}
		if paginationLimit < paginationMin || paginationLimit > paginationMax {
			helpers.RespondWithError(
				w,
				fmt.Errorf("invalid pagination limit=%d, must be between %d and %d",
					paginationLimit, paginationMin, paginationMax),
				http.StatusBadRequest,
			)
			return
		}

		var queryStart time.Time
		queryStartString := r.URL.Query().Get(config.QueryParamAfter)
		if queryStartString == "" {
			helpers.RespondWithError(w, fmt.Errorf("missing next cursor"), http.StatusBadRequest)
			return
		} else {
			queryStart, err = time.Parse(time.RFC3339, queryStartString)
			if err != nil {
				helpers.RespondWithError(w, err, http.StatusBadRequest)
				return
			}
		}

		// make database call
		records, paginationInfo, err := dbm.GetDeletionRequestsAndPaginationInfo(paginationLimit, queryStart)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		responseData := struct {
			Data       []db.RecordOfDeletionRequest `json:"data"`
			Pagination db.PaginationInfo            `json:"pagination"`
		}{
			Data:       records,
			Pagination: paginationInfo,
		}

		// create the json object
		data, err := json.Marshal(responseData)
		if err != nil {
			helpers.RespondWithError(w, err, http.StatusInternalServerError)
		}
		// respond with that json
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		w.Write(data)
	})
}
