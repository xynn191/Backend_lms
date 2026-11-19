package models

import (
	"fmt"
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

func columnExists(db *gorm.DB, tableName, columnName string) bool {
	var count int64
	db.Raw("SELECT COUNT(1) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = ? AND column_name = ?", tableName, columnName).Scan(&count)
	return count > 0
}

func indexExists(db *gorm.DB, tableName, indexName string) bool {
	var count int64
	db.Raw("SELECT COUNT(1) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?", tableName, indexName).Scan(&count)
	return count > 0
}

// setupActiveUniqueIndexes creates partial/conditional unique indexes ensuring that
// unique constraints (email, NIP, NISN) apply only to active records (where deleted_at IS NULL).
// This preserves soft-deleted records in the database while allowing new accounts with
// the same email/NIP/NISN to be created with brand new Primary Keys.
func setupActiveUniqueIndexes(db *gorm.DB) {
	dial := db.Dialector.Name()

	switch dial {
	case "mysql":
		// 1. Migrasi kolom lama n_ip ke nip jika diperlukan
		tablesWithNIP := []string{"teachers", "curriculums", "principals"}
		for _, t := range tablesWithNIP {
			if columnExists(db, t, "n_ip") && !columnExists(db, t, "nip") {
				_ = db.Exec(fmt.Sprintf("ALTER TABLE %s CHANGE COLUMN n_ip nip VARCHAR(50)", t)).Error
			}
		}

		// 2. Drop old single-column unique indexes hanya jika memang ada di database
		dropOldIndex := func(table, idx string) {
			if indexExists(db, table, idx) {
				_ = db.Exec(fmt.Sprintf("ALTER TABLE %s DROP INDEX %s", table, idx)).Error
			}
		}

		dropOldIndex("users", "email")
		dropOldIndex("users", "uni_users_email")
		dropOldIndex("curriculums", "nip")
		dropOldIndex("curriculums", "uni_curriculums_nip")
		dropOldIndex("teachers", "nip")
		dropOldIndex("teachers", "uni_teachers_nip")
		dropOldIndex("principals", "nip")
		dropOldIndex("principals", "uni_principals_nip")
		dropOldIndex("students", "nisn")
		dropOldIndex("students", "uni_students_nisn")

		// 3. Buat partial active unique index hanya jika belum dibuat sebelumnya
		createIndexIfMissing := func(table, idx, query string) {
			if !indexExists(db, table, idx) {
				_ = db.Exec(query).Error
			}
		}

		createIndexIfMissing("users", "idx_users_email_active",
			"CREATE UNIQUE INDEX idx_users_email_active ON users (email, ((CASE WHEN deleted_at IS NULL THEN 1 ELSE NULL END)))")

		createIndexIfMissing("curriculums", "idx_curriculums_nip_active",
			"CREATE UNIQUE INDEX idx_curriculums_nip_active ON curriculums (nip, ((CASE WHEN deleted_at IS NULL THEN 1 ELSE NULL END)))")

		createIndexIfMissing("teachers", "idx_teachers_nip_active",
			"CREATE UNIQUE INDEX idx_teachers_nip_active ON teachers (nip, ((CASE WHEN deleted_at IS NULL THEN 1 ELSE NULL END)))")

		createIndexIfMissing("principals", "idx_principals_nip_active",
			"CREATE UNIQUE INDEX idx_principals_nip_active ON principals (nip, ((CASE WHEN deleted_at IS NULL THEN 1 ELSE NULL END)))")

		createIndexIfMissing("students", "idx_students_nisn_active",
			"CREATE UNIQUE INDEX idx_students_nisn_active ON students (nisn, ((CASE WHEN deleted_at IS NULL THEN 1 ELSE NULL END)))")

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
