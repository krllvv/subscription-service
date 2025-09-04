package service

import (
	"subscription-service/internal/model"
	"subscription-service/internal/repository/sub"

	"github.com/google/uuid"
)

type SubService struct {
	repo sub.SubscriptionRepository
}

func NewSubService(repository sub.SubscriptionRepository) *SubService {
	return &SubService{repo: repository}
}

func (s *SubService) Create(sub *model.Subscription) error {
	return s.repo.Create(sub)
}

func (s *SubService) GetByID(id uuid.UUID) (*model.Subscription, error) {
	return s.repo.GetByID(id)
}

func (s *SubService) Update(id uuid.UUID, sub *model.Subscription) (*model.Subscription, error) {
	err := s.repo.Update(id, sub)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}

func (s *SubService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *SubService) GetAll() ([]model.Subscription, error) {
	return s.repo.GetAll()
}

func (s *SubService) GetTotalSum(start, end string, userID uuid.UUID, name string) (int, error) {
	return s.repo.GetTotalSum(start, end, userID, name)
}
