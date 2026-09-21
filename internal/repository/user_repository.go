package repository

import (
	"github.com/xynn191/Backend_lms/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
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
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
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
