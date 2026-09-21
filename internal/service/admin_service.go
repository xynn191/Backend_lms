package service

import (
	"github.com/xynn191/Backend_lms/internal/repository"
)

type AdminService interface {
	GetDashboardStats() (*repository.DashboardStats, error)
}

type adminService struct {
	repo repository.AdminRepository
}

func NewAdminService(repo repository.AdminRepository) AdminService {
	return &adminService{repo}
}

func (s *adminService) GetDashboardStats() (*repository.DashboardStats, error) {
	return s.repo.GetDashboardStats()
}
