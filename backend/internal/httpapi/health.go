package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

// Health はサーバーが応答可能であることを返す。
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(healthResponse{Status: "ok"}); err != nil {
		slog.ErrorContext(r.Context(), "failed to write health response", "error", err)
	}
}
