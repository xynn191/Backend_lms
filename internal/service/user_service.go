package service

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/xynn191/Backend_lms/internal/models"
	"github.com/xynn191/Backend_lms/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type CreateAdminDTO struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number"`
}

type UpdateAdminDTO struct {
	Email       string `json:"email" binding:"omitempty,email"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number"`
}

type CreateStudentDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	NISN     string `json:"nisn" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Gender   string `json:"gender" binding:"required,oneof=L P"`
	ClassID  uint   `json:"class_id" binding:"required"`
}

type UpdateStudentDTO struct {
	Email   string `json:"email" binding:"omitempty,email"`
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
	Email       string `json:"email" binding:"omitempty,email"`
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
	Email       string `json:"email" binding:"omitempty,email"`
	NIP         string `json:"nip"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number"`
}

type UserService interface {
	// Admins
	GetAllAdmins() ([]models.Admin, error)
	FindAdminByID(id uint) (*models.Admin, error)
	CreateAdmin(dto *CreateAdminDTO) (*models.Admin, error)
	UpdateAdmin(id uint, dto *UpdateAdminDTO) (*models.Admin, error)
	DeleteAdmin(id uint) error

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

	// Superadmin Dashboard
	GetSuperadminDashboardStats() (*repository.SuperadminStats, error)
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

func (s *userService) newUser(roleName, email, password string) (*models.User, error) {
	trimmedEmail := strings.TrimSpace(strings.ToLower(email))
	if available, err := s.repo.IsEmailAvailable(trimmedEmail); err != nil {
		return nil, err
	} else if !available {
		return nil, fmt.Errorf("email %q sudah digunakan oleh akun aktif lain", trimmedEmail)
	}
	hashed, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	role, err := s.repo.FindRoleByName(roleName)
	if err != nil {
		return nil, fmt.Errorf("role %q tidak tersedia: %w", roleName, err)
	}
	return &models.User{
		ID:       uuid.New().String(),
		RoleID:   role.ID,
		Email:    trimmedEmail,
		Password: hashed,
		IsActive: true,
	}, nil
}

// ── Students ──────────────────────────────────────────────────────────────────

func (s *userService) GetAllStudents() ([]models.Student, error) {
	return s.repo.GetAllStudents()
}

func (s *userService) FindStudentByID(id uint) (*models.Student, error) {
	return s.repo.FindStudentByID(id)
}

func (s *userService) CreateStudent(dto *CreateStudentDTO) (*models.Student, error) {
	if exists, err := s.repo.ClassExists(dto.ClassID); err != nil {
		return nil, err
	} else if !exists {
		return nil, fmt.Errorf("kelas dengan ID %d tidak ditemukan di sistem", dto.ClassID)
	}

	if available, err := s.repo.IsNISNAvailable(dto.NISN); err != nil {
		return nil, err
	} else if !available {
		return nil, fmt.Errorf("NISN %q sudah terdaftar pada siswa aktif lain", dto.NISN)
	}
	user, err := s.newUser(models.RoleStudent, dto.Email, dto.Password)
	if err != nil {
		return nil, err
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
	if exists, err := s.repo.ClassExists(dto.ClassID); err != nil {
		return nil, err
	} else if !exists {
		return nil, fmt.Errorf("kelas dengan ID %d tidak ditemukan di sistem", dto.ClassID)
	}

	if dto.NISN != student.NISN {
		if available, err := s.repo.IsNISNAvailable(dto.NISN); err != nil {
			return nil, err
		} else if !available {
			return nil, fmt.Errorf("NISN %q sudah terdaftar pada siswa aktif lain", dto.NISN)
		}
	}

	trimmedEmail := strings.TrimSpace(strings.ToLower(dto.Email))
	if trimmedEmail != "" && trimmedEmail != strings.ToLower(student.User.Email) {
		if available, err := s.repo.IsEmailAvailableForUser(trimmedEmail, student.UserID); err != nil {
			return nil, err
		} else if !available {
			return nil, fmt.Errorf("email %q sudah digunakan oleh akun aktif lain", trimmedEmail)
		}
		if student.UserID != "" {
			if err := s.repo.UpdateUserEmail(student.UserID, trimmedEmail); err != nil {
				return nil, fmt.Errorf("gagal memperbarui email akun: %w", err)
			}
		}
	}

	student.NISN = dto.NISN
	student.Name = dto.Name
	student.Gender = dto.Gender
	student.ClassID = dto.ClassID
	if err := s.repo.UpdateStudent(student); err != nil {
		return nil, err
	}
	return s.repo.FindStudentByID(id)
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
	if dto.NIP != "" {
		if available, err := s.repo.IsNIPAvailable(dto.NIP); err != nil {
			return nil, err
		} else if !available {
			return nil, fmt.Errorf("NIP %q sudah terdaftar pada pengguna aktif lain", dto.NIP)
		}
	}
	user, err := s.newUser(models.RoleTeacher, dto.Email, dto.Password)
	if err != nil {
		return nil, err
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
	if dto.NIP != "" && dto.NIP != teacher.NIP {
		if available, err := s.repo.IsNIPAvailable(dto.NIP); err != nil {
			return nil, err
		} else if !available {
			return nil, fmt.Errorf("NIP %q sudah terdaftar pada pengguna aktif lain", dto.NIP)
		}
	}

	trimmedEmail := strings.TrimSpace(strings.ToLower(dto.Email))
	if trimmedEmail != "" && trimmedEmail != strings.ToLower(teacher.User.Email) {
		if available, err := s.repo.IsEmailAvailableForUser(trimmedEmail, teacher.UserID); err != nil {
			return nil, err
		} else if !available {
			return nil, fmt.Errorf("email %q sudah digunakan oleh akun aktif lain", trimmedEmail)
		}
		if teacher.UserID != "" {
			if err := s.repo.UpdateUserEmail(teacher.UserID, trimmedEmail); err != nil {
				return nil, fmt.Errorf("gagal memperbarui email akun: %w", err)
			}
		}
	}

	teacher.NIP = dto.NIP
	teacher.Name = dto.Name
	teacher.Gender = dto.Gender
	teacher.PhoneNumber = dto.PhoneNumber
	if err := s.repo.UpdateTeacher(teacher); err != nil {
		return nil, err
	}
	return s.repo.FindTeacherByID(id)
}

func (s *userService) DeleteTeacher(id uint) error {
	return s.repo.DeleteTeacher(id)
}

// ── Curriculums ───────────────────────────────────────────────────────────────

func (s *userService) GetAllCurriculums() ([]models.Curriculum, error) {
	return s.repo.GetAllCurriculums()
}

func (s *userService) CreateCurriculum(dto *CreateStaffDTO) (*models.Curriculum, error) {
	if dto.NIP != "" {
		if available, err := s.repo.IsNIPAvailable(dto.NIP); err != nil {
			return nil, err
		} else if !available {
			return nil, fmt.Errorf("NIP %q sudah terdaftar pada pengguna aktif lain", dto.NIP)
		}
	}
	user, err := s.newUser(models.RoleCurriculum, dto.Email, dto.Password)
	if err != nil {
		return nil, err
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
	if dto.NIP != "" && dto.NIP != curriculum.NIP {
		if available, err := s.repo.IsNIPAvailable(dto.NIP); err != nil {
			return nil, err
		} else if !available {
			return nil, fmt.Errorf("NIP %q sudah terdaftar pada pengguna aktif lain", dto.NIP)
		}
	}

	trimmedEmail := strings.TrimSpace(strings.ToLower(dto.Email))
	if trimmedEmail != "" && trimmedEmail != strings.ToLower(curriculum.User.Email) {
		if available, err := s.repo.IsEmailAvailableForUser(trimmedEmail, curriculum.UserID); err != nil {
			return nil, err
		} else if !available {
			return nil, fmt.Errorf("email %q sudah digunakan oleh akun aktif lain", trimmedEmail)
		}
		if curriculum.UserID != "" {
			if err := s.repo.UpdateUserEmail(curriculum.UserID, trimmedEmail); err != nil {
				return nil, fmt.Errorf("gagal memperbarui email akun: %w", err)
			}
		}
	}

	curriculum.NIP = dto.NIP
	curriculum.Name = dto.Name
	curriculum.PhoneNumber = dto.PhoneNumber
	if err := s.repo.UpdateCurriculum(curriculum); err != nil {
		return nil, err
	}
	return s.repo.FindCurriculumByID(id)
}

func (s *userService) DeleteCurriculum(id uint) error {
	return s.repo.DeleteCurriculum(id)
}

// ── Principals ────────────────────────────────────────────────────────────────

func (s *userService) GetAllPrincipals() ([]models.Principal, error) {
	return s.repo.GetAllPrincipals()
}

func (s *userService) CreatePrincipal(dto *CreateStaffDTO) (*models.Principal, error) {
	if dto.NIP != "" {
		if available, err := s.repo.IsNIPAvailable(dto.NIP); err != nil {
			return nil, err
		} else if !available {
			return nil, fmt.Errorf("NIP %q sudah terdaftar pada pengguna aktif lain", dto.NIP)
		}
	}
	user, err := s.newUser(models.RolePrincipal, dto.Email, dto.Password)
	if err != nil {
		return nil, err
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
	if dto.NIP != "" && dto.NIP != principal.NIP {
		if available, err := s.repo.IsNIPAvailable(dto.NIP); err != nil {
			return nil, err
		} else if !available {
			return nil, fmt.Errorf("NIP %q sudah terdaftar pada pengguna aktif lain", dto.NIP)
		}
	}

	trimmedEmail := strings.TrimSpace(strings.ToLower(dto.Email))
	if trimmedEmail != "" && trimmedEmail != strings.ToLower(principal.User.Email) {
		if available, err := s.repo.IsEmailAvailableForUser(trimmedEmail, principal.UserID); err != nil {
			return nil, err
		} else if !available {
			return nil, fmt.Errorf("email %q sudah digunakan oleh akun aktif lain", trimmedEmail)
		}
		if principal.UserID != "" {
			if err := s.repo.UpdateUserEmail(principal.UserID, trimmedEmail); err != nil {
				return nil, fmt.Errorf("gagal memperbarui email akun: %w", err)
			}
		}
	}

	principal.NIP = dto.NIP
	principal.Name = dto.Name
	principal.PhoneNumber = dto.PhoneNumber
	if err := s.repo.UpdatePrincipal(principal); err != nil {
		return nil, err
	}
	return s.repo.FindPrincipalByID(id)
}

func (s *userService) DeletePrincipal(id uint) error {
	return s.repo.DeletePrincipal(id)
}

// ── Admins ───────────────────────────────────────────────────────────────────

func (s *userService) GetAllAdmins() ([]models.Admin, error) {
	return s.repo.GetAllAdmins()
}

func (s *userService) FindAdminByID(id uint) (*models.Admin, error) {
	return s.repo.FindAdminByID(id)
}

func (s *userService) CreateAdmin(dto *CreateAdminDTO) (*models.Admin, error) {
	user, err := s.newUser(models.RoleAdmin, dto.Email, dto.Password)
	if err != nil {
		return nil, err
	}
	admin := &models.Admin{
		Name:        dto.Name,
		PhoneNumber: dto.PhoneNumber,
	}
	if err := s.repo.CreateUserWithAdmin(user, admin); err != nil {
		return nil, err
	}
	return admin, nil
}

func (s *userService) UpdateAdmin(id uint, dto *UpdateAdminDTO) (*models.Admin, error) {
	admin, err := s.repo.FindAdminByID(id)
	if err != nil {
		return nil, err
	}

	trimmedEmail := strings.TrimSpace(strings.ToLower(dto.Email))
	if trimmedEmail != "" && trimmedEmail != strings.ToLower(admin.User.Email) {
		if available, err := s.repo.IsEmailAvailableForUser(trimmedEmail, admin.UserID); err != nil {
			return nil, err
		} else if !available {
			return nil, fmt.Errorf("email %q sudah digunakan oleh akun aktif lain", trimmedEmail)
		}
		if admin.UserID != "" {
			if err := s.repo.UpdateUserEmail(admin.UserID, trimmedEmail); err != nil {
				return nil, fmt.Errorf("gagal memperbarui email akun: %w", err)
			}
		}
	}

	admin.Name = dto.Name
	admin.PhoneNumber = dto.PhoneNumber
	if err := s.repo.UpdateAdmin(admin); err != nil {
		return nil, err
	}
	return s.repo.FindAdminByID(id)
}

func (s *userService) DeleteAdmin(id uint) error {
	return s.repo.DeleteAdmin(id)
}

// ── Superadmin Dashboard Stats ────────────────────────────────────────────────

func (s *userService) GetSuperadminDashboardStats() (*repository.SuperadminStats, error) {
	return s.repo.GetSuperadminDashboardStats()
}
