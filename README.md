# Travel Itinerary — Full-Stack Web Application

**Travel Itinerary** adalah aplikasi web **full-stack** untuk mengelola rencana perjalanan wisata secara terstruktur mulai dari registrasi/login pengguna, menyusun *itinerary*, menambahkan destinasi wisata (*spots*), hingga mencatat informasi penting seperti jam operasional, tiket, dan tips kunjungan pada setiap destinasi (*spot notes*).

Proyek ini menggabungkan backend RESTful API berbasis **Go**, **Fiber v3**, **GORM**, dan **MySQL** (`travel_itinerary_db`) dengan antarmuka web **React + Vite** di folder `frontend/`.

---

## Project Overview

Aplikasi ini dirancang agar pengguna dapat merencanakan perjalanan dengan rapi dan terhubung antar-entitas:

- **Backend REST API:** Menangani logika bisnis, autentikasi, validasi input, relasi GORM, dan operasi database MySQL.
- **Frontend Web App:** Menyediakan antarmuka interaktif (layout 2 kolom dengan tema warna hangat/ceria) untuk mengelola autentikasi/token, *itineraries*, *spots*, dan *spot notes*.

### Tujuan Project

Project ini dikembangkan dengan beberapa tujuan utama:

- Mempermudah pengguna dalam mengelola rencana perjalanan.
- Menyimpan daftar destinasi wisata dalam setiap *itinerary*.
- Menyediakan fitur catatan untuk setiap destinasi wisata.
- Mengelola data pengguna dan relasi antar data secara terstruktur.
- Menyediakan REST API yang terintegrasi dengan frontend.
- Menerapkan pengembangan aplikasi full-stack menggunakan Go, Fiber v3, GORM, MySQL, dan React + Vite.

### Fitur Utama

| Fitur | Deskripsi |
| :--- | :--- |
| **User Management** | Mengelola data pengguna aplikasi |
| **Authentication** | Menangani proses registrasi, login, dan token autentikasi pengguna |
| **Itinerary Management** | Membuat, melihat, mengubah, dan menghapus rencana perjalanan |
| **Spot Management** | Mengelola daftar destinasi wisata dalam *itinerary* |
| **Spot Notes** | Menambahkan dan mengelola catatan (termasuk jam operasional & tips) pada setiap destinasi wisata |

---

## Pembagian Modul Tim

| Modul | Tabel / Area | Penanggung Jawab | Fokus Utama |
| :--- | :--- | :--- | :--- |
| **Auth & JWT** | `users` | Fikri | Autentikasi, middleware/JWT (`auth.go`), & manajemen pengguna |
| **Itineraries** | `itineraries` | Nabil | CRUD rencana perjalanan |
| **Spots** | `spots` | Aurellia | CRUD destinasi wisata dalam *itinerary* |
| **Spot Notes** | `spot_notes` | Reva ar | CRUD catatan spot, `Preload("Spot")`, validasi input, & error handling (`400`, `404`, `422`) |
| **Frontend & Integrasi** | UI React + Vite (`frontend/`) | Frontend | Integrasi API, manajemen token JWT, CRUD berjenjang, styling, & persiapan deployment |

---

## Struktur Direktori Proyek

```text
travel-itinerary-fiber-d1/
├── README.md
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go                 # Entry point, CORS & middleware, AutoMigrate, dan routing
│   └── internal/
│       ├── config/
│       │   └── db.go                   # Konfigurasi environment & koneksi MySQL (GORM)
│       ├── handler/                    # HTTP handler (auth, itineraries, spots, spot_note_handler.go)
│       ├── middleware/
│       │   └── auth.go                 # Middleware autentikasi
│       ├── model/
│       │   └── models.go               # Definisi model GORM (User, Itinerary, Spot, SpotNote)
│       ├── repository/                 # Query & akses data database
│       └── service/                    # Logika bisnis tiap modul
└── frontend/
    ├── .env.example
    ├── index.html
    ├── package.json
    ├── vite.config.js
    └── src/
        ├── App.jsx
        ├── index.css
        ├── main.jsx
        ├── services/
        │   └── api.js                  # Helper HTTP, Bearer token, & normalisasi JSON
        └── components/
            ├── AuthPanel.jsx
            ├── ItinerarySection.jsx
            ├── SpotSection.jsx
            └── SpotNoteSection.jsx
```

---

## How the System Works

Secara umum, sistem bekerja dengan menerima *request* dari frontend, memprosesnya melalui backend Go + Fiber v3, berinteraksi dengan database MySQL melalui GORM, lalu mengembalikan hasilnya dalam format JSON.

### 1. Gambaran Umum Sistem

```text
+----------------------+
|       USER           |
|  Pengguna Aplikasi   |
+----------+-----------+
           |
           v
+----------------------+
|      FRONTEND        |
|  React + Vite UI     |
+----------+-----------+
           |
           | HTTP Request (JSON / Bearer Token)
           v
+----------------------+
|       BACKEND        |
|    Go + Fiber v3     |
+----------+-----------+
           |
           v
+----------------------+
|   HANDLER / API      |
| Menerima Request     |
| Validasi Input       |
| Proses Data          |
+----------+-----------+
           |
           v
+----------------------+
|        GORM          |
|  Database Operations |
+----------+-----------+
           |
           v
+----------------------+
|        MySQL         |
| travel_itinerary_db  |
+----------+-----------+
           |
           | Hasil Query
           v
+----------------------+
|   BACKEND RESPONSE   |
|      JSON            |
+----------+-----------+
           |
           v
+----------------------+
|       FRONTEND       |
| Menampilkan Hasil    |
+----------------------+
```

### 2. Alur Kerja Data Aplikasi

```text
USER
 |
 | Membuat akun / Login
 v
ITINERARY
 |
 | Membuat rencana perjalanan
 v
SPOTS
 |
 | Menambahkan destinasi wisata
 v
SPOT NOTES
 |
 | Menambahkan catatan destinasi
 v
DATA TERSIMPAN DI DATABASE
```

Penjelasan setiap tahap:

- **A. User:** Pengguna melakukan registrasi atau login melalui fitur autentikasi. Data pengguna disimpan pada tabel `users`.
- **B. Itinerary:** Setelah masuk, pengguna membuat rencana perjalanan (misalnya perjalanan ke Bali selama 3 hari). Data disimpan pada tabel `itineraries` dan terhubung dengan `users` melalui `user_id`.
- **C. Spots:** Di dalam *itinerary*, pengguna menambahkan destinasi wisata yang ingin dikunjungi (contoh: Pantai Kuta, Tanah Lot, Ubud, Garuda Wisnu Kencana). Data disimpan pada tabel `spots` dan terhubung dengan `itineraries` melalui `itinerary_id`.
- **D. Spot Notes:** Setiap destinasi dapat memiliki catatan tambahan seperti tips kunjungan, jam operasional, informasi tiket masuk, dan persiapan perjalanan. Data disimpan pada tabel `spot_notes` dan terhubung dengan `spots` melalui `spot_id`.

### 3. Alur Request dan Response API

Komunikasi antara frontend dan backend menggunakan protokol HTTP dengan format JSON. Contoh alur penambahan catatan destinasi:

1. Pengguna mengisi form catatan pada frontend.
2. Frontend mengirim HTTP request `POST` ke endpoint `/api/spot-notes` (menyertakan `Authorization: Bearer <token>` bila diperlukan).
3. Backend menerima request melalui handler Spot Notes.
4. Backend memvalidasi format JSON dan input yang dikirim.
5. Backend memeriksa apakah `spot_id` tersedia di database.
6. Jika valid, GORM menyimpan data ke tabel `spot_notes`.
7. Backend memuat relasi Spot menggunakan `Preload("Spot")` jika diperlukan.
8. Backend mengembalikan response JSON dengan status `201 Created` (atau kode status yang sesuai seperti `400`, `404`, `422` bila terjadi error/validasi gagal).
9. Frontend menerima response dan memperbarui tampilan daftar catatan.

```text
USER
 |
 v
FRONTEND
 |
 | POST /api/spot-notes
 v
SPOT NOTE HANDLER
 |
 v
VALIDASI INPUT
 |
 v
CEK DATA SPOT
 |
 +---- Tidak ditemukan ----> 404 Not Found
 |
 v
SIMPAN DATA DENGAN GORM
 |
 v
MYSQL DATABASE
 |
 v
PRELOAD RELASI SPOT
 |
 v
JSON RESPONSE
 |
 v
FRONTEND MENAMPILKAN HASIL
```

### 4. Struktur Relasi Database

Urutan migrasi otomatis (`AutoMigrate`) di `backend/cmd/api/main.go`:
1. `&model.User{}`
2. `&model.Itinerary{}`
3. `&model.Spot{}`
4. `&model.SpotNote{}`

| Tabel | Primary Key | Foreign Key | Relasi |
| :--- | :--- | :--- | :--- |
| `users` | `id` | - | Induk |
| `itineraries` | `id` | `user_id` | Terhubung ke `users` |
| `spots` | `id` | `itinerary_id` | Terhubung ke `itineraries` |
| `spot_notes` | `id` | `spot_id` | Terhubung ke `spots` |

Contoh hubungan data:

```text
users
  id: 1
  name: User A
       |
       v
itineraries
  id: 1
  user_id: 1
  title: Liburan ke Bali
       |
       v
spots
  id: 1
  itinerary_id: 1
  name: Pantai Kuta
       |
       v
spot_notes
  id: 1
  spot_id: 1
  title: Tips Datang Pagi
```

### 5. Alur Kerja Backend

1. **Application Initialization:** Aplikasi dijalankan melalui `main.go`, memuat konfigurasi *environment*, menginisialisasi Fiber v3, dan mengatur CORS/middleware.
2. **Database Connection:** Backend terhubung ke MySQL menggunakan GORM (`backend/internal/config/db.go`).
3. **Auto Migration:** GORM menjalankan `AutoMigrate()` untuk menyelaraskan tabel dengan model.
4. **Middleware:** Request melewati middleware yang terdaftar (termasuk CORS dan autentikasi pada route yang diproteksi).
5. **Routing:** Fiber mengarahkan request ke handler sesuai HTTP method dan path endpoint.
6. **Handler Processing:** Handler memproses request, memvalidasi payload, dan memanggil service/repository.
7. **Database Operation:** GORM menjalankan query untuk operasi create, read, update, dan delete.
8. **JSON Response:** Backend mengembalikan respons JSON beserta HTTP status code yang sesuai.

---

## Menjalankan Aplikasi Secara Lokal

### 1. Backend (Go + Fiber v3 + MySQL)

Pastikan layanan **MySQL** aktif dan database `travel_itinerary_db` sudah tersedia, lalu jalankan dari folder `backend/`:

```bash
cd backend
go mod download
go run ./cmd/api/main.go
```

Pastikan variabel lingkungan (`APP_PORT`, `CORS_ORIGIN`, `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, serta konfigurasi JWT bila dipakai) sudah sesuai, dan `CORS_ORIGIN` mengizinkan `http://localhost:5173`.

### 2. Frontend (React + Vite)

Jalankan dari folder `frontend/`:

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

Atur `VITE_API_BASE_URL` di `frontend/.env` sesuai base URL API backend (contoh: `http://127.0.0.1:8080/api`). Aplikasi frontend akan berjalan di `http://localhost:5173`.

---

## Contoh Skenario Penggunaan Aplikasi

1. Pengguna membuka aplikasi frontend dan melakukan registrasi atau login.
2. Pengguna membuat *itinerary* dengan judul **"Liburan ke Bali"**.
3. Pengguna memilih *itinerary* tersebut dan menambahkan destinasi seperti **Pantai Kuta** dan **Tanah Lot**.
4. Pengguna memilih destinasi **Pantai Kuta** dan menambahkan catatan berjudul **"Tips Datang Pagi"** beserta jam operasional dan detail catatan.
5. Backend menyimpan catatan ke tabel `spot_notes` dengan relasi ke `spots`.
6. Saat data dimuat ulang, informasi *itinerary*, *spot*, dan catatan tampil terhubung sesuai hierarki database.

---

## Panduan Deployment (Full-Stack)

Ketentuan proyek mewajibkan aplikasi **full-stack** dan **keduanya di-deploy**:

1. **Database MySQL & Backend Go:**
   - Siapkan database MySQL yang dapat diakses oleh layanan backend.
   - *Deploy* backend Go ke platform hosting yang mendukung aplikasi Go/container, lalu atur *environment variable* produksi (`APP_PORT`, `DB_*`, rahasia JWT jika ada, dan `CORS_ORIGIN` yang mengizinkan domain frontend).
2. **Frontend React + Vite:**
   - Set `VITE_API_BASE_URL` mengarah ke URL publik backend API.
   - Jalankan `npm run build` di folder `frontend/` untuk menghasilkan aset statis di `dist/`, lalu *deploy* ke layanan hosting statis/frontend.
3. **Catatan Deployment (Isi Setelah Live):**
   - Frontend URL: `https://<url-frontend>`
   - Backend API URL: `https://<url-backend>/api`

---

## Alur Git

```bash
git status
git pull origin main
git add README.md backend frontend
git commit -m "docs: update full-stack README with architecture, setup, and deployment"
git push origin main
```

---

## Kesimpulan

**Travel Itinerary** menyatukan backend REST API berbasis **Go**, **Fiber v3**, **GORM**, dan **MySQL** dengan antarmuka **React + Vite** dalam satu struktur proyek yang terorganisir. Relasi data berjenjang (`users` → `itineraries` → `spots` → `spot_notes`) memudahkan pengelolaan rencana perjalanan secara terstruktur sekaligus mendukung integrasi dan *deployment* aplikasi secara menyeluruh.
