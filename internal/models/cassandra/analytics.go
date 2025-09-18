package cassandra

import (
	"time"

	"github.com/gocql/gocql"
)

type FactKeuanganBulanan struct {
	ID                      gocql.UUID `json:"id"`
	TenantID                uint64     `json:"tenant_id"`
	KoperasiID              uint64     `json:"koperasi_id"`
	PeriodeTahun            int        `json:"periode_tahun"`
	PeriodeBulan            int        `json:"periode_bulan"`
	TotalSimpanan           float64    `json:"total_simpanan"`
	TotalPinjaman           float64    `json:"total_pinjaman"`
	TotalAngsuran           float64    `json:"total_angsuran"`
	TotalTransaksiPPOB      float64    `json:"total_transaksi_ppob"`
	TotalFeePPOB            float64    `json:"total_fee_ppob"`
	JumlahTransaksiPPOB     int        `json:"jumlah_transaksi_ppob"`
	TotalPendapatanKlinik   float64    `json:"total_pendapatan_klinik"`
	JumlahKunjunganKlinik   int        `json:"jumlah_kunjungan_klinik"`
	TotalAset               float64    `json:"total_aset"`
	TotalKewajiban          float64    `json:"total_kewajiban"`
	TotalEkuitas            float64    `json:"total_ekuitas"`
	TotalPendapatan         float64    `json:"total_pendapatan"`
	TotalBeban              float64    `json:"total_beban"`
	CreatedAt               time.Time  `json:"created_at"`
}

type FactAnggotaBulanan struct {
	ID                    gocql.UUID `json:"id"`
	KoperasiID            uint64     `json:"koperasi_id"`
	PeriodeTahun          int        `json:"periode_tahun"`
	PeriodeBulan          int        `json:"periode_bulan"`
	JumlahAnggotaAktif    int        `json:"jumlah_anggota_aktif"`
	JumlahAnggotaBaru     int        `json:"jumlah_anggota_baru"`
	JumlahAnggotaKeluar   int        `json:"jumlah_anggota_keluar"`
	JumlahPengurus        int        `json:"jumlah_pengurus"`
	JumlahPengawas        int        `json:"jumlah_pengawas"`
	CreatedAt             time.Time  `json:"created_at"`
}

type TransactionLog struct {
	ID               gocql.UUID `json:"id"`
	TenantID         uint64     `json:"tenant_id"`
	KoperasiID       uint64     `json:"koperasi_id"`
	TransactionType  string     `json:"transaction_type"`
	TransactionID    uint64     `json:"transaction_id"`
	UserID           uint64     `json:"user_id"`
	Amount           float64    `json:"amount"`
	Status           string     `json:"status"`
	Description      string     `json:"description"`
	Metadata         string     `json:"metadata"`
	IPAddress        string     `json:"ip_address"`
	UserAgent        string     `json:"user_agent"`
	CreatedAt        time.Time  `json:"created_at"`
	Year             int        `json:"year"`
	Month            int        `json:"month"`
	Day              int        `json:"day"`
}

type PaymentAnalytics struct {
	ID                 gocql.UUID `json:"id"`
	TenantID           uint64     `json:"tenant_id"`
	KoperasiID         uint64     `json:"koperasi_id"`
	PaymentProvider    string     `json:"payment_provider"`
	PaymentMethod      string     `json:"payment_method"`
	TransactionType    string     `json:"transaction_type"`
	TotalAmount        float64    `json:"total_amount"`
	AdminFee           float64    `json:"admin_fee"`
	TransactionCount   int        `json:"transaction_count"`
	SuccessCount       int        `json:"success_count"`
	FailedCount        int        `json:"failed_count"`
	SuccessRate        float64    `json:"success_rate"`
	CreatedAt          time.Time  `json:"created_at"`
	Year               int        `json:"year"`
	Month              int        `json:"month"`
	Day                int        `json:"day"`
}

type PPOBAnalytics struct {
	ID                gocql.UUID `json:"id"`
	TenantID          uint64     `json:"tenant_id"`
	KoperasiID        uint64     `json:"koperasi_id"`
	ProductCategory   string     `json:"product_category"`
	ProductName       string     `json:"product_name"`
	Provider          string     `json:"provider"`
	TransactionCount  int        `json:"transaction_count"`
	TotalRevenue      float64    `json:"total_revenue"`
	TotalCommission   float64    `json:"total_commission"`
	SuccessCount      int        `json:"success_count"`
	FailedCount       int        `json:"failed_count"`
	SuccessRate       float64    `json:"success_rate"`
	CreatedAt         time.Time  `json:"created_at"`
	Year              int        `json:"year"`
	Month             int        `json:"month"`
	Day               int        `json:"day"`
}

type KlinikAnalytics struct {
	ID                    gocql.UUID `json:"id"`
	TenantID              uint64     `json:"tenant_id"`
	KoperasiID            uint64     `json:"koperasi_id"`
	DokterID              uint64     `json:"dokter_id"`
	DokterNama            string     `json:"dokter_nama"`
	Spesialisasi          string     `json:"spesialisasi"`
	JumlahKunjungan       int        `json:"jumlah_kunjungan"`
	TotalPendapatan       float64    `json:"total_pendapatan"`
	RataRataBiayaKunjungan float64   `json:"rata_rata_biaya_kunjungan"`
	CreatedAt             time.Time  `json:"created_at"`
	Year                  int        `json:"year"`
	Month                 int        `json:"month"`
	Day                   int        `json:"day"`
}

type UserActivityLog struct {
	ID               gocql.UUID `json:"id"`
	TenantID         uint64     `json:"tenant_id"`
	KoperasiID       uint64     `json:"koperasi_id"`
	UserID           uint64     `json:"user_id"`
	Username         string     `json:"username"`
	Activity         string     `json:"activity"`
	Module           string     `json:"module"`
	Description      string     `json:"description"`
	IPAddress        string     `json:"ip_address"`
	UserAgent        string     `json:"user_agent"`
	SessionID        string     `json:"session_id"`
	CreatedAt        time.Time  `json:"created_at"`
	Year             int        `json:"year"`
	Month            int        `json:"month"`
	Day              int        `json:"day"`
	Hour             int        `json:"hour"`
}

type ErrorLog struct {
	ID               gocql.UUID `json:"id"`
	TenantID         uint64     `json:"tenant_id"`
	KoperasiID       uint64     `json:"koperasi_id"`
	UserID           uint64     `json:"user_id"`
	ErrorType        string     `json:"error_type"`
	ErrorMessage     string     `json:"error_message"`
	StackTrace       string     `json:"stack_trace"`
	RequestData      string     `json:"request_data"`
	ResponseData     string     `json:"response_data"`
	Module           string     `json:"module"`
	Function         string     `json:"function"`
	IPAddress        string     `json:"ip_address"`
	UserAgent        string     `json:"user_agent"`
	CreatedAt        time.Time  `json:"created_at"`
	Year             int        `json:"year"`
	Month            int        `json:"month"`
	Day              int        `json:"day"`
}

type PerformanceMetrics struct {
	ID               gocql.UUID `json:"id"`
	TenantID         uint64     `json:"tenant_id"`
	KoperasiID       uint64     `json:"koperasi_id"`
	MetricName       string     `json:"metric_name"`
	MetricValue      float64    `json:"metric_value"`
	MetricUnit       string     `json:"metric_unit"`
	EndpointPath     string     `json:"endpoint_path"`
	HttpMethod       string     `json:"http_method"`
	ResponseTime     float64    `json:"response_time"`
	StatusCode       int        `json:"status_code"`
	RequestSize      int64      `json:"request_size"`
	ResponseSize     int64      `json:"response_size"`
	CreatedAt        time.Time  `json:"created_at"`
	Year             int        `json:"year"`
	Month            int        `json:"month"`
	Day              int        `json:"day"`
	Hour             int        `json:"hour"`
	Minute           int        `json:"minute"`
}