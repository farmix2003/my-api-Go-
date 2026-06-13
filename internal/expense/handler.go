package expense

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	var req Expense
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	if err := h.service.CreateExpense(r.Context(), &req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, req)
}

func (h *Handler) GetAllExpensesHandler(w http.ResponseWriter, r *http.Request) {
	expenses, err := h.service.GetAllExpenses(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve expenses")
		return
	}

	writeJSON(w, http.StatusOK, expenses)
}

func (h *Handler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	var req Expense

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid expense id")
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	if err := h.service.UpdateExpense(r.Context(), id, &req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, req)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
