package handlers

import (
	"net/http"
)

func (stx *ServiceContext) GetDashboard(w http.ResponseWriter, r *http.Request) {
	stx.Logger.Info("Dashboard", nil)
}
