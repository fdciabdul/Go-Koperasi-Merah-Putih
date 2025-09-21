# Sistem Manajemen Koperasi

Aplikasi manajemen koperasi berbasis Go untuk operasional koperasi di Indonesia dengan fitur lengkap dan arsitektur yang mudah dipelihara.

## Fitur Utama

| Modul | Deskripsi | Status |
|--------|-------------|--------|
| **Autentikasi** | Registrasi user, login, JWT auth | Selesai |
| **Manajemen Koperasi** | CRUD koperasi, kelola anggota | Selesai |
| **Simpan Pinjam** | Produk simpanan & pinjaman, transaksi | Selesai |
| **Kelola Produk** | Inventori, supplier, penjualan, pembelian | Selesai |
| **Klinik** | Layanan kesehatan, pasien, obat-obatan | Selesai |
| **Keuangan** | Chart of accounts, jurnal, laporan | Selesai |
| **PPOB** | Layanan Payment Point Online Bank | Selesai |
| **Payment Gateway** | Integrasi Midtrans & Xendit | Selesai |
| **Analytics** | Analytics berbasis Cassandra | Dalam Proses |
| **Audit Logging** | Pencatatan audit sistem lengkap | Selesai |

## Teknologi

### Stack Teknologi

| Komponen | Teknologi | Fungsi |
|-----------|------------|---------|
| **Backend** | Go 1.19+, Gin Web Framework | REST API server |
| **Database** | PostgreSQL 13+ | Penyimpanan data utama |
| **Analytics** | Apache Cassandra | Analytics big data |
| **Cache** | Redis (opsional) | Manajemen session & cache |
| **ORM** | GORM v2 | Operasi database |
| **Authentication** | JWT + bcrypt | Layer keamanan |
| **Payment** | Midtrans, Xendit | Pemrosesan pembayaran |

## Instalasi dan Setup

### Persyaratan

| Kebutuhan | Versi | Instalasi |
|-------------|---------|--------------|
| **Go** | 1.19+ | [Download Go](https://golang.org/dl/) |
| **PostgreSQL** | 13+ | [Download PostgreSQL](https://www.postgresql.org/download/) |
| **Git** | Terbaru | [Download Git](https://git-scm.com/downloads) |

### Panduan Instalasi

#### Untuk Unix/Linux/macOS:
```bash
# 1. Clone Repository
git clone https://github.com/fdciabdul/Go-Koperasi-Merah-Putih
cd go_koperasi

# 2. Install Dependencies
go mod download

# 3. Konfigurasi Environment
cp .env.example .env
# Edit .env dengan kredensial database Anda

# 4. Setup Database
createdb koperasi_db
make migrate-fresh

# 5. Jalankan Aplikasi
make run
```

#### Untuk Windows:
```cmd
REM 1. Clone Repository
git clone https://github.com/fdciabdul/Go-Koperasi-Merah-Putih
cd go_koperasi

REM 2. Install Dependencies
go mod download

REM 3. Konfigurasi Environment
copy .env.example .env
REM Edit .env dengan kredensial database Anda

REM 4. Setup Database
createdb koperasi_db
make.bat migrate-fresh

REM 5. Jalankan Aplikasi
make.bat run
```

#### Setup Satu Perintah:
```bash
# Unix/Linux/macOS
make quick-start

# Windows
make.bat quick-start
```

## Perintah yang Tersedia

Proyek ini mendukung environment **Unix/Linux/macOS** (Makefile) dan **Windows** (make.bat).

### Cara Penggunaan

| Platform | Penggunaan | Contoh |
|----------|-------|---------|
| **Unix/Linux/macOS** | `make <perintah>` | `make run` |
| **Windows** | `make.bat <perintah>` | `make.bat run` |

### Perintah Development

| Perintah | Deskripsi | Unix | Windows |
|---------|-------------|------|---------|
| `help` | Tampilkan semua perintah yang tersedia | `make help` | `make.bat help` |
| `build` | Build aplikasi | `make build` | `make.bat build` |
| `run` | Jalankan aplikasi | `make run` | `make.bat run` |
| `test` | Jalankan semua test | `make test` | `make.bat test` |
| `fmt` | Format kode | `make fmt` | `make.bat fmt` |
| `lint` | Lint kode | `make lint` | `make.bat lint` |
| `dev` | Hot reload development | `make dev` | `make.bat dev` |

### Perintah Database

| Perintah | Deskripsi | Unix | Windows |
|---------|-------------|------|---------|
| `migrate` | Jalankan GORM auto-migrations | `make migrate` | `make.bat migrate` |
| `seed` | Jalankan database seeders | `make seed` | `make.bat seed` |
| `migrate-fresh` | Drop, migrate, dan seed | `make migrate-fresh` | `make.bat migrate-fresh` |
| `migrate-drop` | Drop semua tabel dan migrate | `make migrate-drop` | `make.bat migrate-drop` |
| `dev-setup` | Setup development lengkap | `make dev-setup` | `make.bat dev-setup` |

### Perintah Tools

| Perintah | Deskripsi | Unix | Windows |
|---------|-------------|------|---------|
| `install-tools` | Install development tools | `make install-tools` | `make.bat install-tools` |
| `quick-start` | Setup lengkap untuk developer baru | `make quick-start` | `make.bat quick-start` |
| `clean` | Bersihkan build artifacts | `make clean` | `make.bat clean` |
| `env-info` | Tampilkan informasi environment | `make env-info` | `make.bat env-info` |

## Migrasi Database

### GORM Auto-Migration

Proyek ini menggunakan **GORM auto-migration** daripada file SQL:

```go
// Jalankan migrasi
go run cmd/migrate/main.go

// Fresh migration (drop + migrate + seed)
go run cmd/migrate/main.go -fresh

// Drop tabel dan migrate
go run cmd/migrate/main.go -drop
```

### Fitur Migrasi

| Fitur | Deskripsi |
|---------|-------------|
| **Auto-Migration** | GORM otomatis membuat/update tabel |
| **Berbasis Model** | Migrasi berdasarkan Go struct models |
| **Pembuatan Index** | Pembuatan index otomatis untuk performa |
| **Penambahan Constraint** | Custom business rule constraints |
| **Integrasi Seeder** | Seeding otomatis setelah migrasi |

### Opsi Perintah Migrasi

| Flag | Deskripsi | Contoh |
|------|-------------|---------|
| `-drop` | Drop semua tabel sebelum migrasi | `go run cmd/migrate/main.go -drop` |
| `-seed` | Jalankan seeders setelah migrasi | `go run cmd/migrate/main.go -seed` |
| `-fresh` | Drop, migrate, dan seed | `go run cmd/migrate/main.go -fresh` |

## API Endpoints

### Authentication

| Method | Endpoint | Deskripsi | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/v1/users/register` | Registrasi user | Public |
| `POST` | `/api/v1/auth/login` | Login user | Public |
| `PUT` | `/api/v1/users/verify-payment/:id` | Verifikasi pembayaran | Public |

### Manajemen Koperasi

| Method | Endpoint | Deskripsi | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/v1/koperasi` | Buat koperasi | SuperAdmin |
| `GET` | `/api/v1/koperasi` | List koperasi | Authenticated |
| `GET` | `/api/v1/koperasi/:id` | Detail koperasi | Authenticated |
| `PUT` | `/api/v1/koperasi/:id` | Update koperasi | Admin |

### Manajemen Produk

| Method | Endpoint | Deskripsi | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/v1/produk` | Buat produk | Admin |
| `GET` | `/api/v1/produk/:koperasi_id` | List produk | Authenticated |
| `POST` | `/api/v1/produk/purchase-order` | Buat purchase order | Admin |
| `POST` | `/api/v1/produk/penjualan` | Buat transaksi penjualan | Authenticated |

### Manajemen Keuangan

| Method | Endpoint | Deskripsi | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/v1/financial/jurnal` | Buat jurnal entry | Financial |
| `GET` | `/api/v1/financial/:id/neraca-saldo` | Neraca saldo | Financial |
| `GET` | `/api/v1/financial/:id/laba-rugi` | Laba rugi | Financial |

## Arsitektur Routes Modular

Routes diorganisir dalam struktur modular untuk kemudahan maintenance:

```
internal/routes/
├── routes.go              # Main orchestrator
└── modules/               # Domain-specific modules
    ├── auth_routes.go     # Authentication & payments
    ├── koperasi_routes.go # Manajemen koperasi
    ├── produk_routes.go   # Manajemen produk
    ├── financial_routes.go # Operasi keuangan
    └── ...                # Route domain lainnya
```

### Keuntungan

| Keuntungan | Deskripsi |
|---------|-------------|
| **Separation of Concerns** | Setiap domain punya file route sendiri |
| **Maintainability** | Mudah mencari dan memodifikasi endpoint |
| **Scalability** | Mudah menambah modul domain baru |
| **Kolaborasi Tim** | Mengurangi konflik saat bekerja pada fitur berbeda |

## Authentication & Authorization

### Role-Based Access Control (RBAC)

| Role | Permission | Level Akses |
|------|-------------|--------------|
| **SuperAdmin** | Akses sistem penuh | Semua operasi |
| **Admin** | Manajemen koperasi | Spesifik koperasi |
| **Financial** | Operasi keuangan | Modul keuangan |
| **User** | Operasi dasar | Akses terbatas |
| **Operator** | Entry data | Modul spesifik |

## Fitur Bisnis

### Kepatuhan Indonesia

| Fitur | Deskripsi | Implementasi |
|---------|-------------|----------------|
| **NIAK Generation** | ID koperasi otomatis | Berbasis algoritma |
| **Validasi NIK** | Validasi ID Indonesia | Validasi 16 digit |
| **Data Regional** | Data wilayah Indonesia lengkap | Provinsi hingga Kelurahan |
| **Integrasi KBLI** | Klasifikasi bisnis | Kepatuhan standar |

### Manajemen Produk

| Fitur | Deskripsi | Keuntungan |
|---------|-------------|----------|
| **12 Kategori Produk** | Makanan, minuman, ternak, dll | Inventori terorganisir |
| **Support Barcode** | Generasi EAN-13 | Tracking efisien |
| **Manajemen Supplier** | Dukungan multi-supplier | Optimasi biaya |
| **Tracking Kedaluwarsa** | Manajemen tanggal expire | Pengurangan waste |
| **Purchase Orders** | Workflow procurement lengkap | Pembelian terorganisir |
| **Transaksi Penjualan** | Pemrosesan penjualan style POS | Transaksi mudah |
| **Pergerakan Stok** | Tracking inventori real-time | Level stok akurat |

## Testing

```bash
# Jalankan semua test
make test

# Jalankan test package spesifik
go test ./internal/services/... -v

# Jalankan dengan coverage
go test ./... -cover
```

## Deployment

### Environment Variables

| Variable | Deskripsi | Contoh |
|----------|-------------|---------|
| `DB_HOST` | Host database | `localhost` |
| `DB_PORT` | Port database | `5432` |
| `DB_NAME` | Nama database | `koperasi_db` |
| `DB_USER` | User database | `postgres` |
| `DB_PASSWORD` | Password database | `password` |
| `JWT_SECRET` | JWT signing key | `your-secret-key` |

## Kontribusi

1. Fork repository ini
2. Buat feature branch (`git checkout -b feature/fitur-keren`)
3. Commit perubahan (`git commit -m 'Tambah fitur keren'`)
4. Push ke branch (`git push origin feature/fitur-keren`)
5. Buka Pull Request

## Lisensi

Proyek ini menggunakan lisensi MIT License.

## Kontak

- Repository: https://github.com/fdciabdul/Go-Koperasi-Merah-Putih
- Telegram: cp@imtaqin.id
- Issues: GitHub Issues
- Diskusi: GitHub Discussions