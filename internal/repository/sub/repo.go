package sub

import (
	"subscription-service/internal/model"

	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(sub *model.Subscription) error
	GetByID(id uuid.UUID) (*model.Subscription, error)
	Update(id uuid.UUID, sub *model.Subscription) error
	Delete(id uuid.UUID) error
	GetAll() ([]model.Subscription, error)
	GetTotalSum(start, end string, userID uuid.UUID, name string) (int, error)
}
