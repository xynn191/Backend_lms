package repository

import (
	"github.com/xynn191/Backend_lms/internal/models"
	"gorm.io/gorm"
)

type SuperadminStats struct {
	TotalAdmin     int64 `json:"total_admin"`
	TotalSiswa     int64 `json:"total_siswa"`
	TotalGuru      int64 `json:"total_guru"`
	TotalKurikulum int64 `json:"total_kurikulum"`
	TotalKepsek    int64 `json:"total_kepsek"`
	TotalUsers     int64 `json:"total_users"`
}

type UserRepository interface {
	FindRoleByName(name string) (*models.Role, error)
	// Admins
	GetAllAdmins() ([]models.Admin, error)
	FindAdminByID(id uint) (*models.Admin, error)
	CreateUserWithAdmin(user *models.User, admin *models.Admin) error
	UpdateAdmin(admin *models.Admin) error
	DeleteAdmin(id uint) error

	// Students
	GetAllStudents() ([]models.Student, error)
	FindStudentByID(id uint) (*models.Student, error)
	CreateUserWithStudent(user *models.User, student *models.Student) error
	UpdateStudent(student *models.Student) error
	DeleteStudent(id uint) error

	// Teachers
	GetAllTeachers() ([]models.Teacher, error)
	FindTeacherByID(id uint) (*models.Teacher, error)
	CreateUserWithTeacher(user *models.User, teacher *models.Teacher) error
	UpdateTeacher(teacher *models.Teacher) error
	DeleteTeacher(id uint) error

	// Curriculums
	GetAllCurriculums() ([]models.Curriculum, error)
	FindCurriculumByID(id uint) (*models.Curriculum, error)
	CreateUserWithCurriculum(user *models.User, curriculum *models.Curriculum) error
	UpdateCurriculum(curriculum *models.Curriculum) error
	DeleteCurriculum(id uint) error

	// Principals
	GetAllPrincipals() ([]models.Principal, error)
	FindPrincipalByID(id uint) (*models.Principal, error)
	CreateUserWithPrincipal(user *models.User, principal *models.Principal) error
	UpdatePrincipal(principal *models.Principal) error
	DeletePrincipal(id uint) error

	// Active Uniqueness Checks
	UpdateUserEmail(userID string, email string) error
	IsEmailAvailable(email string) (bool, error)
	IsEmailAvailableForUser(email string, excludeUserID string) (bool, error)
	IsNIPAvailable(nip string) (bool, error)
	IsNISNAvailable(nisn string) (bool, error)
	ClassExists(id uint) (bool, error)

	// Superadmin Dashboard
	GetSuperadminDashboardStats() (*SuperadminStats, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindRoleByName(name string) (*models.Role, error) {
	var role models.Role
	if err := r.db.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// ── Admins ───────────────────────────────────────────────────────────────────

func (r *userRepository) GetAllAdmins() ([]models.Admin, error) {
	var admins []models.Admin
	err := r.db.Preload("User").Find(&admins).Error
	return admins, err
}

func (r *userRepository) FindAdminByID(id uint) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.Preload("User").First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *userRepository) CreateUserWithAdmin(user *models.User, admin *models.Admin) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		admin.UserID = user.ID
		return tx.Create(admin).Error
	})
}

func (r *userRepository) UpdateAdmin(admin *models.Admin) error {
	return r.db.Save(admin).Error
}

func (r *userRepository) DeleteAdmin(id uint) error {
	var admin models.Admin
	if err := r.db.First(&admin, id).Error; err != nil {
		return err
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.User{}, "id = ?", admin.UserID).Error; err != nil {
			return err
		}
		return tx.Delete(&admin).Error
	})
}

// ── Students ──────────────────────────────────────────────────────────────────

func (r *userRepository) GetAllStudents() ([]models.Student, error) {
	var students []models.Student
	err := r.db.Preload("Class").Preload("User").Find(&students).Error
	return students, err
}

func (r *userRepository) FindStudentByID(id uint) (*models.Student, error) {
	var student models.Student
	err := r.db.Preload("Class").Preload("User").First(&student, id).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *userRepository) CreateUserWithStudent(user *models.User, student *models.Student) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		student.UserID = user.ID
		return tx.Create(student).Error
	})
}

func (r *userRepository) UpdateStudent(student *models.Student) error {
	return r.db.Save(student).Error
}

func (r *userRepository) DeleteStudent(id uint) error {
	var student models.Student
	if err := r.db.First(&student, id).Error; err != nil {
		return err
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.User{}, "id = ?", student.UserID).Error; err != nil {
			return err
		}
		return tx.Delete(&student).Error
	})
}

// ── Teachers ──────────────────────────────────────────────────────────────────

func (r *userRepository) GetAllTeachers() ([]models.Teacher, error) {
	var teachers []models.Teacher
	err := r.db.Preload("User").Find(&teachers).Error
	return teachers, err
}

func (r *userRepository) FindTeacherByID(id uint) (*models.Teacher, error) {
	var teacher models.Teacher
	err := r.db.Preload("User").First(&teacher, id).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (r *userRepository) CreateUserWithTeacher(user *models.User, teacher *models.Teacher) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		teacher.UserID = user.ID
		return tx.Create(teacher).Error
	})
}

func (r *userRepository) UpdateTeacher(teacher *models.Teacher) error {
	return r.db.Save(teacher).Error
}

func (r *userRepository) DeleteTeacher(id uint) error {
	var teacher models.Teacher
	if err := r.db.First(&teacher, id).Error; err != nil {
		return err
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.User{}, "id = ?", teacher.UserID).Error; err != nil {
			return err
		}
		return tx.Delete(&teacher).Error
	})
}

// ── Curriculums ───────────────────────────────────────────────────────────────

func (r *userRepository) GetAllCurriculums() ([]models.Curriculum, error) {
	var curriculums []models.Curriculum
	err := r.db.Preload("User").Find(&curriculums).Error
	return curriculums, err
}

func (r *userRepository) FindCurriculumByID(id uint) (*models.Curriculum, error) {
	var curriculum models.Curriculum
	err := r.db.Preload("User").First(&curriculum, id).Error
	if err != nil {
		return nil, err
	}
	return &curriculum, nil
}

func (r *userRepository) CreateUserWithCurriculum(user *models.User, curriculum *models.Curriculum) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		curriculum.UserID = user.ID
		return tx.Create(curriculum).Error
	})
}

func (r *userRepository) UpdateCurriculum(curriculum *models.Curriculum) error {
	return r.db.Save(curriculum).Error
}

func (r *userRepository) DeleteCurriculum(id uint) error {
	var curriculum models.Curriculum
	if err := r.db.First(&curriculum, id).Error; err != nil {
		return err
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.User{}, "id = ?", curriculum.UserID).Error; err != nil {
			return err
		}
		return tx.Delete(&curriculum).Error
	})
}

// ── Principals ────────────────────────────────────────────────────────────────

func (r *userRepository) GetAllPrincipals() ([]models.Principal, error) {
	var principals []models.Principal
	err := r.db.Preload("User").Find(&principals).Error
	return principals, err
}

func (r *userRepository) FindPrincipalByID(id uint) (*models.Principal, error) {
	var principal models.Principal
	err := r.db.Preload("User").First(&principal, id).Error
	if err != nil {
		return nil, err
	}
	return &principal, nil
}

func (r *userRepository) CreateUserWithPrincipal(user *models.User, principal *models.Principal) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		principal.UserID = user.ID
		return tx.Create(principal).Error
	})
}

func (r *userRepository) UpdatePrincipal(principal *models.Principal) error {
	return r.db.Save(principal).Error
}

func (r *userRepository) DeletePrincipal(id uint) error {
	var principal models.Principal
	if err := r.db.First(&principal, id).Error; err != nil {
		return err
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.User{}, "id = ?", principal.UserID).Error; err != nil {
			return err
		}
		return tx.Delete(&principal).Error
	})
}

// ── Active Uniqueness Availability Checks ─────────────────────────────────────

func (r *userRepository) UpdateUserEmail(userID string, email string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("email", email).Error
}

func (r *userRepository) IsEmailAvailable(email string) (bool, error) {
	return r.IsEmailAvailableForUser(email, "")
}

func (r *userRepository) IsEmailAvailableForUser(email string, excludeUserID string) (bool, error) {
	var count int64
	query := r.db.Model(&models.User{}).Where("email = ?", email)
	if excludeUserID != "" {
		query = query.Where("id != ?", excludeUserID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *userRepository) IsNIPAvailable(nip string) (bool, error) {
	if nip == "" {
		return true, nil
	}
	var countTeacher, countCurriculum, countPrincipal int64
	_ = r.db.Model(&models.Teacher{}).Where("nip = ?", nip).Count(&countTeacher).Error
	_ = r.db.Model(&models.Curriculum{}).Where("nip = ?", nip).Count(&countCurriculum).Error
	_ = r.db.Model(&models.Principal{}).Where("nip = ?", nip).Count(&countPrincipal).Error
	return (countTeacher + countCurriculum + countPrincipal) == 0, nil
}

func (r *userRepository) IsNISNAvailable(nisn string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Student{}).Where("nisn = ?", nisn).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *userRepository) ClassExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Class{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ── Superadmin Dashboard Stats ────────────────────────────────────────────────

func (r *userRepository) GetSuperadminDashboardStats() (*SuperadminStats, error) {
	var stats SuperadminStats
	if err := r.db.Model(&models.Admin{}).Count(&stats.TotalAdmin).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&models.Student{}).Count(&stats.TotalSiswa).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&models.Teacher{}).Count(&stats.TotalGuru).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&models.Curriculum{}).Count(&stats.TotalKurikulum).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&models.Principal{}).Count(&stats.TotalKepsek).Error; err != nil {
		return nil, err
	}
	stats.TotalUsers = stats.TotalAdmin + stats.TotalSiswa + stats.TotalGuru + stats.TotalKurikulum + stats.TotalKepsek
	return &stats, nil
}
