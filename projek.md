# EduLMS Backend Project Tracker & Context

## Integrasi Brand Veyora
Frontend Veyora menggunakan backend ini melalui API `/api`; perubahan design system hanya memengaruhi presentasi frontend dan tidak mengubah kontrak endpoint.

## Konteks Proyek
EduLMS Backend adalah REST API untuk sistem manajemen pembelajaran sekolah berbasis Role-Based Access Control (Admin, Guru, Siswa, Kurikulum, Kepala Sekolah).
Dibangun menggunakan **Golang**, **Gin Framework**, **GORM** (ORM), dan **MySQL**.
Menggunakan pola **Clean Architecture**: Handler → Service → Repository.

## Tech Stack
- Language: Go 1.26
- Framework: Gin
- ORM: GORM (dengan Soft Deletes)
- Database: MySQL / MariaDB
- Auth: JWT (1 hari, via golang-jwt/v5)
- Config: .env + godotenv

## Konfigurasi Default (.env)
- APP_PORT=8080
- DB_NAME=edulms_db
- Default Admin: admin@gmail.com / Admin@123

## Struktur Direktori
- `internal/config/` — Load .env dan DSN builder
- `internal/models/` — GORM structs + AutoMigrate + Seeder
- `internal/repository/` — Query database
- `internal/service/` — Logika bisnis
- `internal/handler/` — HTTP handler (Gin)
- `internal/middleware/` — Auth JWT + Role check
- `internal/router/` — Pendaftaran semua routes

---

## Riwayat Pembuatan

### [SELESAI] Fondasi & Database
- [x] Inisialisasi Go Module (`github.com/xynn191/Backend_lms`)
- [x] Instalasi dependensi: Gin, GORM, MySQL driver, JWT, bcrypt, godotenv, uuid
- [x] `internal/config/config.go` — Config struct + Load() + DSN()
- [x] `internal/models/auth.go` — Model Role dan User (UUID primary key) + konstanta role name
- [x] `internal/models/profiles.go` — Admin, Teacher, Student, Curriculum, Principal
- [x] `internal/models/master.go` — Major, Class, Subject, ClassSubject
- [x] `internal/models/learning.go` — Material, Task, TaskSubmission, Assessment, AssessmentSubmission
- [x] `internal/models/migrate.go` — AutoMigrate() semua tabel
- [x] `internal/models/seed.go` — Seed roles (5 role) + user admin default

### [SELESAI] Auth Layer
- [x] `internal/middleware/auth_middleware.go` — AuthRequired() validasi JWT, RoleRequired() cek role
- [x] `internal/repository/auth_repository.go` — FindUserByEmail, FindUserByID
- [x] `internal/service/auth_service.go` — Login (bcrypt + JWT generate), GetProfile()
- [x] `internal/handler/auth_handler.go` — POST /api/auth/login, GET /api/auth/me

### [SELESAI] Admin Endpoints
- [x] `internal/repository/admin_repository.go` — GetDashboardStats()
- [x] `internal/service/admin_service.go` — GetDashboardStats()
- [x] `internal/handler/admin_handler.go` — GET /api/admin/dashboard

### [SELESAI] Master Data Endpoints (Admin Only)
- [x] `internal/repository/major_repository.go` — CRUD Jurusan
- [x] `internal/service/major_service.go` — CRUD Jurusan
- [x] `internal/handler/major_handler.go` — GET/POST/PUT/DELETE /api/admin/majors

- [x] `internal/repository/class_repository.go` — CRUD Kelas (dengan Preload)
- [x] `internal/service/class_service.go` — CRUD Kelas
- [x] `internal/handler/class_handler.go` — GET/POST/PUT/DELETE /api/admin/classes

- [x] `internal/repository/subject_repository.go` — CRUD Mata Pelajaran
- [x] `internal/service/subject_service.go` — CRUD Mata Pelajaran
- [x] `internal/handler/subject_handler.go` — GET/POST/PUT/DELETE /api/admin/subjects

### [SELESAI] User Management Endpoints (Admin Only)
- [x] `internal/repository/user_repository.go` — CRUD Student, Teacher, Curriculum, Principal + FindRoleByName
- [x] `internal/service/user_service.go` — CRUD semua 4 tipe user (menggunakan uuid, atomic transaction)
- [x] `internal/handler/user_handler.go` — CRUD /api/admin/students, /teachers, /curriculums, /principals

### [SELESAI] Student API Endpoints (Siswa Only)
- [x] `internal/repository/student_api_repository.go` — GetStudentProfile, GetMaterials, GetTasks, GetDashboardStats (query berdasarkan classID)
- [x] `internal/service/student_api_service.go` — GetDashboardStats, GetMaterials, GetTasks
- [x] `internal/handler/student_api_handler.go` — GET /api/student/dashboard, /materials, /tasks

### [SELESAI] Router
- [x] `internal/router/router.go` — Wire semua dependency + register semua 28 endpoint

---

## Semua Endpoint Aktif
| Method | Endpoint | Role |
|---|---|---|
| POST | /api/auth/login | Public |
| GET | /api/auth/me | Any |
| GET | /api/admin/dashboard | Admin |
| GET/POST/PUT/DELETE | /api/admin/majors | Admin |
| GET/POST/PUT/DELETE | /api/admin/classes | Admin |
| GET/POST/PUT/DELETE | /api/admin/subjects | Admin |
| GET/POST/PUT/DELETE | /api/admin/students | Admin |
| GET/POST/PUT/DELETE | /api/admin/teachers | Admin |
| GET/POST/PUT/DELETE | /api/admin/curriculums | Admin |
| GET/POST/PUT/DELETE | /api/admin/principals | Admin |
| GET | /api/student/dashboard | Siswa |
| GET | /api/student/materials | Siswa |
| GET | /api/student/tasks | Siswa |

---

## Todo Berikutnya
- [ ] Guru API: endpoint untuk guru (lihat kelas, buat materi, buat tugas, beri nilai)
- [ ] Kurikulum & Kepsek API: endpoint monitoring dan laporan
- [ ] File upload untuk lampiran materi dan pengumpulan tugas
