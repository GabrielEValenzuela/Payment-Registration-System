package bank

import "github.com/GabrielEValenzuela/Payment-Registration-System/src/internal/models"

type Service interface {
	AddFinancingPromotionToBank(promotionFinancing models.Financing) error
}

type service struct {
	repo IRepository
}

func NewService(repo IRepository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) AddFinancingPromotionToBank(promotionFinancing models.Financing) error {
	return s.repo.AddFinancingPromotionToBank(promotionFinancing)
}
