package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/alienxp03/teya-ledger/types"
)

func (a *APIImpl) respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		a.l.Error("respond", slog.String("error", err.Error()))
	}
}

func (a *APIImpl) respondError(w http.ResponseWriter, status int, err error, msg string) {
	type response struct {
		Message string `json:"message"`
	}

	if serviceError, ok := err.(*types.ServiceError); ok {
		a.respond(w, serviceError.Status, serviceError)
		return
	}
	a.respond(w, status, response{Message: msg})
}

func parseBody(r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("invalid body: %w", err)
	}

	if err := validate.Struct(dst); err != nil {
		return fmt.Errorf("invalid body: %w", err)
	}

	return nil
}
