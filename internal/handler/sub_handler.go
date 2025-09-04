package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"subscription-service/internal/model"
	"subscription-service/internal/repository/sub/postgres"
	"subscription-service/internal/service"
	"subscription-service/pkg/utils"
	"subscription-service/pkg/validator"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	paramSubID     = "subID"
	errDecodeMsg   = "decode JSON error"
	errEncodeMsg   = "encode JSON error"
	errInternalMsg = "internal error occurred"
	errInvalidID   = "invalid subscription ID"
	errNotFound    = "subscription not found"
)

type SubHandler struct {
	srv    *service.SubService
	logger *log.Logger
}

func NewSubHandler(srv *service.SubService, logger *log.Logger) *SubHandler {
	return &SubHandler{srv: srv, logger: logger}
}

func (h *SubHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/subscriptions", h.create).Methods("POST")
	r.HandleFunc("/subscriptions", h.getAll).Methods("GET")
	r.HandleFunc("/subscription/{subID}", h.get).Methods("GET")
	r.HandleFunc("/subscription/{subID}", h.update).Methods("PUT")
	r.HandleFunc("/subscription/{subID}", h.delete).Methods("DELETE")
	r.HandleFunc("/subscriptions/total", h.totalSum).Methods("GET")
}

// @Summary		Create Subscription
// @Description	Create a new subscription record. Field 'end_data' is optional
// @Tags		Subscriptions
// @Accept		json
// @Produce		json
// @Param		subscription	body		model.SubRequest	true	"Subscription payload"
// @Success		201				{object}	model.Subscription	"Successfully created subscription"
// @Failure		400				{object}	utils.ErrorResponse	"Validation error or invalid request body"
// @Failure		500				{object}	utils.ErrorResponse	"Internal server error"
// @Router		/subscriptions [post]
func (h *SubHandler) create(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("CREATE subscription request")

	var req model.SubRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Println("Create: decode error:", err)
		utils.WriteError(w, http.StatusInternalServerError, errDecodeMsg)
		return
	}

	if validationErrs := validator.ValidateSubRequest(req); validationErrs != nil {
		h.logger.Println("Create: validation error", validationErrs)
		utils.WriteValidationErrors(w, validationErrs)
		return
	}

	sub := model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	if err := h.srv.Create(&sub); err != nil {
		h.logger.Println("Failed to create subscription:", err)
		utils.WriteError(w, http.StatusInternalServerError, errInternalMsg)
		return
	}

	err := utils.WriteJSON(w, http.StatusCreated, sub)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errEncodeMsg)
	}
}

// @Summary		Get All Subscription
// @Description	Get list of all subscriptions
// @Tags		Subscriptions
// @Produce		json
// @Success		200	{array}		model.Subscription	"A list of subscriptions"
// @Failure		500	{object}	utils.ErrorResponse	"Internal server error"
// @Router		/subscriptions [get]
func (h *SubHandler) getAll(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("GET all subscriptions request")

	subs, err := h.srv.GetAll()
	if err != nil {
		h.logger.Println("Failed to get subscriptions:", err)
		utils.WriteError(w, http.StatusInternalServerError, errInternalMsg)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, subs)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errEncodeMsg)
	}
}

// @Summary		Get Subscription
// @Description	Get subscription by ID
// @Tags		Subscriptions
// @Produce		json
// @Param		subID	path		string				true	"Subscription ID"	format(uuid)
// @Success		200		{object}	model.Subscription	"Requested subscription"
// @Failure		400		{object}	utils.ErrorResponse	"Invalid subscription ID"
// @Failure		404		{object}	utils.ErrorResponse	"Subscription not found"
// @Failure		500		{object}	utils.ErrorResponse	"Internal server error"
// @Router		/subscription/{subID} [get]
func (h *SubHandler) get(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("GET subscription request")

	vars := mux.Vars(r)
	subID := vars[paramSubID]

	id, err := uuid.Parse(subID)
	if err != nil {
		h.logger.Println("Invalid subscription ID:", err)
		utils.WriteError(w, http.StatusBadRequest, errInvalidID)
		return
	}

	sub, err := h.srv.GetByID(id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			h.logger.Println("Get: subscription not found:", err)
			utils.WriteError(w, http.StatusNotFound, errNotFound)
			return
		}
		h.logger.Println("Failed to get subscription:", err)
		utils.WriteError(w, http.StatusInternalServerError, errInternalMsg)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, sub)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errEncodeMsg)
	}
}

// @Summary		Update subscription
// @Description	Update subscription by ID
// @Tags		Subscriptions
// @Accept		json
// @Produce		json
// @Param		subID			path		string				true	"Subscription ID"	format(uuid)
// @Param		subscription	body		model.SubRequest	true	"Updated subscription payload"
// @Success		200				{object}	model.Subscription	"Successfully updated subscription"
// @Failure		400				{object}	utils.ErrorResponse	"Invalid subscription ID or validation error"
// @Failure		404				{object}	utils.ErrorResponse	"Subscription not found"
// @Failure		500				{object}	utils.ErrorResponse	"Internal server error"
// @Router		/subscription/{subID} [put]
func (h *SubHandler) update(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("UPDATE subscription request")

	vars := mux.Vars(r)
	subID := vars[paramSubID]

	id, err := uuid.Parse(subID)
	if err != nil {
		h.logger.Println("Invalid subscription ID:", err)
		utils.WriteError(w, http.StatusBadRequest, errInvalidID)
		return
	}

	var req model.SubRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Println("Create: decode error:", err)
		utils.WriteError(w, http.StatusInternalServerError, errDecodeMsg)
		return
	}

	if validationErrs := validator.ValidateSubRequest(req); validationErrs != nil {
		h.logger.Println("Update: validation error", validationErrs)
		utils.WriteValidationErrors(w, validationErrs)
		return
	}

	sub := model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	newSub, err := h.srv.Update(id, &sub)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			h.logger.Println("Update error, subscription not found:", err)
			utils.WriteError(w, http.StatusNotFound, errNotFound)
			return
		}
		h.logger.Println("Failed to update subscription:", err)
		utils.WriteError(w, http.StatusInternalServerError, errInternalMsg)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, newSub)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errEncodeMsg)
		return
	}
}

// @Summary		Delete Subscription
// @Description	Delete subscription by ID
// @Tags		Subscriptions
// @Produce		json
// @Param		subID	path	string	true	"Subscription ID"	format(uuid)
// @Success		204		"Successfully deleted subscription"
// @Failure		400		{object}	utils.ErrorResponse	"Invalid subscription ID"
// @Failure		404		{object}	utils.ErrorResponse	"Subscription not found"
// @Failure		500		{object}	utils.ErrorResponse	"Internal server error"
// @Router		/subscription/{subID} [delete]
func (h *SubHandler) delete(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("DELETE subscription request")

	vars := mux.Vars(r)
	subID := vars[paramSubID]

	id, err := uuid.Parse(subID)
	if err != nil {
		h.logger.Println("Invalid subscription ID:", err)
		utils.WriteError(w, http.StatusBadRequest, errInvalidID)
		return
	}

	err = h.srv.Delete(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Println("Delete error, subscription not found:", err)
			utils.WriteError(w, http.StatusNotFound, errNotFound)
			return
		}
		h.logger.Println("Failed to delete subscription:", err)
		utils.WriteError(w, http.StatusInternalServerError, errInternalMsg)
		return
	}
}

// @Summary		Calculate Total Sum
// @Description	Get the total cost of subscriptions for a given period. Optional filters for user and subscription name
// @Tags		Subscriptions
// @Produce		json
// @Param		start_date		query		string				true	"Start date of the period (MM-YYYY)"	Example("01-2025")
// @Param		end_date		query		string				true	"End date of the period (MM-YYYY)"		Example("12-2025")
// @Param		user_id			query		string				false	"Filter by User ID (UUID)"				format(uuid)
// @Param		service_name	query		string				false	"Filter by subscription name"
// @Success		200				{object}	map[string]int		"Total sum"
// @Failure		400				{object}	utils.ErrorResponse	"Invalid parameters"
// @Failure		500				{object}	utils.ErrorResponse	"Internal server error"
// @Router		/subscriptions/total [get]
func (h *SubHandler) totalSum(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("GET total sum of subscriptions request")

	params := r.URL.Query()
	startDate := params.Get("start_date")
	endDate := params.Get("end_date")
	userID := params.Get("user_id")
	name := params.Get("service_name")

	var id uuid.UUID
	if userID != "" {
		var err error
		id, err = uuid.Parse(userID)
		if err != nil {
			h.logger.Println("Invalid subscription ID:", err)
			utils.WriteError(w, http.StatusBadRequest, errInvalidID)
			return
		}
	}

	if startDate == "" || endDate == "" {
		h.logger.Println("start_date or end_date is empty")
		utils.WriteError(w, http.StatusBadRequest, "start_date and end_date must be in query")
		return
	}

	if !validator.ValidateMonthYear(startDate) || !validator.ValidateMonthYear(endDate) {
		h.logger.Println("start_date or end_data is incorrect")
		utils.WriteError(w, http.StatusBadRequest, "dates must be valid")
		return
	}

	sum, err := h.srv.GetTotalSum(startDate, endDate, id, name)
	if err != nil {
		h.logger.Println("Failed to get total sum:", err)
		utils.WriteError(w, http.StatusInternalServerError, errInternalMsg)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, map[string]int{"total_sum": sum})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, errEncodeMsg)
		return
	}
}
