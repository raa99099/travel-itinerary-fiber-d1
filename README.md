## Project Overview

**Travel Itinerary Backend API** merupakan backend dari aplikasi Travel Itinerary yang berfungsi untuk mengelola data dan kebutuhan perjalanan wisata pengguna.

Aplikasi ini dibuat untuk membantu pengguna dalam menyusun rencana perjalanan secara lebih terstruktur, mulai dari membuat akun, menyusun itinerary perjalanan, menambahkan destinasi wisata, hingga menyimpan catatan penting pada setiap destinasi.

Backend ini berperan sebagai penghubung antara frontend (tampilan aplikasi) dengan database. Seluruh proses pengelolaan data dilakukan melalui REST API yang dibangun menggunakan Go dan Fiber v3, dengan GORM sebagai ORM dan MySQL sebagai database.

### Tujuan Project

Project ini dikembangkan dengan beberapa tujuan utama:

- Mempermudah pengguna dalam mengelola rencana perjalanan.
- Menyimpan daftar destinasi wisata dalam setiap itinerary.
- Menyediakan fitur catatan untuk setiap destinasi wisata.
- Mengelola data pengguna dan relasi antar data secara terstruktur.
- Menyediakan REST API yang dapat diintegrasikan dengan frontend.
- Menerapkan konsep backend development menggunakan Go, Fiber, GORM, dan MySQL.

### Fitur Utama

| Fitur | Deskripsi |
|---|---|
| User Management | Mengelola data pengguna aplikasi |
| Authentication | Menangani proses autentikasi pengguna |
| Itinerary Management | Membuat, melihat, mengubah, dan menghapus rencana perjalanan |
| Spot Management | Mengelola daftar destinasi wisata dalam itinerary |
| Spot Notes | Menambahkan dan mengelola catatan pada setiap destinasi wisata |

---

## How the System Works

Secara umum, sistem bekerja dengan menerima request dari frontend, memproses request melalui backend, berinteraksi dengan database MySQL, kemudian mengembalikan hasilnya dalam format JSON.

### 1. Gambaran Umum Sistem

Berikut gambaran komunikasi antarbagian aplikasi:

```text
+----------------------+
|       USER           |
|  Pengguna Aplikasi   |
+----------+-----------+
           |
           v
+----------------------+
|      FRONTEND        |
|  Tampilan Aplikasi   |
+----------+-----------+
           |
           | HTTP Request
           | (JSON)
           v
+----------------------+
|       BACKEND        |
|       Go + Fiber     |
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

Pengelolaan data dalam aplikasi mengikuti struktur relasi berikut:

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

**A. User**

Pengguna melakukan registrasi atau login melalui fitur autentikasi. Data pengguna disimpan dalam tabel `users`.

**B. Itinerary**

Setelah pengguna berhasil masuk, pengguna dapat membuat rencana perjalanan. Contohnya, pengguna membuat itinerary untuk perjalanan ke Bali selama 3 hari.

Data itinerary disimpan dalam tabel `itineraries` dan terhubung dengan user melalui relasi foreign key.

**C. Spots**

Di dalam itinerary, pengguna dapat menambahkan beberapa destinasi wisata yang ingin dikunjungi.

Contohnya:
- Pantai Kuta
- Tanah Lot
- Ubud
- Garuda Wisnu Kencana

Data destinasi disimpan dalam tabel `spots` dan terhubung dengan itinerary.

**D. Spot Notes**

Setiap destinasi dapat memiliki catatan tambahan yang berisi informasi penting, seperti:
- Tips mengunjungi tempat wisata.
- Jam operasional.
- Informasi tiket masuk.
- Hal-hal yang perlu dipersiapkan sebelum berkunjung.

Data catatan disimpan dalam tabel `spot_notes` dan terhubung dengan spot melalui `spot_id`.

### 3. Alur Request dan Response API

Setiap komunikasi antara frontend dan backend menggunakan protokol HTTP dengan format JSON.

Contoh alur ketika pengguna menambahkan catatan destinasi:

1. Pengguna mengisi form catatan pada frontend.
2. Frontend mengirim HTTP request `POST` ke endpoint `/api/spot-notes`.
3. Backend menerima request melalui handler Spot Notes.
4. Backend memeriksa format JSON dan melakukan validasi input.
5. Backend memeriksa apakah `spot_id` yang dikirim tersedia di database.
6. Jika valid, GORM menyimpan data ke tabel `spot_notes`.
7. Backend mengambil relasi Spot menggunakan `Preload("Spot")` jika diperlukan.
8. Backend mengirim response JSON dengan status `201 Created`.
9. Frontend menerima response dan menampilkan informasi bahwa catatan berhasil dibuat.

**Ilustrasi alur:**

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

Database menggunakan konsep relasi antar tabel untuk menjaga keterhubungan data.

| Tabel | Primary Key | Foreign Key | Relasi |
|---|---|---|---|
| `users` | `id` | - | Induk |
| `itineraries` | `id` | `user_id` | Terhubung ke Users |
| `spots` | `id` | `itinerary_id` | Terhubung ke Itineraries |
| `spot_notes` | `id` | `spot_id` | Terhubung ke Spots |

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

Dengan struktur ini, satu pengguna dapat memiliki beberapa itinerary, satu itinerary dapat memiliki beberapa destinasi, dan satu destinasi dapat memiliki beberapa catatan.

### 5. Alur Kerja Backend

Dari sisi teknis, proses backend berjalan melalui beberapa tahapan:

1. **Application Initialization**  
   Aplikasi dijalankan melalui `main.go`, kemudian memuat konfigurasi environment dan menginisialisasi Fiber.

2. **Database Connection**  
   Backend membuat koneksi ke MySQL menggunakan GORM.

3. **Auto Migration**  
   GORM menjalankan `AutoMigrate()` untuk memastikan tabel sesuai dengan struktur model.

4. **Middleware**  
   Request melewati middleware yang digunakan, termasuk middleware autentikasi jika diterapkan pada route terkait.

5. **Routing**  
   Fiber mengarahkan request ke handler berdasarkan HTTP method dan endpoint.

6. **Handler Processing**  
   Handler memproses request, memvalidasi input, dan menjalankan operasi database.

7. **Database Operation**  
   GORM menjalankan query untuk mengambil, menambahkan, memperbarui, atau menghapus data.

8. **JSON Response**  
   Backend mengembalikan hasil pemrosesan kepada frontend dalam format JSON beserta HTTP status code yang sesuai.

---

## Contoh Skenario Penggunaan Aplikasi

Misalnya seorang pengguna ingin membuat perjalanan wisata ke Bali.

1. Pengguna melakukan registrasi dan login.
2. Pengguna membuat itinerary dengan judul "Liburan ke Bali".
3. Pengguna menambahkan beberapa destinasi seperti Pantai Kuta dan Tanah Lot.
4. Pengguna membuka detail destinasi Pantai Kuta.
5. Pengguna menambahkan catatan berjudul "Tips Datang Pagi".
6. Backend menyimpan catatan tersebut dan menghubungkannya dengan Pantai Kuta.
7. Ketika pengguna membuka kembali detail destinasi, backend mengambil informasi spot beserta catatan yang berkaitan.

Dengan alur tersebut, seluruh informasi perjalanan dapat dikelola secara terstruktur dan saling terhubung melalui database.

---

## Kesimpulan

Travel Itinerary Backend API merupakan sistem backend yang dirancang untuk mendukung pengelolaan rencana perjalanan wisata secara terstruktur.

Dengan menggunakan Go, Fiber v3, GORM, dan MySQL, aplikasi ini menyediakan layanan REST API untuk mengelola pengguna, itinerary, destinasi wisata, serta catatan pada setiap destinasi.

Struktur data yang berjenjang memungkinkan setiap informasi saling terhubung melalui relasi database, sehingga mempermudah pengelolaan data dan integrasi dengan frontend.