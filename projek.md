# Veyora LMS Backend Project Tracker & Context

## [SELESAI] Penjabaran Nama Guru Pengampu Mapel (`SubjectMonitoringItem.Teachers`)
- [x] **Perbaikan Tag Mapping GORM (`internal/repository/monitoring_repository.go`)**:
  - Menambahkan tag mapping eksplisit `gorm:"column:id"` dan `gorm:"column:name"` pada struct `SubjectTeacherInfo` agar GORM sukses memetakan hasil query SQL `DISTINCT teachers.id, teachers.name, teachers.nip` ke struct Go tanpa menghasilkan string kosong.
- [x] **Skema DTO Guru Pengampu (`SubjectTeacherInfo`) & Integrasi Kepala Sekolah**:
  - Memperbarui `SubjectMonitoringItem` dengan field `Teachers []SubjectTeacherInfo` dan `TeacherCount int64`.
  - Memperbarui `PrincipalMonitoringData` dengan menyertakan `SubjectStats []SubjectMonitoringItem` agar Kepala Sekolah dan Kurikulum mendapatkan data yang identik.
- [x] **Query Relasi Penugasan Guru per Mata Pelajaran (`monitoringRepository.getSubjectStats`)**:
  - Mengimplementasikan query `DISTINCT` yang mengambil seluruh guru yang ditugaskan ke setiap mapel melalui tabel relasi `class_subjects`.
  - Mendukung kondisi multi-guru (jika suatu mapel diajar oleh lebih dari 1 guru di kelas-kelas berbeda, semua guru dikembalikan dalam array `teachers`).
- [x] **Unit Testing & Kompilasi Binary (`internal/test/monitoring_test.go`, `server`)**:
  - Menambahkan unit test `TestSubjectTeachersInMonitoring` (PASS) yang menguji pemetaan nama guru pengampu mapel.
  - Mengompilasi binary `server` terbaru dengan `go build -o server .`.

---

## [SELESAI] Hidden Role "Superadmin" (Manajemen Pengguna Murni & Proteksi RBAC Eksklusif)
- [x] **Konstanta & Seeding Role Superadmin (`internal/models/auth.go`, `internal/models/seed.go`)**:
  - Menambahkan konstanta `RoleSuperadmin = "superadmin"` dan mendaftarkannya ke slice `SupportedRoles`.
  - Mendaftarkan Role ID 6 untuk `RoleSuperadmin` pada `seedRoles(db)`.
  - Mengimplementasikan `seedSuperadminUser(db)` dengan kredensial default `superadmin@veyora.id` / `superadmin123`.
- [x] **Repository CRUD Admin & Agregasi 5 Pengguna (`internal/repository/user_repository.go`)**:
  - Menambahkan method `GetAllAdmins()`, `FindAdminByID()`, `CreateUserWithAdmin()`, `UpdateAdmin()`, dan `DeleteAdmin()`.
  - Mengimplementasikan `GetSuperadminDashboardStats()` yang menghitung agregasi riil 5 entitas: Administrator, Siswa, Guru, Tim Kurikulum, dan Kepala Sekolah beserta total pengguna aktif.
- [x] **Service & DTO Admin (`internal/service/user_service.go`)**:
  - Menambahkan DTO `CreateAdminDTO` dan `UpdateAdminDTO`.
  - Menambahkan method service `GetAllAdmins`, `CreateAdmin`, `UpdateAdmin`, `DeleteAdmin`, dan `GetSuperadminDashboardStats`.
- [x] **Handler & Rute RBAC Superadmin (`internal/handler/user_handler.go`, `internal/router/router.go`)**:
  - Mendaftarkan grup rute eksklusif `/api/superadmin` terproteksi `middleware.RoleRequired(models.RoleSuperadmin)`:
    - `GET /api/superadmin/dashboard` (Statistik 5 pengguna)
    - `GET /api/superadmin/admins`, `POST /api/superadmin/admins`, `PUT /api/superadmin/admins/:id`, `DELETE /api/superadmin/admins/:id`
  - Membuka akses manajemen user sekolah (`/api/admin/students`, `/api/admin/teachers`, `/api/admin/curriculums`, `/api/admin/principals`, serta read-only `/api/admin/classes`) untuk `RoleAdmin` dan `RoleSuperadmin`.
  - Menolak akses `RoleSuperadmin` terhadap seluruh master data kurikulum (`/api/admin/majors`, `/api/admin/subjects`, write kelas, penugasan) dengan `403 Forbidden`.
  - Menolak akses `RoleAdmin` terhadap `/api/superadmin/*` dengan `403 Forbidden`.
- [x] **Unit Testing RBAC Superadmin (`internal/test/superadmin_rbac_test.go`)**:
  - Memvalidasi 6 skenario pengujian komprehensif (akses dashboard, pencegahan akses admin biasa ke dashboard/CRUD admin, pembuatan admin baru, pencegahan akses superadmin ke master data jurusan, dan izin akses ke daftar siswa) dengan status PASS 100%.
- [x] Verifikasi build backend: `go test ./...` sukses dengan Exit Code 0.

---

## [SELESAI] Dukungan Multi-File & Screenshot Lampiran Tugas Guru (`Task.AttachmentURL text`)
- [x] **Peningkatan Skema Database `tasks.attachment_url` (`internal/models/learning.go`)**:
  - Kolom `AttachmentURL` pada struct `Task` ditingkatkan dari `varchar(255)` menjadi `text` untuk mendukung serialisasi JSON array multi-file lampiran dokumen dan gambar tangkapan layar (screenshot).
- [x] Verifikasi build backend: `go build -o server main.go` dan `go test ./...` sukses dengan Exit Code 0.

---

## [SELESAI] Isolasi Route RBAC & Integrasi Aktivitas Guru untuk Kurikulum & Kepala Sekolah
- [x] **Aktivitas Guru pada Endpoint Kurikulum (`TeacherActivityItem`)**:
  - `internal/repository/monitoring_repository.go`: Menambahkan field `TeacherActivity []TeacherActivityItem` pada struct `CurriculumMonitoringData`.
  - Mengimplementasikan penghitungan aktivitas riil per guru (`subject_count`, `material_count`, `task_count`) pada method `GetCurriculumMonitoring()`.
- [x] **Eksposur Route Read-Only Master Data untuk Kurikulum & Kepala Sekolah**:
  - `internal/router/router.go`: Mendaftarkan endpoint baca kelas & mapel tanpa melanggar proteksi admin:
    - `curriculum.GET("/classes", classH.GetAll)`
    - `curriculum.GET("/subjects", subjectH.GetAll)`
    - `principal.GET("/classes", classH.GetAll)`
    - `principal.GET("/subjects", subjectH.GetAll)`
  - Menghilangkan potensi error `403 Forbidden` saat akun non-admin (Kurikulum / Kepala Sekolah) mengakses data referensi kelas & mapel.
- [x] Verifikasi build backend: `go build -o server main.go` dan `go test ./...` sukses dengan Exit Code 0.

---

## [SELESAI] API Nilai Siswa Sekolah untuk Kurikulum & Kepala Sekolah
- [x] **Repository Nilai Tingkat Sekolah (`SchoolGradeItem`)**:
  - `internal/repository/monitoring_repository.go`: Menambahkan struct `SchoolGradeItem` (student_id, student_name, student_nisn, class_id, class_name, subject_id, subject_name, task_average, assessment_average, final_score).
  - Mengimplementasikan `GetSchoolGrades(classID uint, subjectID uint) ([]SchoolGradeItem, error)` yang menghitung agregasi tugas & ujian riil siswa per kelas dan mata pelajaran.
- [x] **Service & Handler Integrasi**:
  - `internal/service/monitoring_service.go`: Meneruskan pemanggilan `GetSchoolGrades`.
  - `internal/handler/monitoring_handler.go`: Menambahkan handler `SchoolGrades` dengan dukungan filter query string `class_id` dan `subject_id`.
  - `internal/router/router.go`: Mendaftarkan endpoint `GET /api/curriculum/grades` dan `GET /api/principal/grades` terlindungi RBAC.
- [x] Verifikasi build backend: `go build -o server main.go` dan `go test ./...` sukses dengan Exit Code 0.

---

## [SELESAI] Pembersihan Startup Migrations & Eliminasi Warning SQL Backend
- [x] **Standardisasi Kolom NIP (`column:nip`)**:
  - `internal/models/profiles.go`: Menambahkan tag eksplisit `gorm:"type:varchar(50);column:nip"` pada struct `Teacher`, `Curriculum`, dan `Principal`. Mencegah GORM membuat kolom default `n_ip` (penyebab Error 1072).
- [x] **Migrasi Skema & Indeks Idempoten (`internal/models/migrate.go`)**:
  - Menambahkan helper `columnExists(db, table, col)` dan `indexExists(db, table, idx)` berbasis `information_schema`.
  - Migrasi otomatis jika kolom lama `n_ip` terdeteksi di tabel MySQL (`ALTER TABLE ... CHANGE COLUMN n_ip nip VARCHAR(50)`).
  - Melakukan `DROP INDEX` hanya jika nama index lama benar-benar ditemukan di database (mengeliminasi Error 1091).
  - Melakukan `CREATE UNIQUE INDEX` hanya jika index belum dibuat sebelumnya (mengeliminasi Error 1061).
- [x] **Pembersihan Log Framework Gin (`internal/router/router.go`)**:
  - Menambahkan `_ = r.SetTrustedProxies(nil)` untuk menghilangkan peringatan `You trusted all proxies, this is NOT safe`.
- [x] Verifikasi build backend: `go build -o server main.go` dan `go test ./...` sukses dengan Exit Code 0.

---

## [SELESAI] Fitur Edit Email pada CRUD User Admin (Siswa, Guru, Kurikulum, Kepsek)
- [x] **Dukungan Update Email pada DTO User**:
  - `internal/service/user_service.go`: Menambahkan field `Email string json:"email" binding:"omitempty,email"` pada `UpdateStudentDTO`, `UpdateTeacherDTO`, dan `UpdateStaffDTO`.
- [x] **Pemeriksaan Keunikan Email Aktif (`IsEmailAvailableForUser`)**:
  - `internal/repository/user_repository.go`: Menambahkan method `IsEmailAvailableForUser(email string, excludeUserID string) (bool, error)` untuk memastikan email baru tidak bertabrakan dengan akun aktif lain, sekaligus memperbolehkan akun yang sedang diedit tetap mempertahankan emailnya sendiri tanpa false positive error.
  - `internal/repository/user_repository.go`: Menambahkan method `UpdateUserEmail(userID string, email string) error` untuk memperbarui kolom `users.email` berdasarkan `user_id`.
- [x] **Logika Pembaruan Email di Service**:
  - `internal/service/user_service.go`: Pada `UpdateStudent`, `UpdateTeacher`, `UpdateCurriculum`, dan `UpdatePrincipal`, jika admin mengubah email pengguna, sistem memvalidasi ketersediaan email baru, menolak jika duplikat (`"email ... sudah digunakan oleh akun aktif lain"`), dan memperbarui email pada tabel `users`.
  - Mengembalikan record data teranyar dengan relasi `User` yang telah diperbarui.
- [x] Verifikasi build backend: `go build -o server main.go` dan `go test ./...` sukses dengan Exit Code 0.

---

## [SELESAI] Perbaikan Validasi User & Penanganan Error CRUD (Import Excel)
- [x] **Pemeriksaan Keberadaan Kelas (`ClassExists`)**:
  - `internal/repository/user_repository.go`: Menambahkan method `ClassExists(id uint) (bool, error)` untuk memvalidasi `class_id` sebelum data siswa disimpan.
  - `internal/service/user_service.go`: Mencegah error raw MySQL constraint (Error 1452) pada `CreateStudent` dan `UpdateStudent`. Jika `class_id` tidak ada di database, sistem mengembalikan pesan deskriptif: *"kelas dengan ID X tidak ditemukan di sistem"*.
- [x] **Standarisasi Status Code HTTP pada User Handler**:
  - `internal/handler/user_handler.go`: Mengubah status error pada `CreateStudent`, `UpdateStudent`, `CreateTeacher`, `UpdateTeacher`, `CreateCurriculum`, `UpdateCurriculum`, `CreatePrincipal`, `UpdatePrincipal` dari `500 Internal Server Error` menjadi `400 Bad Request` dengan payload JSON `{"error": "..."}`.
  - Dampak: Error validasi input, email duplikat, NISN/NIP duplikat, dan foreign key kelas kini mengembalikan status code dan pesan yang jelas bagi client/frontend saat impor massal Excel.
- [x] Verifikasi build backend: `go build -o server main.go` sukses dengan Exit Code 0.

---

## [SELESAI] Multi-File Attachment Materi & Penjadwalan Assessment Fleksibel
- [x] **Skema Database `materials.attachment_url` (`internal/models/learning.go`)**:
  - Tipe kolom `AttachmentURL` ditingkatkan dari `varchar(255)` menjadi `text` untuk mendukung serialisasi multi-lampiran (JSON string array) beserta metadata nama file dan tipe file.
- [x] **Endpoint Cepat Jadwal Assessment (`PATCH /api/teacher/assessments/:id/schedule`)**:
  - `internal/handler/teacher_api_handler.go`: Menambahkan handler `UpdateAssessmentSchedule` untuk update `start_time`, `end_time`, dan `duration_minutes`.
  - `internal/service/teacher_api_service.go`: Menambahkan validasi `EndTime > StartTime` tanpa memicu re-writing daftar soal.
  - `internal/repository/teacher_api_repository.go`: Mengimplementasikan `UpdateAssessmentSchedule` menggunakan GORM updates parsial.
  - `internal/router/router.go`: Mendaftarkan route `PATCH /assessments/:id/schedule`.
- [x] **Proteksi Akses Soal Ujian Siswa Sebelum Jadwal Dimulai**:
  - `internal/repository/student_api_repository.go`: Pada `GetAssessmentByID`, akses soal otomatis ditolak jika `time.Now().Before(a.StartTime)` dan siswa belum memiliki submission (mencegah kebocoran soal sebelum jam ujian dibuka).
- [x] Verifikasi build backend: `go build -o server main.go` sukses dengan Exit Code 0.

---

## [SELESAI] Database Migration & Schema Fixes
- [x] **Soft Delete & Partial Unique Index Architecture (Re-create Account with Fresh PK)**:
  - `internal/models/auth.go` & `internal/models/profiles.go`: Melepaskan tag rigid `;unique` yang sebelumnya memblokir soft-deleted records pada kolom `users.email`, `teachers.nip`, `students.nisn`, `curriculums.nip`, dan `principals.nip`.
  - `internal/models/migrate.go`: Mengimplementasikan migrasi otomatis `setupActiveUniqueIndexes(db)` yang menghapus index lama dan membuat **Partial Unique Index** (hanya berlaku jika `deleted_at IS NULL`):
    - MySQL 8.0.13+: Menggunakan functional expression index `((CASE WHEN deleted_at IS NULL THEN 1 ELSE NULL END))`.
    - PostgreSQL / SQLite: Menggunakan conditional unique index `WHERE deleted_at IS NULL`.
  - `internal/repository/user_repository.go`: Menambahkan helper verifikasi ketersediaan aktif: `IsEmailAvailable`, `IsNIPAvailable`, dan `IsNISNAvailable`.
  - `internal/service/user_service.go`: Integrasi validasi pra-pembuatan akun untuk Siswa, Guru, Tim Kurikulum, dan Kepala Sekolah.
  - **Hasil**: Data lama yang dihapus tetap tersimpan rapi di database (soft delete/arsip riwayat). Saat admin mendaftarkan ulang akun dengan email atau NIP yang sama (misal `kurikulum@gmail.com`), sistem sukses membuat baris baru dengan **Primary Key (PK) baru** tanpa terjadi konflik constraint database.
- [x] **Fix Foreign Key Constraint `classes.academic_year_id` (Error 1452)**:
  - `internal/models/master.go`: Field `AcademicYearID` pada model `Class` diubah menjadi pointer nullable `*uint` dengan `gorm:"foreignKey:AcademicYearID;constraint:-;"`.
  - `internal/models/migrate.go`: Alur migrasi disempurnakan:
    1. Migrasi tabel `academic_years` terlebih dahulu.
    2. Seed tahun akademik aktif default ("2026/2027 Ganjil") jika belum ada data.
    3. Eksekusi `UPDATE classes SET academic_year_id = ? WHERE academic_year_id IS NULL OR academic_year_id = 0` sebelum constraint diterapkan untuk mencegah kegagalan FK pada database yang sudah berisi data.
- [x] Verifikasi build backend: `go build ./...` berjalan sukses tanpa error.

---

## Konteks Proyek
Veyora LMS Backend adalah REST API untuk sistem manajemen pembelajaran sekolah berbasis Role-Based Access Control (Admin, Guru, Siswa, Kurikulum, Kepala Sekolah).
Dibangun menggunakan **Golang**, **Gin Framework**, **GORM** (ORM), dan **MySQL**.
Menggunakan arsitektur berlapis (**Clean Architecture**): Handler → Service → Repository.

## Tech Stack
- Language: Go 1.26
- Framework: Gin Web Framework
- ORM: GORM (dengan Soft Deletes & Migrations)
- Database: MySQL / MariaDB
- Auth: JWT (1 hari, via golang-jwt/v5) + bcrypt
- Config: .env + godotenv

## Konfigurasi Default (.env)
- APP_PORT=8080
- DB_NAME=lms
- Default Admin: admin@gmail.com / Admin@123

## Struktur Direktori
- `internal/config/` — Load .env dan DSN builder
- `internal/models/` — GORM structs, AutoMigrate, Seeder, dan Schema relationships
- `internal/repository/` — Query database (Admin, Auth, User, Master, Teacher, Student, Monitoring)
- `internal/service/` — Logika bisnis (Auth, Admin, User, Master, Teacher, Student, Monitoring)
- `internal/handler/` — HTTP handler Gin (Auth, Admin, User, Master, Teacher, Student, Monitoring, Upload)
- `internal/middleware/` — Auth JWT + Variadic Role checking (`RoleRequired(roles ...string)`)
- `internal/router/` — Pendaftaran semua routes dan CORS

---

## Riwayat Pembuatan & Modul Selesai

### [SELESAI] Fondasi & Database
- [x] Inisialisasi Go Module (`github.com/xynn191/Backend_lms`)
- [x] Instalasi dependensi: Gin, GORM, MySQL driver, JWT, bcrypt, godotenv, uuid
- [x] `internal/config/config.go` — Config struct + Load() + DSN()
- [x] `internal/models/auth.go` — Model Role dan User (UUID primary key) + konstanta role name
- [x] `internal/models/profiles.go` — Admin, Teacher, Student, Curriculum, Principal
- [x] `internal/models/master.go` — AcademicYear, Major, Class, Subject, ClassSubject
- [x] `internal/models/learning.go` — Material, Task, TaskSubmission, Assessment, AssessmentQuestion, AssessmentOption, AssessmentSubmission, AssessmentAnswer
- [x] `internal/models/migrate.go` — AutoMigrate() semua tabel dengan dependency ordering aman
- [x] `internal/models/seed.go` — Seed roles (5 role) + user admin default + tahun akademik default

### [SELESAI] Auth Layer
- [x] `internal/middleware/auth_middleware.go` — AuthRequired() validasi JWT, RoleRequired(roles ...string)
- [x] `internal/repository/auth_repository.go` — FindUserByEmail, FindUserByID
- [x] `internal/service/auth_service.go` — Login (bcrypt + JWT generate), GetProfile()
- [x] `internal/handler/auth_handler.go` — POST /api/auth/login, GET /api/auth/me

### [SELESAI] Master Data & User Management (Admin Only)
- [x] `internal/repository/major_repository.go` & `service/` & `handler/` — CRUD Jurusan (`/api/admin/majors`)
- [x] `internal/repository/class_repository.go` & `service/` & `handler/` — CRUD Kelas (`/api/admin/classes`)
- [x] `internal/repository/subject_repository.go` & `service/` & `handler/` — CRUD Mata Pelajaran (`/api/admin/subjects`)
- [x] `internal/repository/class_subject_repository.go` & `service/` & `handler/` — CRUD Penugasan Guru (`/api/admin/class-subjects`)
- [x] `internal/repository/user_repository.go` & `service/` & `handler/` — CRUD Siswa, Guru, Tim Kurikulum, Kepala Sekolah

### [SELESAI] Pembelajaran & Penilaian Guru (Teacher API)
- [x] `internal/repository/teacher_api_repository.go`, `service/`, `handler/`:
  - GET `/api/teacher/dashboard`: Ringkasan kelas diampu, total materi, tugas, dan submission pending
  - CRUD `/api/teacher/materials` + PATCH `/materials/:id/publish`
  - CRUD `/api/teacher/tasks` + PATCH `/tasks/:id/publish`
  - GET `/api/teacher/tasks/:id/submissions`: Daftar pengumpulan tugas siswa
  - PUT `/api/teacher/submissions/:submission_id/grade`: Penilaian nilai manual + feedback guru
  - CRUD `/api/teacher/assessments` + PATCH `/assessments/:id/publish` + GET `/assessments/:id/results`
  - GET `/api/teacher/grades` & `/grades/detail`: Rekap nilai terpadu per kelas dan per siswa

### [SELESAI] Portal Siswa (Student API)
- [x] `internal/repository/student_api_repository.go`, `service/`, `handler/`:
  - GET `/api/student/dashboard`: Ringkasan mapel, tugas pending, dan status kuis
  - GET `/api/student/materials` & `/materials/:id`: Akses modul bahan ajar
  - GET `/api/student/tasks` & `/tasks/:id`: Daftar dan detail instruksi tugas
  - POST `/api/student/tasks/:id/submit`: Pengumpulan tugas (URL file lampiran dan tautan link)
  - GET `/api/student/assessments` & `/assessments/:id`: Akses kuis/ujian dengan opsi soal teracak dan flag kunci jawaban disembunyikan
  - POST `/api/student/assessments/:id/submit`: Pengerjaan kuis dengan evaluasi dan perhitungan skor otomatis (auto-scoring 0-100)

### [SELESAI] File Upload Engine (Local Storage)
- [x] `internal/handler/upload_handler.go`:
  - POST `/api/upload`: Menerima multipart/form-data
  - Batasan ukuran file: Maksimum 25 MB
  - Whitelist ekstensi aman: PDF, DOC, DOCX, XLS, XLSX, PPT, PPTX, JPG, JPEG, PNG, WEBP, GIF, MP3, WAV, MP4, MKV
  - Simpan lokal di direktori `./uploads` dengan UUID penamaan unik
  - Static endpoint: `r.Static("/uploads", "./uploads")`

### [SELESAI] Monitoring Kurikulum & Kepala Sekolah
- [x] `internal/repository/monitoring_repository.go`, `service/`, `handler/`:
  - GET `/api/curriculum/monitoring`: Keterisian silabus pembelajaran, persentase pengumpulan tugas siswa, rata-rata nilai sekolah, data progres per kelas, dan status materi per mata pelajaran
  - GET `/api/principal/monitoring`: Partisipasi akademik sekolah, efektivitas pendidik, rerata nilai mutu sekolah, distribusi rombel per jurusan, dan rekap keaktifan guru

### [SELESAI] Router & Security
- [x] `internal/router/router.go`: Wire seluruh repositori, service, dan handler dengan dependency injection
- [x] Middleware JWT Bearer Authentication & Role-Based Authorization
- [x] Konfigurasi CORS responsif untuk `http://localhost:3000`

---

## Daftar Lengkap Endpoint Aktif (30 Endpoint)

| Method | Endpoint | Akses Role | Deskripsi |
|---|---|---|---|
| POST | `/api/auth/login` | Public | Autentikasi email & password, mengembalikan JWT token |
| GET | `/api/auth/me` | Authenticated | Profil pengguna aktif & role saat ini |
| POST | `/api/upload` | Authenticated | Upload berkas materi/tugas (maks 25 MB, disimpan di `./uploads`) |
| GET | `/api/admin/dashboard` | Admin | Statistik total siswa, guru, kelas, mata pelajaran |
| GET/POST/PUT/DELETE | `/api/admin/majors` | Admin | Master data program keahlian / jurusan |
| GET/POST/PUT/DELETE | `/api/admin/classes` | Admin | Master data kelas, tingkat, jurusan, dan wali kelas |
| GET/POST/PUT/DELETE | `/api/admin/subjects` | Admin | Master data mata pelajaran sekolah |
| GET/POST/PUT/DELETE | `/api/admin/class-subjects` | Admin | Penugasan guru pengampu per kelas dan mapel |
| GET/POST/PUT/DELETE | `/api/admin/students` | Admin | Manajemen akun dan data profil siswa |
| GET/POST/PUT/DELETE | `/api/admin/teachers` | Admin | Manajemen akun dan data profil guru |
| GET/POST/PUT/DELETE | `/api/admin/curriculums` | Admin | Manajemen akun tim pengembang kurikulum |
| GET/POST/PUT/DELETE | `/api/admin/principals` | Admin | Manajemen akun kepala sekolah |
| GET | `/api/teacher/dashboard` | Guru | Ringkasan kegiatan mengajar guru |
| GET/POST/PUT/DELETE | `/api/teacher/materials` | Guru | Kelola modul pembelajaran guru |
| PATCH | `/api/teacher/materials/:id/publish` | Guru | Terbitkan / tarik materi ajar |
| GET/POST/PUT/DELETE | `/api/teacher/tasks` | Guru | Kelola penugasan siswa |
| PATCH | `/api/teacher/tasks/:id/publish` | Guru | Terbitkan / sembunyikan tugas |
| GET | `/api/teacher/tasks/:id/submissions` | Guru | Daftar submission tugas siswa |
| PUT | `/api/teacher/submissions/:submission_id/grade` | Guru | Input nilai angka & catatan evaluasi guru |
| GET/POST/PUT/DELETE | `/api/teacher/assessments` | Guru | Kelola ujian online pilihan ganda |
| GET | `/api/student/dashboard` | Siswa | Ringkasan materi, tugas, dan ujian siswa |
| GET | `/api/student/materials` | Siswa | Daftar materi mata pelajaran kelas siswa |
| GET | `/api/student/materials/:id` | Siswa | Detail isi modul materi dan unduh lampiran |
| GET | `/api/student/tasks` | Siswa | Daftar tugas yang harus dikerjakan |
| GET | `/api/student/tasks/:id` | Siswa | Detail instruksi tugas & status pengumpulan |
| POST | `/api/student/tasks/:id/submit` | Siswa | Kirim berkas tugas / tautan link |
| GET | `/api/student/assessments` | Siswa | Daftar kuis online terjadwal |
| GET | `/api/student/assessments/:id` | Siswa | Lembar butir soal ujian online |
| POST | `/api/student/assessments/:id/submit` | Siswa | Kumpulkan jawaban ujian dengan auto-scoring instan |
| GET | `/api/curriculum/monitoring` | Kurikulum, Admin | Metrik ketercapaian silabus & mutu akademik |
| GET | `/api/principal/monitoring` | Kepala Sekolah, Admin | Ringkasan eksekutif performa dan mutu sekolah |

---

## Verifikasi Kualitas
- Build Check: `go build ./...` -> **0 Error (Exit Code 0)**
