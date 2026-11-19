package service

import "github.com/xynn191/Backend_lms/internal/repository"

type MonitoringService interface {
	GetCurriculumMonitoring() (*repository.CurriculumMonitoringData, error)
	GetPrincipalMonitoring() (*repository.PrincipalMonitoringData, error)
	GetSchoolGrades(classID uint, subjectID uint) ([]repository.SchoolGradeItem, error)
}

type monitoringService struct {
	repo repository.MonitoringRepository
}

func NewMonitoringService(repo repository.MonitoringRepository) MonitoringService {
	return &monitoringService{repo: repo}
}

func (s *monitoringService) GetCurriculumMonitoring() (*repository.CurriculumMonitoringData, error) {
	return s.repo.GetCurriculumMonitoring()
}

func (s *monitoringService) GetPrincipalMonitoring() (*repository.PrincipalMonitoringData, error) {
	return s.repo.GetPrincipalMonitoring()
}

func (s *monitoringService) GetSchoolGrades(classID uint, subjectID uint) ([]repository.SchoolGradeItem, error) {
	return s.repo.GetSchoolGrades(classID, subjectID)
}
