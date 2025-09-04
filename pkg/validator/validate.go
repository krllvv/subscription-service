package validator

import (
	"strconv"
	"strings"
	"subscription-service/internal/model"

	"github.com/google/uuid"
)

func ValidateSubRequest(req model.SubRequest) []string {
	var errors []string

	if req.ServiceName == "" {
		errors = append(errors, "name is required")
	}

	if req.Price <= 0 {
		errors = append(errors, "price must be positive")
	}

	if req.UserID == uuid.Nil {
		errors = append(errors, "user_id is required")
	}

	if req.StartDate == "" {
		errors = append(errors, "start_date is required")
	} else if !ValidateMonthYear(req.StartDate) {
		errors = append(errors, "start_date has invalid format, must be 'MM-YYYY'")
	}

	if req.EndDate != nil && *req.EndDate != "" {
		if !ValidateMonthYear(*req.EndDate) {
			errors = append(errors, "end_date has invalid format, must be 'MM-YYYY'")
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func ValidateMonthYear(dateStr string) bool {
	parts := strings.Split(dateStr, "-")
	if len(parts) != 2 {
		return false
	}

	month, year := parts[0], parts[1]

	if len(month) != 2 {
		return false
	}
	m, err := strconv.Atoi(month)
	if err != nil {
		return false
	}
	if m < 1 || m > 12 {
		return false
	}

	if len(year) != 4 {
		return false
	}
	y, err := strconv.Atoi(year)
	if err != nil {
		return false
	}
	if y < 2000 {
		return false
	}

	return true
}
