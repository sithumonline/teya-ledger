package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/alienxp03/teya-ledger/handler/transaction"
)

func (a *APIImpl) setupRoutes() {
	a.mux = http.NewServeMux()

	a.mux.Handle("POST /api/v1/deposits", AuthMiddleware(RequestIDMiddleware(http.HandlerFunc(a.createDeposit))))
	a.mux.Handle("POST /api/v1/withdrawals", AuthMiddleware(RequestIDMiddleware(http.HandlerFunc(a.createWithdrawal))))
	a.mux.Handle("GET /api/v1/balances", AuthMiddleware(RequestIDMiddleware(http.HandlerFunc(a.getBalance))))
	a.mux.Handle("GET /api/v1/transactions", AuthMiddleware(RequestIDMiddleware(http.HandlerFunc(a.getTransactions))))
	a.mux.Handle("GET /api/v1/transactions/{transactionID}", AuthMiddleware(RequestIDMiddleware(http.HandlerFunc(a.getTransaction))))
}

func (a *APIImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.setupRoutes()
	a.l.InfoContext(r.Context(), "Request received", slog.String("method", r.Method), slog.String("path", r.URL.Path))
	a.mux.ServeHTTP(w, r)
}

func (a *APIImpl) createDeposit(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(HeaderUserID).(string)

	params, err := createDepositParams(r)
	if err != nil {
		a.respondError(w, http.StatusBadRequest, err, fmt.Sprintf("Invalid body request %+v", err))
		return
	}

	result, err := a.transactioner.CreateDeposit(userID, *params)
	if err != nil {
		a.respondError(w, http.StatusBadRequest, err, fmt.Sprintf("Failed to create deposit: %+v", err))
		return
	}

	resp := CreateDepositResponse{
		Transaction: Transaction{
			TransactionID: result.Transaction.TransactionID,
			Status:        result.Transaction.Status,
			Amount:        result.Transaction.Amount,
			Currency:      result.Transaction.Currency,
			Description:   result.Transaction.Description,
			CreatedAt:     result.Transaction.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     result.Transaction.UpdatedAt.Format(time.RFC3339),
		},
	}

	a.respond(w, http.StatusOK, resp)
}

func (a *APIImpl) createWithdrawal(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(HeaderUserID).(string)

	params, err := createWithdrawalParams(r)
	if err != nil {
		a.respondError(w, http.StatusBadRequest, err, fmt.Sprintf("Invalid body request %+v", err))
		return
	}

	result, err := a.transactioner.CreateWithdrawal(userID, *params)
	if err != nil {
		a.respondError(w, http.StatusBadRequest, err, fmt.Sprintf("Failed to create withdrawal: %+v", err))
		return
	}

	resp := CreateWithdrawalResponse{
		Transaction: Transaction{
			TransactionID: result.Transaction.TransactionID,
			Status:        result.Transaction.Status,
			Amount:        result.Transaction.Amount,
			Currency:      result.Transaction.Currency,
			Description:   result.Transaction.Description,
			CreatedAt:     result.Transaction.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     result.Transaction.UpdatedAt.Format(time.RFC3339),
		},
	}

	a.respond(w, http.StatusOK, resp)
}

func (a *APIImpl) getTransactions(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(HeaderUserID).(string)
	params := getTransactionsParams(r)

	result, err := a.transactioner.GetTransactions(userID, *params)
	if err != nil {
		a.respondError(w, http.StatusBadRequest, err, fmt.Sprintf("Failed to get transactions %+v", err))
		return
	}

	resp := GetTransactionsResponse{}
	transactions := []Transaction{}
	for _, transaction := range result.Transactions {
		transactions = append(transactions, Transaction{
			TransactionID: transaction.TransactionID,
			Status:        transaction.Status,
			Amount:        transaction.Amount,
			Currency:      transaction.Currency,
			Description:   transaction.Description,
			CreatedAt:     transaction.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     transaction.UpdatedAt.Format(time.RFC3339),
		})
	}

	resp.Transactions = transactions

	a.respond(w, http.StatusOK, resp)
}

func (a *APIImpl) getBalance(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(HeaderUserID).(string)

	req := getBalancesParams(r)

	resp, err := a.transactioner.GetBalance(userID, *req)
	if err != nil {
		a.respondError(w, http.StatusBadRequest, err, fmt.Sprintf("Failed to get balance: %+v", err))
		return
	}

	result := GetBalanceResponse{
		Balance: Balance{Amount: resp.Amount, Currency: resp.Currency},
	}

	a.respond(w, http.StatusOK, result)
}

func getTransactionsParams(r *http.Request) *transaction.GetTransactionsRequest {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 0
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 0
	}

	req := &transaction.GetTransactionsRequest{
		AccountNumber: r.URL.Query().Get("accountNumber"),
		Limit:         limit,
		Page:          page,
	}

	return req
}

func createDepositParams(r *http.Request) (*transaction.CreateDepositRequest, error) {
	var req CreateDepositRequest
	if err := parseBody(r, &req); err != nil {
		return nil, err
	}

	result := &transaction.CreateDepositRequest{
		TransactionID: req.TransactionID,
		AccountNumber: req.AccountNumber,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Description:   req.Description,
	}
	return result, nil
}

func createWithdrawalParams(r *http.Request) (*transaction.CreateWithdrawalRequest, error) {
	var req CreateWithdrawalRequest
	if err := parseBody(r, &req); err != nil {
		return nil, err
	}

	result := &transaction.CreateWithdrawalRequest{
		TransactionID: req.TransactionID,
		AccountNumber: req.AccountNumber,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Description:   req.Description,
	}
	return result, nil
}

func getBalancesParams(r *http.Request) *transaction.GetBalanceRequest {
	req := &transaction.GetBalanceRequest{
		AccountNumber: r.URL.Query().Get("accountNumber"),
	}

	return req
}

func (a *APIImpl) getTransaction(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(HeaderUserID).(string)
	transactionID := path.Base(r.URL.Path)

	transaction, err := a.transactioner.GetTransaction(userID, transactionID)
	if err != nil {
		a.respondError(w, http.StatusBadRequest, err, fmt.Sprintf("Failed to get transaction: %+v", err))
		return
	}
	a.l.InfoContext(r.Context(), "Request received", slog.String("transaction", fmt.Sprintf("%+v", transaction)))

	result := GetTransactionResponse{
		Transaction: Transaction{
			TransactionID: transaction.TransactionID,
			Status:        transaction.Status,
			Amount:        transaction.Amount,
			Currency:      transaction.Currency,
			Description:   transaction.Description,
			CreatedAt:     transaction.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     transaction.UpdatedAt.Format(time.RFC3339),
		},
	}

	a.respond(w, http.StatusOK, result)
}
