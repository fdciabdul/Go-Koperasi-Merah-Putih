package postgres

import (
	"time"

	"gorm.io/gorm"
)

type KategoriProduk struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Kode      string         `gorm:"uniqueIndex;size:20;not null" json:"kode"`
	Nama      string         `gorm:"size:100;not null" json:"nama"`
	Deskripsi string         `gorm:"type:text" json:"deskripsi"`
	Icon      string         `gorm:"size:255" json:"icon"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Produk []Produk `gorm:"foreignKey:KategoriProdukID" json:"produk,omitempty"`
}

type SatuanProduk struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Kode      string         `gorm:"uniqueIndex;size:10;not null" json:"kode"`
	Nama      string         `gorm:"size:50;not null" json:"nama"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Produk []Produk `gorm:"foreignKey:SatuanProdukID" json:"produk,omitempty"`
}

type Supplier struct {
	ID             uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID     uint64         `gorm:"not null;index" json:"koperasi_id"`
	Kode           string         `gorm:"size:20;not null" json:"kode"`
	Nama           string         `gorm:"size:255;not null" json:"nama"`
	KontakPerson   string         `gorm:"size:100" json:"kontak_person"`
	Telepon        string         `gorm:"size:20" json:"telepon"`
	Email          string         `gorm:"size:255" json:"email"`
	Alamat         string         `gorm:"type:text" json:"alamat"`
	ProvinsiID     uint64         `json:"provinsi_id"`
	KabupatenID    uint64         `json:"kabupaten_id"`
	NoRekening     string         `gorm:"size:50" json:"no_rekening"`
	NamaBank       string         `gorm:"size:100" json:"nama_bank"`
	AtasNamaBank   string         `gorm:"size:100" json:"atas_nama_bank"`
	NPWP           string         `gorm:"size:20" json:"npwp"`
	JenisSupplier  string         `gorm:"type:varchar(20);default:'individu'" json:"jenis_supplier"`
	Status         string         `gorm:"type:varchar(20);default:'aktif'" json:"status"`
	TermPembayaran int            `gorm:"default:30;comment:hari" json:"term_pembayaran"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedBy      uint64         `json:"created_by"`
	UpdatedBy      uint64         `json:"updated_by"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Koperasi        Koperasi        `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Provinsi        WilayahProvinsi `gorm:"foreignKey:ProvinsiID" json:"provinsi,omitempty"`
	Kabupaten       WilayahKabupaten `gorm:"foreignKey:KabupatenID" json:"kabupaten,omitempty"`
	PurchaseOrder   []PurchaseOrder `gorm:"foreignKey:SupplierID" json:"purchase_order,omitempty"`
	SupplierProduk  []SupplierProduk `gorm:"foreignKey:SupplierID" json:"supplier_produk,omitempty"`
}

type Produk struct {
	ID                uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID        uint64         `gorm:"not null;index" json:"koperasi_id"`
	KategoriProdukID  uint64         `gorm:"not null" json:"kategori_produk_id"`
	SatuanProdukID    uint64         `gorm:"not null" json:"satuan_produk_id"`
	KodeProduk        string         `gorm:"size:50;not null" json:"kode_produk"`
	Barcode           string         `gorm:"size:50;uniqueIndex" json:"barcode"`
	NamaProduk        string         `gorm:"size:255;not null" json:"nama_produk"`
	Deskripsi         string         `gorm:"type:text" json:"deskripsi"`
	Brand             string         `gorm:"size:100" json:"brand"`
	Varian            string         `gorm:"size:100" json:"varian"`
	BeratBersih       float64        `gorm:"type:decimal(10,2)" json:"berat_bersih"`
	Dimensi           string         `gorm:"size:50" json:"dimensi"`
	FotoProduk        string         `gorm:"size:500" json:"foto_produk"`
	HargaBeli         float64        `gorm:"type:decimal(15,2);default:0" json:"harga_beli"`
	HargaJual         float64        `gorm:"type:decimal(15,2);not null" json:"harga_jual"`
	MarginPersen      float64        `gorm:"type:decimal(5,2);default:0" json:"margin_persen"`
	StokMinimal       int            `gorm:"default:0" json:"stok_minimal"`
	StokMaksimal      int            `gorm:"default:0" json:"stok_maksimal"`
	StokCurrent       int            `gorm:"default:0" json:"stok_current"`
	IsPerishable      bool           `gorm:"default:false;comment:produk mudah rusak" json:"is_perishable"`
	ShelfLife         int            `gorm:"default:0;comment:masa simpan dalam hari" json:"shelf_life"`
	IsProduksi        bool           `gorm:"default:false;comment:produk hasil produksi sendiri" json:"is_produksi"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	IsReadyStock      bool           `gorm:"default:true" json:"is_ready_stock"`
	TanggalExpired    *time.Time     `json:"tanggal_expired"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedBy         uint64         `json:"created_by"`
	UpdatedBy         uint64         `json:"updated_by"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Koperasi          Koperasi          `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	KategoriProduk    KategoriProduk    `gorm:"foreignKey:KategoriProdukID" json:"kategori_produk,omitempty"`
	SatuanProduk      SatuanProduk      `gorm:"foreignKey:SatuanProdukID" json:"satuan_produk,omitempty"`
	SupplierProduk    []SupplierProduk  `gorm:"foreignKey:ProdukID" json:"supplier_produk,omitempty"`
	PembelianDetail   []PembelianDetail `gorm:"foreignKey:ProdukID" json:"pembelian_detail,omitempty"`
	PenjualanDetail   []PenjualanDetail `gorm:"foreignKey:ProdukID" json:"penjualan_detail,omitempty"`
	StokMovement      []StokMovement    `gorm:"foreignKey:ProdukID" json:"stok_movement,omitempty"`
	ProdukDiskon      []ProdukDiskon    `gorm:"foreignKey:ProdukID" json:"produk_diskon,omitempty"`
}

type SupplierProduk struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	SupplierID    uint64         `gorm:"not null" json:"supplier_id"`
	ProdukID      uint64         `gorm:"not null" json:"produk_id"`
	KodeSupplier  string         `gorm:"size:50" json:"kode_supplier"`
	HargaSupplier float64        `gorm:"type:decimal(15,2)" json:"harga_supplier"`
	MinOrder      int            `gorm:"default:1" json:"min_order"`
	LeadTime      int            `gorm:"default:1;comment:hari" json:"lead_time"`
	IsPreferred   bool           `gorm:"default:false" json:"is_preferred"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Supplier Supplier `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Produk   Produk   `gorm:"foreignKey:ProdukID" json:"produk,omitempty"`
}

type PurchaseOrder struct {
	ID               uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID       uint64         `gorm:"not null;index" json:"koperasi_id"`
	SupplierID       uint64         `gorm:"not null" json:"supplier_id"`
	NomorPO          string         `gorm:"size:50;not null;uniqueIndex" json:"nomor_po"`
	TanggalPO        time.Time      `json:"tanggal_po"`
	TanggalKirim     *time.Time     `json:"tanggal_kirim"`
	TotalItem        int            `gorm:"default:0" json:"total_item"`
	SubTotal         float64        `gorm:"type:decimal(15,2);default:0" json:"sub_total"`
	PajakPersen      float64        `gorm:"type:decimal(5,2);default:0" json:"pajak_persen"`
	TotalPajak       float64        `gorm:"type:decimal(15,2);default:0" json:"total_pajak"`
	BiayaKirim       float64        `gorm:"type:decimal(15,2);default:0" json:"biaya_kirim"`
	Diskon           float64        `gorm:"type:decimal(15,2);default:0" json:"diskon"`
	GrandTotal       float64        `gorm:"type:decimal(15,2);default:0" json:"grand_total"`
	Status           string         `gorm:"type:varchar(20);default:'draft'" json:"status"`
	Keterangan       string         `gorm:"type:text" json:"keterangan"`
	ApprovedBy       uint64         `json:"approved_by"`
	ApprovedAt       *time.Time     `json:"approved_at"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedBy        uint64         `json:"created_by"`
	UpdatedBy        uint64         `json:"updated_by"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Koperasi           Koperasi           `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Supplier           Supplier           `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	PurchaseOrderDetail []PurchaseOrderDetail `gorm:"foreignKey:PurchaseOrderID" json:"purchase_order_detail,omitempty"`
}

type PurchaseOrderDetail struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	PurchaseOrderID uint64         `gorm:"not null" json:"purchase_order_id"`
	ProdukID        uint64         `gorm:"not null" json:"produk_id"`
	Qty             int            `gorm:"not null" json:"qty"`
	HargaSatuan     float64        `gorm:"type:decimal(15,2);not null" json:"harga_satuan"`
	Subtotal        float64        `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	QtyReceived     int            `gorm:"default:0" json:"qty_received"`
	Keterangan      string         `gorm:"type:text" json:"keterangan"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	PurchaseOrder PurchaseOrder `gorm:"foreignKey:PurchaseOrderID" json:"purchase_order,omitempty"`
	Produk        Produk        `gorm:"foreignKey:ProdukID" json:"produk,omitempty"`
}

type PembelianHeader struct {
	ID               uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID       uint64         `gorm:"not null;index" json:"koperasi_id"`
	SupplierID       uint64         `gorm:"not null" json:"supplier_id"`
	PurchaseOrderID  uint64         `json:"purchase_order_id"`
	NomorFaktur      string         `gorm:"size:50;not null;uniqueIndex" json:"nomor_faktur"`
	TanggalFaktur    time.Time      `json:"tanggal_faktur"`
	TanggalJatuhTempo *time.Time    `json:"tanggal_jatuh_tempo"`
	TotalItem        int            `gorm:"default:0" json:"total_item"`
	SubTotal         float64        `gorm:"type:decimal(15,2);default:0" json:"sub_total"`
	PajakPersen      float64        `gorm:"type:decimal(5,2);default:0" json:"pajak_persen"`
	TotalPajak       float64        `gorm:"type:decimal(15,2);default:0" json:"total_pajak"`
	BiayaKirim       float64        `gorm:"type:decimal(15,2);default:0" json:"biaya_kirim"`
	Diskon           float64        `gorm:"type:decimal(15,2);default:0" json:"diskon"`
	GrandTotal       float64        `gorm:"type:decimal(15,2);default:0" json:"grand_total"`
	StatusPembayaran string         `gorm:"type:varchar(20);default:'unpaid'" json:"status_pembayaran"`
	TotalBayar       float64        `gorm:"type:decimal(15,2);default:0" json:"total_bayar"`
	Keterangan       string         `gorm:"type:text" json:"keterangan"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedBy        uint64         `json:"created_by"`
	UpdatedBy        uint64         `json:"updated_by"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Koperasi        Koperasi        `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Supplier        Supplier        `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	PurchaseOrder   PurchaseOrder   `gorm:"foreignKey:PurchaseOrderID" json:"purchase_order,omitempty"`
	PembelianDetail []PembelianDetail `gorm:"foreignKey:PembelianHeaderID" json:"pembelian_detail,omitempty"`
	PembayaranPembelian []PembayaranPembelian `gorm:"foreignKey:PembelianHeaderID" json:"pembayaran_pembelian,omitempty"`
}

type PembelianDetail struct {
	ID                 uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	PembelianHeaderID  uint64         `gorm:"not null" json:"pembelian_header_id"`
	ProdukID           uint64         `gorm:"not null" json:"produk_id"`
	Qty                int            `gorm:"not null" json:"qty"`
	HargaSatuan        float64        `gorm:"type:decimal(15,2);not null" json:"harga_satuan"`
	Subtotal           float64        `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	TanggalExpired     *time.Time     `json:"tanggal_expired"`
	BatchNumber        string         `gorm:"size:50" json:"batch_number"`
	Keterangan         string         `gorm:"type:text" json:"keterangan"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	PembelianHeader PembelianHeader `gorm:"foreignKey:PembelianHeaderID" json:"pembelian_header,omitempty"`
	Produk          Produk          `gorm:"foreignKey:ProdukID" json:"produk,omitempty"`
}

type PembayaranPembelian struct {
	ID                uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	PembelianHeaderID uint64         `gorm:"not null" json:"pembelian_header_id"`
	TanggalBayar      time.Time      `json:"tanggal_bayar"`
	JumlahBayar       float64        `gorm:"type:decimal(15,2);not null" json:"jumlah_bayar"`
	MetodePembayaran  string         `gorm:"type:varchar(20);default:'cash'" json:"metode_pembayaran"`
	NomorReferensi    string         `gorm:"size:100" json:"nomor_referensi"`
	Keterangan        string         `gorm:"type:text" json:"keterangan"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy         uint64         `json:"created_by"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	PembelianHeader PembelianHeader `gorm:"foreignKey:PembelianHeaderID" json:"pembelian_header,omitempty"`
}

type PenjualanHeader struct {
	ID               uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID       uint64         `gorm:"not null;index" json:"koperasi_id"`
	AnggotaID        uint64         `json:"anggota_id"`
	NomorTransaksi   string         `gorm:"size:50;not null;uniqueIndex" json:"nomor_transaksi"`
	TanggalTransaksi time.Time      `json:"tanggal_transaksi"`
	TotalItem        int            `gorm:"default:0" json:"total_item"`
	SubTotal         float64        `gorm:"type:decimal(15,2);default:0" json:"sub_total"`
	PajakPersen      float64        `gorm:"type:decimal(5,2);default:0" json:"pajak_persen"`
	TotalPajak       float64        `gorm:"type:decimal(15,2);default:0" json:"total_pajak"`
	Diskon           float64        `gorm:"type:decimal(15,2);default:0" json:"diskon"`
	GrandTotal       float64        `gorm:"type:decimal(15,2);default:0" json:"grand_total"`
	MetodePembayaran string         `gorm:"type:varchar(20);default:'cash'" json:"metode_pembayaran"`
	StatusPembayaran string         `gorm:"type:varchar(20);default:'pending'" json:"status_pembayaran"`
	JumlahBayar      float64        `gorm:"type:decimal(15,2);default:0" json:"jumlah_bayar"`
	JumlahKembalian  float64        `gorm:"type:decimal(15,2);default:0" json:"jumlah_kembalian"`
	Kasir            string         `gorm:"size:100" json:"kasir"`
	Keterangan       string         `gorm:"type:text" json:"keterangan"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedBy        uint64         `json:"created_by"`
	UpdatedBy        uint64         `json:"updated_by"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Koperasi        Koperasi        `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Anggota         AnggotaKoperasi `gorm:"foreignKey:AnggotaID" json:"anggota,omitempty"`
	PenjualanDetail []PenjualanDetail `gorm:"foreignKey:PenjualanHeaderID" json:"penjualan_detail,omitempty"`
}

type PenjualanDetail struct {
	ID                uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	PenjualanHeaderID uint64         `gorm:"not null" json:"penjualan_header_id"`
	ProdukID          uint64         `gorm:"not null" json:"produk_id"`
	Qty               int            `gorm:"not null" json:"qty"`
	HargaSatuan       float64        `gorm:"type:decimal(15,2);not null" json:"harga_satuan"`
	DiskonPersen      float64        `gorm:"type:decimal(5,2);default:0" json:"diskon_persen"`
	DiskonRupiah      float64        `gorm:"type:decimal(15,2);default:0" json:"diskon_rupiah"`
	Subtotal          float64        `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	Keterangan        string         `gorm:"type:text" json:"keterangan"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	PenjualanHeader PenjualanHeader `gorm:"foreignKey:PenjualanHeaderID" json:"penjualan_header,omitempty"`
	Produk          Produk          `gorm:"foreignKey:ProdukID" json:"produk,omitempty"`
}

type StokMovement struct {
	ID             uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID     uint64         `gorm:"not null;index" json:"koperasi_id"`
	ProdukID       uint64         `gorm:"not null" json:"produk_id"`
	TipeMovement   string         `gorm:"type:varchar(20);not null" json:"tipe_movement"`
	ReferensiTipe  string         `gorm:"type:varchar(20);not null" json:"referensi_tipe"`
	ReferensiID    uint64         `json:"referensi_id"`
	TanggalMovement time.Time     `json:"tanggal_movement"`
	QtyBefore      int            `gorm:"not null" json:"qty_before"`
	QtyMovement    int            `gorm:"not null" json:"qty_movement"`
	QtyAfter       int            `gorm:"not null" json:"qty_after"`
	HargaSatuan    float64        `gorm:"type:decimal(15,2)" json:"harga_satuan"`
	TotalNilai     float64        `gorm:"type:decimal(15,2)" json:"total_nilai"`
	Keterangan     string         `gorm:"type:text" json:"keterangan"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy      uint64         `json:"created_by"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Koperasi Koperasi `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Produk   Produk   `gorm:"foreignKey:ProdukID" json:"produk,omitempty"`
}

type ProdukDiskon struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProdukID      uint64         `gorm:"not null" json:"produk_id"`
	NamaDiskon    string         `gorm:"size:100;not null" json:"nama_diskon"`
	TipeDiskon    string         `gorm:"type:varchar(20);default:'percentage'" json:"tipe_diskon"`
	NilaiDiskon   float64        `gorm:"type:decimal(15,2);not null" json:"nilai_diskon"`
	TanggalMulai  time.Time      `json:"tanggal_mulai"`
	TanggalSelesai time.Time     `json:"tanggal_selesai"`
	MinimumBeli   int            `gorm:"default:1" json:"minimum_beli"`
	MaksimumDiskon float64       `gorm:"type:decimal(15,2);default:0" json:"maksimum_diskon"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedBy     uint64         `json:"created_by"`
	UpdatedBy     uint64         `json:"updated_by"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Produk Produk `gorm:"foreignKey:ProdukID" json:"produk,omitempty"`
}