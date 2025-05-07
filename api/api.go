package api

import (
	"log/slog"
	"net/http"

	"github.com/alienxp03/teya-ledger/handler/transaction"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

const (
	HeaderUserID = "userID"
)

type APIImpl struct {
	transactioner transaction.Transactioner
	l             *slog.Logger

	mux *http.ServeMux
}

func New(transactioner transaction.Transactioner, l *slog.Logger) *APIImpl {
	return &APIImpl{
		transactioner: transactioner,
		l:             l,
	}
}
