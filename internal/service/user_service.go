package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/xynn191/Backend_lms/internal/models"
	"github.com/xynn191/Backend_lms/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type CreateStudentDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	NISN     string `json:"nisn" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Gender   string `json:"gender" binding:"required,oneof=L P"`
	ClassID  uint   `json:"class_id" binding:"required"`
}

type UpdateStudentDTO struct {
	NISN    string `json:"nisn" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Gender  string `json:"gender" binding:"required,oneof=L P"`
	ClassID uint   `json:"class_id" binding:"required"`
}

type CreateTeacherDTO struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	NIP         string `json:"nip"`
	Name        string `json:"name" binding:"required"`
	Gender      string `json:"gender" binding:"required,oneof=L P"`
	PhoneNumber string `json:"phone_number"`
}

type UpdateTeacherDTO struct {
	NIP         string `json:"nip"`
	Name        string `json:"name" binding:"required"`
	Gender      string `json:"gender" binding:"required,oneof=L P"`
	PhoneNumber string `json:"phone_number"`
}

type CreateStaffDTO struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	NIP         string `json:"nip"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number"`
}

type UpdateStaffDTO struct {
	NIP         string `json:"nip"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number"`
}

type UserService interface {
	// Students
	GetAllStudents() ([]models.Student, error)
	FindStudentByID(id uint) (*models.Student, error)
	CreateStudent(dto *CreateStudentDTO) (*models.Student, error)
	UpdateStudent(id uint, dto *UpdateStudentDTO) (*models.Student, error)
	DeleteStudent(id uint) error

	// Teachers
	GetAllTeachers() ([]models.Teacher, error)
	FindTeacherByID(id uint) (*models.Teacher, error)
	CreateTeacher(dto *CreateTeacherDTO) (*models.Teacher, error)
	UpdateTeacher(id uint, dto *UpdateTeacherDTO) (*models.Teacher, error)
	DeleteTeacher(id uint) error

	// Curriculums
	GetAllCurriculums() ([]models.Curriculum, error)
	CreateCurriculum(dto *CreateStaffDTO) (*models.Curriculum, error)
	UpdateCurriculum(id uint, dto *UpdateStaffDTO) (*models.Curriculum, error)
	DeleteCurriculum(id uint) error

	// Principals
	GetAllPrincipals() ([]models.Principal, error)
	CreatePrincipal(dto *CreateStaffDTO) (*models.Principal, error)
	UpdatePrincipal(id uint, dto *UpdateStaffDTO) (*models.Principal, error)
	DeletePrincipal(id uint) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("gagal hash password: %w", err)
	}
	return string(hashed), nil
}

// ── Students ──────────────────────────────────────────────────────────────────

func (s *userService) GetAllStudents() ([]models.Student, error) {
	return s.repo.GetAllStudents()
}

func (s *userService) FindStudentByID(id uint) (*models.Student, error) {
	return s.repo.FindStudentByID(id)
}

func (s *userService) CreateStudent(dto *CreateStudentDTO) (*models.Student, error) {
	hashed, err := hashPassword(dto.Password)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		ID:       uuid.New().String(),
		RoleID:   3, // siswa
		Email:    dto.Email,
		Password: hashed,
		IsActive: true,
	}
	student := &models.Student{
		NISN:    dto.NISN,
		Name:    dto.Name,
		Gender:  dto.Gender,
		ClassID: dto.ClassID,
	}
	if err := s.repo.CreateUserWithStudent(user, student); err != nil {
		return nil, err
	}
	return student, nil
}

func (s *userService) UpdateStudent(id uint, dto *UpdateStudentDTO) (*models.Student, error) {
	student, err := s.repo.FindStudentByID(id)
	if err != nil {
		return nil, err
	}
	student.NISN = dto.NISN
	student.Name = dto.Name
	student.Gender = dto.Gender
	student.ClassID = dto.ClassID
	if err := s.repo.UpdateStudent(student); err != nil {
		return nil, err
	}
	return student, nil
}

func (s *userService) DeleteStudent(id uint) error {
	return s.repo.DeleteStudent(id)
}

// ── Teachers ──────────────────────────────────────────────────────────────────

func (s *userService) GetAllTeachers() ([]models.Teacher, error) {
	return s.repo.GetAllTeachers()
}

func (s *userService) FindTeacherByID(id uint) (*models.Teacher, error) {
	return s.repo.FindTeacherByID(id)
}

func (s *userService) CreateTeacher(dto *CreateTeacherDTO) (*models.Teacher, error) {
	hashed, err := hashPassword(dto.Password)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		ID:       uuid.New().String(),
		RoleID:   2, // guru
		Email:    dto.Email,
		Password: hashed,
		IsActive: true,
	}
	teacher := &models.Teacher{
		NIP:         dto.NIP,
		Name:        dto.Name,
		Gender:      dto.Gender,
		PhoneNumber: dto.PhoneNumber,
	}
	if err := s.repo.CreateUserWithTeacher(user, teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func (s *userService) UpdateTeacher(id uint, dto *UpdateTeacherDTO) (*models.Teacher, error) {
	teacher, err := s.repo.FindTeacherByID(id)
	if err != nil {
		return nil, err
	}
	teacher.NIP = dto.NIP
	teacher.Name = dto.Name
	teacher.Gender = dto.Gender
	teacher.PhoneNumber = dto.PhoneNumber
	if err := s.repo.UpdateTeacher(teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func (s *userService) DeleteTeacher(id uint) error {
	return s.repo.DeleteTeacher(id)
}

// ── Curriculums ───────────────────────────────────────────────────────────────

func (s *userService) GetAllCurriculums() ([]models.Curriculum, error) {
	return s.repo.GetAllCurriculums()
}

func (s *userService) CreateCurriculum(dto *CreateStaffDTO) (*models.Curriculum, error) {
	hashed, err := hashPassword(dto.Password)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		ID:       uuid.New().String(),
		RoleID:   4, // kurikulum
		Email:    dto.Email,
		Password: hashed,
		IsActive: true,
	}
	curriculum := &models.Curriculum{
		NIP:         dto.NIP,
		Name:        dto.Name,
		PhoneNumber: dto.PhoneNumber,
	}
	if err := s.repo.CreateUserWithCurriculum(user, curriculum); err != nil {
		return nil, err
	}
	return curriculum, nil
}

func (s *userService) UpdateCurriculum(id uint, dto *UpdateStaffDTO) (*models.Curriculum, error) {
	curriculum, err := s.repo.FindCurriculumByID(id)
	if err != nil {
		return nil, err
	}
	curriculum.NIP = dto.NIP
	curriculum.Name = dto.Name
	curriculum.PhoneNumber = dto.PhoneNumber
	if err := s.repo.UpdateCurriculum(curriculum); err != nil {
		return nil, err
	}
	return curriculum, nil
}

func (s *userService) DeleteCurriculum(id uint) error {
	return s.repo.DeleteCurriculum(id)
}

// ── Principals ────────────────────────────────────────────────────────────────

func (s *userService) GetAllPrincipals() ([]models.Principal, error) {
	return s.repo.GetAllPrincipals()
}

func (s *userService) CreatePrincipal(dto *CreateStaffDTO) (*models.Principal, error) {
	hashed, err := hashPassword(dto.Password)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		ID:       uuid.New().String(),
		RoleID:   5, // kepala_sekolah
		Email:    dto.Email,
		Password: hashed,
		IsActive: true,
	}
	principal := &models.Principal{
		NIP:         dto.NIP,
		Name:        dto.Name,
		PhoneNumber: dto.PhoneNumber,
	}
	if err := s.repo.CreateUserWithPrincipal(user, principal); err != nil {
		return nil, err
	}
	return principal, nil
}

func (s *userService) UpdatePrincipal(id uint, dto *UpdateStaffDTO) (*models.Principal, error) {
	principal, err := s.repo.FindPrincipalByID(id)
	if err != nil {
		return nil, err
	}
	principal.NIP = dto.NIP
	principal.Name = dto.Name
	principal.PhoneNumber = dto.PhoneNumber
	if err := s.repo.UpdatePrincipal(principal); err != nil {
		return nil, err
	}
	return principal, nil
}

func (s *userService) DeletePrincipal(id uint) error {
	return s.repo.DeletePrincipal(id)
}
