# EduLMS Backend Project Tracker & Context

## Konteks Proyek
EduLMS Backend adalah REST API yang melayani sistem manajemen pembelajaran sekolah berbasis Role-Based Access Control (Admin, Guru, Siswa, Kurikulum, Kepala Sekolah).
Dibangun menggunakan **Golang**, **Gin Framework**, **GORM** (ORM), dan **MySQL** dengan mengadaptasi pola **Clean Architecture** (Handler -> Service -> Repository).

## Riwayat & Status (Backend)
- [x] Inisialisasi Go Module (`github.com/xynn191/Backend_lms`)
- [x] Instalasi dependensi (Gin, GORM, MySQL Driver, JWT, godotenv, bcrypt)
- [x] Setup struktur direktori internal (config, handler, middleware, models, repository, router, service)
- [x] **Models & Database**: Pembuatan schema `Role`, `User`, Profil Role (`Admin`, `Teacher`, `Student`, dll), Master Data (`Major`, `Class`, `Subject`), dan Modul Pembelajaran (`Material`, `Task`, `Assessment`).
- [x] Auto-Migrate script & Seeder awal (Roles + Akun Admin Default)
- [x] **Auth Layer**: Middleware verifikasi JWT token dan verifikasi Role.
- [x] **Auth Endpoints**: Login (`POST /api/auth/login`) dan Profile (`GET /api/auth/me`).
- [x] **Admin Endpoints**: CRUD lengkap untuk Master Data (Jurusan, Kelas, Mapel) dan Manajemen User.
- [x] **Student Endpoints**: API untuk Dashboard Siswa (`GET /api/student/dashboard`), Materi (`GET /api/student/materials`), dan Tugas (`GET /api/student/tasks`) berdasarkan relasi kelas siswa.

## Todo Berikutnya
- Pembuatan API untuk Role Guru (Membuat materi, tugas, memberi nilai, manajemen kelas)
- Pembuatan API untuk Role Kurikulum & Kepala Sekolah (Laporan & Monitoring)
- Implementasi fitur upload file untuk Attachment Materi / Jawaban Tugas
