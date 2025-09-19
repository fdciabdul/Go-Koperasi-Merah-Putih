# Models Layer - Data Structure Definitions

## Purpose
Models layer mendefinisikan struktur data untuk aplikasi. Layer ini berisi definisi database models dengan GORM tags, request/response structs, dan data transfer objects.

## Structure Pattern
```
models/
├── postgres/
│   ├── user.go              # User dan authentication models
│   ├── koperasi.go          # Koperasi dan anggota models
│   ├── financial.go         # Financial dan accounting models
│   ├── simpan_pinjam.go     # Savings/loans models
│   ├── klinik.go            # Healthcare models
│   ├── master_data.go       # Master data models
│   └── audit.go             # Audit dan logging models
└── cassandra/               # Future analytics models
    └── analytics.go
```

## Model Definition Standards

### Base Model Pattern
```go
// Common fields untuk semua models
type BaseModel struct {
    ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// Multi-tenant model
type TenantModel struct {
    BaseModel
    TenantID uint64 `gorm:"not null;index" json:"tenant_id"`
}
```

### Model Structure
```go
type Entity struct {
    BaseModel                              // Embed base fields
    TenantID    uint64 `gorm:"not null;index" json:"tenant_id"`

    // Required fields
    Name        string `gorm:"size:255;not null" json:"name" binding:"required"`
    Email       string `gorm:"size:255;unique;not null" json:"email" binding:"required,email"`

    // Optional fields
    Description *string `gorm:"type:text" json:"description,omitempty"`
    Status      string  `gorm:"size:50;default:'active'" json:"status"`

    // Foreign keys
    CategoryID  uint64   `gorm:"not null;index" json:"category_id"`
    Category    Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`

    // One-to-many relationship
    Items       []Item `gorm:"foreignKey:EntityID" json:"items,omitempty"`

    // Validation tags
    Phone       string `gorm:"size:20" json:"phone" binding:"omitempty,numeric"`
    Website     string `gorm:"size:255" json:"website" binding:"omitempty,url"`
}
```

## GORM Tags Reference

### Primary Key & Auto Increment
```go
type Entity struct {
    ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
}
```

### Column Constraints
```go
type Entity struct {
    Name     string  `gorm:"size:255;not null;unique" json:"name"`
    Email    string  `gorm:"size:255;not null;index" json:"email"`
    Price    float64 `gorm:"precision:10;scale:2" json:"price"`
    IsActive bool    `gorm:"default:true" json:"is_active"`
    Status   string  `gorm:"size:50;default:'pending';check:status IN ('pending','active','inactive')" json:"status"`
}
```

### Relationships
```go
// One-to-One
type User struct {
    ID      uint64  `gorm:"primaryKey" json:"id"`
    Profile Profile `gorm:"foreignKey:UserID" json:"profile,omitempty"`
}

type Profile struct {
    ID     uint64 `gorm:"primaryKey" json:"id"`
    UserID uint64 `gorm:"not null;unique" json:"user_id"`
    Bio    string `gorm:"type:text" json:"bio"`
}

// One-to-Many
type Koperasi struct {
    ID      uint64            `gorm:"primaryKey" json:"id"`
    Anggota []AnggotaKoperasi `gorm:"foreignKey:KoperasiID" json:"anggota,omitempty"`
}

// Many-to-Many
type User struct {
    ID    uint64 `gorm:"primaryKey" json:"id"`
    Roles []Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}
```

### Indexes
```go
type Entity struct {
    Email    string `gorm:"index" json:"email"`                    // Simple index
    TenantID uint64 `gorm:"index:idx_tenant" json:"tenant_id"`     // Named index
    Status   string `gorm:"index:idx_tenant_status" json:"status"` // Composite index
}

// Multiple column index
type Entity struct {
    TenantID uint64 `gorm:"index:idx_tenant_entity,priority:1" json:"tenant_id"`
    EntityID uint64 `gorm:"index:idx_tenant_entity,priority:2" json:"entity_id"`
}
```

## JSON Tags & Validation

### JSON Serialization
```go
type Entity struct {
    ID          uint64     `json:"id"`
    Name        string     `json:"name"`
    Email       string     `json:"email"`
    Password    string     `json:"-"`                          // Never serialize
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty"`       // Only if not nil
    Description *string    `json:"description,omitempty"`      // Only if not nil
}
```

### Validation Tags
```go
type CreateEntityRequest struct {
    Name        string  `json:"name" binding:"required,min=3,max=255"`
    Email       string  `json:"email" binding:"required,email"`
    Phone       string  `json:"phone" binding:"omitempty,numeric,len=12"`
    Age         int     `json:"age" binding:"required,min=17,max=99"`
    Website     string  `json:"website" binding:"omitempty,url"`
    Status      string  `json:"status" binding:"required,oneof=active inactive"`
    Price       float64 `json:"price" binding:"required,min=0"`
    TenantID    uint64  `json:"tenant_id" binding:"required"`
    CategoryID  uint64  `json:"category_id" binding:"required"`
}
```

## Indonesian Models

### Koperasi Models
```go
type Koperasi struct {
    BaseModel
    TenantID         uint64 `gorm:"not null;index" json:"tenant_id"`

    // Basic information
    Nama             string `gorm:"size:255;not null" json:"nama"`
    NIAK             string `gorm:"size:16;unique;not null" json:"niak"`
    Email            string `gorm:"size:255" json:"email"`
    Telepon          string `gorm:"size:20" json:"telepon"`
    Website          string `gorm:"size:255" json:"website"`

    // Address
    Alamat           string `gorm:"type:text" json:"alamat"`
    ProvinsiID       uint64 `gorm:"not null" json:"provinsi_id"`
    KabupatenID      uint64 `gorm:"not null" json:"kabupaten_id"`
    KecamatanID      uint64 `gorm:"not null" json:"kecamatan_id"`
    KelurahanID      uint64 `gorm:"not null" json:"kelurahan_id"`
    KodePos          string `gorm:"size:10" json:"kode_pos"`

    // Classification
    JenisKoperasiID  uint64 `gorm:"not null" json:"jenis_koperasi_id"`
    BentukKoperasiID uint64 `gorm:"not null" json:"bentuk_koperasi_id"`
    KBLIID           uint64 `gorm:"not null" json:"kbli_id"`

    // Operational
    TanggalBerdiri   time.Time `json:"tanggal_berdiri"`
    IsActive         bool      `gorm:"default:true" json:"is_active"`

    // Relationships
    Provinsi         Provinsi         `gorm:"foreignKey:ProvinsiID" json:"provinsi,omitempty"`
    Kabupaten        Kabupaten        `gorm:"foreignKey:KabupatenID" json:"kabupaten,omitempty"`
    Kecamatan        Kecamatan        `gorm:"foreignKey:KecamatanID" json:"kecamatan,omitempty"`
    Kelurahan        Kelurahan        `gorm:"foreignKey:KelurahanID" json:"kelurahan,omitempty"`
    JenisKoperasi    JenisKoperasi    `gorm:"foreignKey:JenisKoperasiID" json:"jenis_koperasi,omitempty"`
    BentukKoperasi   BentukKoperasi   `gorm:"foreignKey:BentukKoperasiID" json:"bentuk_koperasi,omitempty"`
    KBLI             KBLI             `gorm:"foreignKey:KBLIID" json:"kbli,omitempty"`
    Anggota          []AnggotaKoperasi `gorm:"foreignKey:KoperasiID" json:"anggota,omitempty"`
}

type AnggotaKoperasi struct {
    BaseModel
    KoperasiID       uint64 `gorm:"not null;index" json:"koperasi_id"`

    // Member information
    NomorAnggota     string     `gorm:"size:50;not null;index:idx_koperasi_nomor,unique" json:"nomor_anggota"`
    NIK              string     `gorm:"size:16;not null" json:"nik"`
    NamaLengkap      string     `gorm:"size:255;not null" json:"nama_lengkap"`
    JenisKelamin     string     `gorm:"size:1" json:"jenis_kelamin"`
    TempatLahir      string     `gorm:"size:100" json:"tempat_lahir"`
    TanggalLahir     *time.Time `json:"tanggal_lahir"`

    // Contact information
    Alamat           string `gorm:"type:text" json:"alamat"`
    Telepon          string `gorm:"size:20" json:"telepon"`
    Email            string `gorm:"size:255" json:"email"`

    // Additional information
    Pekerjaan        string `gorm:"size:100" json:"pekerjaan"`
    StatusPernikahan string `gorm:"size:20" json:"status_pernikahan"`
    Status           string `gorm:"size:20;default:'aktif'" json:"status"`
    TanggalBergabung time.Time `json:"tanggal_bergabung"`

    // Relationships
    Koperasi         Koperasi `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
}
```

### Financial Models
```go
type COAAkun struct {
    BaseModel
    TenantID    uint64 `gorm:"not null;index" json:"tenant_id"`
    KoperasiID  uint64 `gorm:"not null;index" json:"koperasi_id"`

    KodeAkun    string `gorm:"size:20;not null;index:idx_koperasi_kode,unique" json:"kode_akun"`
    NamaAkun    string `gorm:"size:255;not null" json:"nama_akun"`
    KategoriID  uint64 `gorm:"not null" json:"kategori_id"`
    ParentID    uint64 `json:"parent_id"`
    LevelAkun   int    `gorm:"default:1" json:"level_akun"`
    SaldoNormal string `gorm:"size:10;not null" json:"saldo_normal"`
    IsKas       bool   `gorm:"default:false" json:"is_kas"`
    IsAktif     bool   `gorm:"default:true" json:"is_aktif"`

    // Relationships
    Kategori    COAKategori `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
    Parent      *COAAkun    `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
    Children    []COAAkun   `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

type JurnalUmum struct {
    BaseModel
    TenantID         uint64 `gorm:"not null;index" json:"tenant_id"`
    KoperasiID       uint64 `gorm:"not null;index" json:"koperasi_id"`

    NomorJurnal      string    `gorm:"size:50;not null;unique" json:"nomor_jurnal"`
    TanggalTransaksi time.Time `gorm:"not null;index" json:"tanggal_transaksi"`
    Referensi        string    `gorm:"size:100" json:"referensi"`
    Keterangan       string    `gorm:"type:text" json:"keterangan"`
    TotalDebit       float64   `gorm:"precision:15;scale:2" json:"total_debit"`
    TotalKredit      float64   `gorm:"precision:15;scale:2" json:"total_kredit"`
    Status           string    `gorm:"size:20;default:'draft'" json:"status"`

    // Audit fields
    CreatedBy        uint64     `gorm:"not null" json:"created_by"`
    PostedBy         uint64     `json:"posted_by"`
    PostedAt         *time.Time `json:"posted_at"`

    // Relationships
    JurnalDetail     []JurnalDetail `gorm:"foreignKey:JurnalID" json:"jurnal_detail,omitempty"`
}
```

### Healthcare Models
```go
type KlinikPasien struct {
    BaseModel
    KoperasiID      uint64 `gorm:"not null;index" json:"koperasi_id"`

    NomorRM         string     `gorm:"size:20;not null;unique" json:"nomor_rm"`
    NIK             string     `gorm:"size:16" json:"nik"`
    NamaLengkap     string     `gorm:"size:255;not null" json:"nama_lengkap"`
    JenisKelamin    string     `gorm:"size:1" json:"jenis_kelamin"`
    TempatLahir     string     `gorm:"size:100" json:"tempat_lahir"`
    TanggalLahir    *time.Time `json:"tanggal_lahir"`

    // Contact
    Alamat          string `gorm:"type:text" json:"alamat"`
    Telepon         string `gorm:"size:20" json:"telepon"`
    Email           string `gorm:"size:255" json:"email"`

    // Medical information
    GolonganDarah   string `gorm:"size:5" json:"golongan_darah"`
    Alergi          string `gorm:"type:text" json:"alergi"`
    RiwayatPenyakit string `gorm:"type:text" json:"riwayat_penyakit"`

    // Koperasi link
    AnggotaID       uint64 `json:"anggota_id"`

    // Relationships
    Anggota         *AnggotaKoperasi `gorm:"foreignKey:AnggotaID" json:"anggota,omitempty"`
    Kunjungan       []KlinikKunjungan `gorm:"foreignKey:PasienID" json:"kunjungan,omitempty"`
}
```

## Table Names & Conventions

### Table Naming
```go
// Custom table names
func (Koperasi) TableName() string {
    return "koperasi"
}

func (AnggotaKoperasi) TableName() string {
    return "anggota_koperasi"
}

func (COAAkun) TableName() string {
    return "coa_akun"
}
```

### Naming Conventions
- **Tables**: snake_case, plural (users, koperasi, anggota_koperasi)
- **Columns**: snake_case (created_at, tenant_id, nama_lengkap)
- **Foreign Keys**: {table}_id (koperasi_id, user_id)
- **Indexes**: idx_{table}_{columns} (idx_koperasi_niak, idx_tenant_status)

## Validation Patterns

### Custom Validation
```go
func (k *Koperasi) BeforeCreate(tx *gorm.DB) error {
    if len(k.NIAK) != 16 {
        return fmt.Errorf("NIAK harus 16 karakter")
    }

    if k.TanggalBerdiri.After(time.Now()) {
        return fmt.Errorf("tanggal berdiri tidak boleh di masa depan")
    }

    return nil
}

func (a *AnggotaKoperasi) BeforeSave(tx *gorm.DB) error {
    if len(a.NIK) != 16 {
        return fmt.Errorf("NIK harus 16 digit")
    }

    if a.JenisKelamin != "L" && a.JenisKelamin != "P" {
        return fmt.Errorf("jenis kelamin harus L atau P")
    }

    return nil
}
```

## Request/Response Models

### Request Models
```go
type CreateKoperasiRequest struct {
    TenantID         uint64 `json:"tenant_id" binding:"required"`
    Nama             string `json:"nama" binding:"required,min=3,max=255"`
    Email            string `json:"email" binding:"required,email"`
    Telepon          string `json:"telepon" binding:"required"`
    Alamat           string `json:"alamat" binding:"required"`
    ProvinsiID       uint64 `json:"provinsi_id" binding:"required"`
    KabupatenID      uint64 `json:"kabupaten_id" binding:"required"`
    KecamatanID      uint64 `json:"kecamatan_id" binding:"required"`
    KelurahanID      uint64 `json:"kelurahan_id" binding:"required"`
    JenisKoperasiID  uint64 `json:"jenis_koperasi_id" binding:"required"`
    BentukKoperasiID uint64 `json:"bentuk_koperasi_id" binding:"required"`
    KBLIID           uint64 `json:"kbli_id" binding:"required"`
}

type UpdateKoperasiRequest struct {
    Nama    string `json:"nama" binding:"required,min=3,max=255"`
    Email   string `json:"email" binding:"required,email"`
    Telepon string `json:"telepon" binding:"required"`
    Alamat  string `json:"alamat" binding:"required"`
    Website string `json:"website" binding:"omitempty,url"`
}
```

### Response Models
```go
type KoperasiResponse struct {
    ID               uint64    `json:"id"`
    Nama             string    `json:"nama"`
    NIAK             string    `json:"niak"`
    Email            string    `json:"email"`
    Status           string    `json:"status"`
    TanggalBerdiri   time.Time `json:"tanggal_berdiri"`
    JumlahAnggota    int       `json:"jumlah_anggota"`
    CreatedAt        time.Time `json:"created_at"`
}

type PaginatedResponse struct {
    Data  interface{} `json:"data"`
    Page  int         `json:"page"`
    Limit int         `json:"limit"`
    Total int64       `json:"total"`
}
```

## Migration Helpers

### Model Migrations
```go
func AutoMigrateModels(db *gorm.DB) error {
    return db.AutoMigrate(
        &User{},
        &Tenant{},
        &Koperasi{},
        &AnggotaKoperasi{},
        &COAAkun{},
        &JurnalUmum{},
        &JurnalDetail{},
        &KlinikPasien{},
        &KlinikKunjungan{},
        // Add all models here
    )
}
```

## Best Practices

1. **Use consistent naming conventions**
2. **Always include audit fields (created_at, updated_at)**
3. **Use appropriate GORM tags untuk constraints**
4. **Implement soft deletes untuk important data**
5. **Use foreign key constraints**
6. **Add proper indexes untuk performance**
7. **Validate data at model level**
8. **Use proper JSON tags**
9. **Handle nullable fields correctly**
10. **Document complex relationships**