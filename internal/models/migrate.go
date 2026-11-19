package models

import (
	"log"

	"gorm.io/gorm"
)

// AutoMigrate all models
func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&Role{},
		&User{},
		&Admin{},
		&Teacher{},
		&Major{},
		&AcademicYear{},
		&Class{},   // Dependent on Major and Teacher
		&Student{}, // Dependent on Class
		&Curriculum{},
		&Principal{},
		&Subject{},
		&ClassSubject{},
		&Material{},
		&Task{},
		&TaskSubmission{},
		&Assessment{},
		&AssessmentQuestion{},
		&AssessmentOption{},
		&AssessmentSubmission{},
		&AssessmentAnswer{},
	); err != nil {
		return err
	}

	setupActiveUniqueIndexes(db)
	return nil
}

// setupActiveUniqueIndexes creates partial/conditional unique indexes ensuring that
// unique constraints (email, NIP, NISN) apply only to active records (where deleted_at IS NULL).
// This preserves soft-deleted records in the database while allowing new accounts with
// the same email/NIP/NISN to be created with brand new Primary Keys.
func setupActiveUniqueIndexes(db *gorm.DB) {
	dial := db.Dialector.Name()

	switch dial {
	case "mysql":
		// Drop old single-column unique indexes if they exist
		_ = db.Exec("ALTER TABLE users DROP INDEX email").Error
		_ = db.Exec("ALTER TABLE users DROP INDEX uni_users_email").Error
		_ = db.Exec("ALTER TABLE curriculums DROP INDEX nip").Error
		_ = db.Exec("ALTER TABLE curriculums DROP INDEX uni_curriculums_nip").Error
		_ = db.Exec("ALTER TABLE teachers DROP INDEX nip").Error
		_ = db.Exec("ALTER TABLE teachers DROP INDEX uni_teachers_nip").Error
		_ = db.Exec("ALTER TABLE principals DROP INDEX nip").Error
		_ = db.Exec("ALTER TABLE principals DROP INDEX uni_principals_nip").Error
		_ = db.Exec("ALTER TABLE students DROP INDEX nisn").Error
		_ = db.Exec("ALTER TABLE students DROP INDEX uni_students_nisn").Error

		// In MySQL 8.0.13+, functional expression index:
		// When deleted_at IS NULL, the expression is 1 (uniqueness enforced among active rows).
		// When deleted_at IS NOT NULL, the expression is NULL (ignored by MySQL unique index, allowing multiple deleted entries).
		_ = db.Exec("CREATE UNIQUE INDEX idx_users_email_active ON users (email, ((CASE WHEN deleted_at IS NULL THEN 1 ELSE NULL END)))").Error
		_ = db.Exec("CREATE UNIQUE INDEX idx_curriculums_nip_active ON curriculums (nip, ((CASE WHEN deleted_at IS NULL THEN 1 ELSE NULL END)))").Error
		_ = db.Exec("CREATE UNIQUE INDEX idx_teachers_nip_active ON teachers (nip, ((CASE WHEN deleted_at IS NULL THEN 1 ELSE NULL END)))").Error
		_ = db.Exec("CREATE UNIQUE INDEX idx_principals_nip_active ON principals (nip, ((CASE WHEN deleted_at IS NULL THEN 1 ELSE NULL END)))").Error
		_ = db.Exec("CREATE UNIQUE INDEX idx_students_nisn_active ON students (nisn, ((CASE WHEN deleted_at IS NULL THEN 1 ELSE NULL END)))").Error

	case "postgres":
		_ = db.Exec("ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key").Error
		_ = db.Exec("DROP INDEX IF EXISTS idx_users_email").Error
		_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_active ON users (email) WHERE deleted_at IS NULL").Error

		_ = db.Exec("ALTER TABLE curriculums DROP CONSTRAINT IF EXISTS curriculums_nip_key").Error
		_ = db.Exec("DROP INDEX IF EXISTS idx_curriculums_nip").Error
		_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_curriculums_nip_active ON curriculums (nip) WHERE deleted_at IS NULL").Error

		_ = db.Exec("ALTER TABLE teachers DROP CONSTRAINT IF EXISTS teachers_nip_key").Error
		_ = db.Exec("DROP INDEX IF EXISTS idx_teachers_nip").Error
		_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_teachers_nip_active ON teachers (nip) WHERE deleted_at IS NULL").Error

		_ = db.Exec("ALTER TABLE principals DROP CONSTRAINT IF EXISTS principals_nip_key").Error
		_ = db.Exec("DROP INDEX IF EXISTS idx_principals_nip").Error
		_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_principals_nip_active ON principals (nip) WHERE deleted_at IS NULL").Error

		_ = db.Exec("ALTER TABLE students DROP CONSTRAINT IF EXISTS students_nisn_key").Error
		_ = db.Exec("DROP INDEX IF EXISTS idx_students_nisn").Error
		_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_students_nisn_active ON students (nisn) WHERE deleted_at IS NULL").Error

	case "sqlite":
		_ = db.Exec("DROP INDEX IF EXISTS uni_users_email").Error
		_ = db.Exec("DROP INDEX IF EXISTS idx_users_email").Error
		_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_active ON users (email) WHERE deleted_at IS NULL").Error

		_ = db.Exec("DROP INDEX IF EXISTS uni_curriculums_nip").Error
		_ = db.Exec("DROP INDEX IF EXISTS idx_curriculums_nip").Error
		_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_curriculums_nip_active ON curriculums (nip) WHERE deleted_at IS NULL").Error

		_ = db.Exec("DROP INDEX IF EXISTS uni_teachers_nip").Error
		_ = db.Exec("DROP INDEX IF EXISTS idx_teachers_nip").Error
		_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_teachers_nip_active ON teachers (nip) WHERE deleted_at IS NULL").Error

		_ = db.Exec("DROP INDEX IF EXISTS uni_principals_nip").Error
		_ = db.Exec("DROP INDEX IF EXISTS idx_principals_nip").Error
		_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_principals_nip_active ON principals (nip) WHERE deleted_at IS NULL").Error

		_ = db.Exec("DROP INDEX IF EXISTS uni_students_nisn").Error
		_ = db.Exec("DROP INDEX IF EXISTS idx_students_nisn").Error
		_ = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_students_nisn_active ON students (nisn) WHERE deleted_at IS NULL").Error

	default:
		log.Printf("Dialector %s: active partial unique indexes skipped", dial)
	}
}
