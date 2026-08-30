# Selada V2 - Revamp Go API Backend

Ini adalah backend API Go yang modern, aman, dan berkinerja tinggi (dibangun menggunakan Gin + GORM + MySQL + PostgreSQL) untuk melakukan revamp pada aplikasi lama PHP CodeIgniter 3 (`seladabaselama`).

Proyek ini mencakup integrasi yang kuat, perutean multi-database, serta kontainerisasi menggunakan Docker Compose yang dirancang khusus untuk mematuhi **Aturan Keamanan & Kepatuhan SonarQube** yang ketat.

---

## 📂 Struktur Direktori Proyek

API Go ini disusun dalam arsitektur berlapis / MVC yang bersih di bawah folder `/go-api`:

```text
go-api/
├── config/
│   └── db.go            # Konfigurasi koneksi pool untuk MySQL dan PostgreSQL
├── middleware/
│   ├── auth.go          # Pembuatan token bearer kustom dan verifikasi tanda tangan (signature)
│   └── security.go      # Konfigurasi CORS dan header respons HTTP yang aman
├── models/
│   ├── user.go          # Skema GORM untuk autentikasi (User POS & PostgreSQL)
│   ├── rencana.go       # Skema GORM untuk rencana, detail rencana, dan aktivitas
│   ├── customer.go      # Skema GORM untuk manajemen data customer
│   ├── katalog.go       # Skema GORM untuk katalog motor, spesifikasi, dan harga
│   └── taksiran.go      # Skema GORM untuk estimasi/taksiran gadai berbasis PostgreSQL
├── handlers/
│   ├── auth.go          # Handler untuk login dan logout
│   ├── rencana.go       # Handler untuk list rencana, kalkulasi jumlah, input, dan soft-delete
│   ├── katalog.go       # Handler untuk perhitungan poin performa sales dan detail spesifikasi
│   └── taksiran.go      # Handler untuk kalkulasi biaya admin dan pencarian barang gadai
├── routes/
│   └── routes.go        # Pemetaan rute publik/terproteksi dan pengikatan middleware
├── Dockerfile           # Konfigurasi build kontainer produksi multi-stage (Non-root execution)
├── go.mod               # Registri dependensi Go
└── main.go              # Bootstrapper utama dan pengikatan port server
```

---

## 🔒 Arsitektur Keamanan (Kepatuhan SonarQube)

Kami telah menerapkan standar keamanan tinggi berdasarkan panduan SonarQube:

### 1. Perlindungan SQL Injection (SQLi)
- **Kerentanan:** Penggabungan string secara langsung pada query SQL (seperti `db.Raw("SELECT ... WHERE name = " + input)`) dilarang keras oleh SonarQube.
- **Solusi:** Semua query GORM menggunakan pengikatan argumen berparameter (parameterized queries) bawaan:
  ```go
  config.DBMysql.Where("username = ? AND void = ?", username, 1).First(&user)
  ```

### 2. Token Sesi Bearer yang Aman (HMAC-SHA256)
- **Kerentanan:** Sesi berbasis cookie rentan terhadap serangan CSRF. Pustaka JWT pihak ketiga yang tidak diperbarui dapat membawa celah keamanan baru.
- **Solusi:** Kami menerapkan token sesi aman kustom menggunakan pustaka standar Go (`crypto/hmac` dan `crypto/sha256`).
- **Pencegahan Timing Attack:** Verifikasi tanda tangan token menggunakan perbandingan waktu-konstan (`hmac.Equal`) untuk menggagalkan upaya peretasan berbasis perbedaan waktu respons.

### 3. Migrasi Hash Kata Sandi (BCrypt + Fallback MD5 Legacy)
- **Kerentanan:** Aplikasi CodeIgniter lama menyimpan kata sandi menggunakan MD5 (`md5($pass)`), yang sudah tidak aman secara kriptografi.
- **Solusi:** Di handler login, kami memverifikasi kata sandi menggunakan **BCrypt** terlebih dahulu. Jika gagal, sistem akan memeriksa kecocokan menggunakan hash **MD5**. Jika MD5 cocok, sistem akan **secara otomatis memperbarui hash kata sandi tersebut menjadi BCrypt** di database MySQL POS.

### 4. Pembatasan CORS & Pengerasan Header Respons HTTP
- Menghindari penggunaan header wildcard dinamis `Access-Control-Allow-Origin: *` pada rute terproteksi. CORS secara ketat memverifikasi asal permintaan (Origin) terhadap daftar putih (whitelist).
- Menambahkan header keamanan berikut secara otomatis pada setiap respons HTTP:
  - `X-Frame-Options: DENY` (mencegah Clickjacking)
  - `X-Content-Type-Options: nosniff` (mencegah MIME sniffing)
  - `X-XSS-Protection: 1; mode=block` (mengaktifkan filter XSS bawaan browser)
  - `Content-Security-Policy` & `Strict-Transport-Security`

### 5. Pengerasan Keamanan Kontainer (Container Hardening)
- `Dockerfile` menggunakan teknik **multi-stage build** untuk membuang compiler compiler Go dari gambar akhir, sehingga menghasilkan kontainer Alpine yang sangat minimal.
- Kontainer dikonfigurasi untuk berjalan di bawah akun sistem non-root kustom (`appuser:appgroup`), meminimalisir risiko eskalasi hak akses (privilege escalation).

---

## 🛢️ Perutean Multi-Database & Konfigurasi Lingkungan

API backend ini terhubung ke tiga koneksi database sekaligus:
1. **MySQL Default DB (`selada`)** -> Menyimpan informasi customer, rencana aktivitas sales, dan poin performa.
2. **PostgreSQL DB (`sam_live`)** -> Menyimpan data jenis barang, jenis produk gadai, dan tarif biaya admin.
3. **MySQL POS DB (`sam_pos`)** -> Menyimpan data kredensial user/sales dan katalog sepeda motor.

### Variabel Lingkungan (Environment Variables)
Sesuaikan variabel ini pada lingkungan lokal Anda atau file `docker-compose.yml`:

| Variabel Lingkungan | Deskripsi | Nilai Default |
| :--- | :--- | :--- |
| `PORT` | Port aplikasi Go | `3000` |
| `JWT_SECRET` | Kunci rahasia untuk enkripsi token sesi | `super-secret-key-selada-v2-123456789` |
| `DB_MYSQL_HOST` | Host database MySQL default | `host.docker.internal` (Laragon Host) |
| `DB_MYSQL_USER` | Username database MySQL default | `root` |
| `DB_MYSQL_PASS` | Password database MySQL default | `aditya@123` |
| `DB_MYSQL_NAME` | Nama database MySQL default | `selada` |
| `DB_PG_HOST` | Host database PostgreSQL | `host.docker.internal` (Laragon Host) |
| `DB_PG_PORT` | Port database PostgreSQL | `5432` |
| `DB_PG_USER` | Username database PostgreSQL | `postgres` |
| `DB_PG_PASS` | Password database PostgreSQL | `<uX42)#YQ#8D;t]/` |
| `DB_PG_NAME` | Nama database PostgreSQL | `sam_live` |
| `DB_MYSQL_POS_HOST`| Host database MySQL POS | `host.docker.internal` (Laragon Host) |
| `DB_MYSQL_POS_NAME`| Nama database MySQL POS | `sam_pos` |

---

## ⚡ Katalog Endpoint API

Semua rute memiliki awalan `/api`. Endpoint terproteksi memerlukan token pada header `Authorization` dengan format `Bearer <Token>`.

### Endpoint Publik
- **`POST /api/login`**: Menerima JSON `{"user": "...", "pass": "..."}`. Mengembalikan token sesi kustom dan data profil user.
- **`POST /api/logout`**: Mengembalikan pesan sukses logout.

### Endpoint Terproteksi (Memerlukan Token Bearer)
#### 📊 Rencana & Inquiry
- **`GET /api/inquiry/counts`**: Mengembalikan data perhitungan jumlah untuk indikator badge dashboard (leads, prospect, hotprospect, sbg, batal).
- **`GET /api/inquiry/filter`**: List data rencana terpaginasi lengkap dengan filter status dan fitur pencarian nama customer.
- **`POST /api/rencana/input`**: Menambahkan rencana baru. Melakukan normalisasi nomor telepon (mengubah awalan `62` ke `0`), pengecekan duplikasi global, dan validasi kepemilikan.
- **`DELETE /api/rencana/:id`**: Soft-delete data rencana dengan mengubah status detil rencana menjadi tidak aktif.

#### 🏍️ Katalog Sepeda Motor
- **`GET /api/katalog/point`**: Menghitung akumulasi poin performa aktivitas sales yang sedang login.
- **`GET /api/katalog/list/:category`**: Mengembalikan list motor (sport, trail, moped) beserta gambar utamanya.
- **`GET /api/katalog/spec/:id`**: Mengembalikan spesifikasi lengkap motor, daftar brosur, pilihan warna, detail fitur, dan harga.

#### ⚖️ Estimasi / Taksiran Gadai
- **`GET /api/taksiran/jenis-barang`**: Mengembalikan list jenis barang gadai yang aktif.
- **`GET /api/taksiran/barang-umum`**: Melakukan pencarian spesifik detail barang umum berdasarkan tipe, grade, dan kecocokan nama.
- **`GET /api/taksiran/barang-emas`**: Mengembalikan informasi harga taksir emas yang berlaku.
- **`GET /api/taksiran/products`**: Mengembalikan daftar produk berdasarkan tipe (misal: gadai emas vs gadai umum).
- **`GET /api/taksiran/biaya-admin`**: Menghitung secara otomatis biaya administrasi gadai berdasarkan jenis produk dan nilai pinjaman.
- **`POST /api/taksiran/simpan`**: Menyimpan hasil estimasi transaksi gadai.

---

## 🐳 Cara Menjalankan dengan Docker Compose

Kontainerisasi dengan Docker Compose adalah cara tercepat dan paling konsisten untuk menjalankan proyek revamp ini.

### Prasyarat: Terhubung ke Database Host (Laragon)
If your databases (MySQL and PostgreSQL) are running locally via Laragon on your host machine, the Docker container bridges connections automatically to your host network using the `host.docker.internal` gateway address, configured in `docker-compose.yml`.

### Menjalankan Kontainer
Jalankan perintah berikut di direktori root proyek (`d:\02_Projects\selada-v2`):

```bash
# Build image API Go dan jalankan kontainer di latar belakang (detached mode)
docker compose up --build -d
```

Melihat status kontainer yang sedang berjalan:
```bash
docker compose ps
```

Melihat log aplikasi backend secara real-time:
```bash
docker compose logs -f backend
```

Menghentikan semua layanan kontainer:
```bash
docker compose down
```

---

## 🛠️ Menjalankan Secara Lokal (Tanpa Docker)

Jika Anda memiliki compiler Go (Go versi >= 1.23) di komputer Anda:

```bash
# Masuk ke folder go-api
cd go-api

# Unduh dan rapikan dependensi proyek
go mod tidy

# Jalankan server backend Go
go run main.go
```
Aplikasi akan aktif dan mendengarkan permintaan di `http://localhost:3000`.
