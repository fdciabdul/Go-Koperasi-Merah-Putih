package services

import (
	"fmt"
	"time"

	"koperasi-merah-putih/internal/models/postgres"
	postgresRepo "koperasi-merah-putih/internal/repository/postgres"
)

type KlinikService struct {
	klinikRepo      *postgresRepo.KlinikRepository
	sequenceService *SequenceService
}

func NewKlinikService(
	klinikRepo *postgresRepo.KlinikRepository,
	sequenceService *SequenceService,
) *KlinikService {
	return &KlinikService{
		klinikRepo:      klinikRepo,
		sequenceService: sequenceService,
	}
}

func (s *KlinikService) CreatePasien(req *CreatePasienRequest) (*postgres.KlinikPasien, error) {
	existing, _ := s.klinikRepo.GetPasienByNomorRM(req.NomorRM)
	if existing != nil {
		return nil, fmt.Errorf("pasien with nomor RM %s already exists", req.NomorRM)
	}

	if req.NomorRM == "" {
		nomorRM, err := s.generateNomorRM(req.KoperasiID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate nomor RM: %v", err)
		}
		req.NomorRM = nomorRM
	}

	pasien := &postgres.KlinikPasien{
		KoperasiID:      req.KoperasiID,
		NomorRM:         req.NomorRM,
		NIK:             req.NIK,
		NamaLengkap:     req.NamaLengkap,
		JenisKelamin:    req.JenisKelamin,
		TempatLahir:     req.TempatLahir,
		TanggalLahir:    req.TanggalLahir,
		Alamat:          req.Alamat,
		Telepon:         req.Telepon,
		Email:           req.Email,
		GolonganDarah:   req.GolonganDarah,
		Alergi:          req.Alergi,
		RiwayatPenyakit: req.RiwayatPenyakit,
		AnggotaID:       req.AnggotaID,
	}

	err := s.klinikRepo.CreatePasien(pasien)
	if err != nil {
		return nil, fmt.Errorf("failed to create pasien: %v", err)
	}

	return pasien, nil
}

func (s *KlinikService) GetPasienByID(id uint64) (*postgres.KlinikPasien, error) {
	return s.klinikRepo.GetPasienByID(id)
}

func (s *KlinikService) GetPasienList(koperasiID uint64, page, limit int) ([]postgres.KlinikPasien, error) {
	offset := (page - 1) * limit
	return s.klinikRepo.GetPasienByKoperasi(koperasiID, limit, offset)
}

func (s *KlinikService) SearchPasien(koperasiID uint64, search string) ([]postgres.KlinikPasien, error) {
	return s.klinikRepo.SearchPasien(koperasiID, search)
}

func (s *KlinikService) CreateTenagaMedis(req *CreateTenagaMedisRequest) (*postgres.KlinikTenagaMedis, error) {
	tenagaMedis := &postgres.KlinikTenagaMedis{
		KoperasiID:       req.KoperasiID,
		NIP:              req.NIP,
		NamaLengkap:      req.NamaLengkap,
		JenisKelamin:     req.JenisKelamin,
		Spesialisasi:     req.Spesialisasi,
		NoSTR:            req.NoSTR,
		NoSIP:            req.NoSIP,
		Telepon:          req.Telepon,
		Email:            req.Email,
		JadwalPraktik:    req.JadwalPraktik,
		TarifKonsultasi:  req.TarifKonsultasi,
		Status:           "aktif",
	}

	err := s.klinikRepo.CreateTenagaMedis(tenagaMedis)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenaga medis: %v", err)
	}

	return tenagaMedis, nil
}

func (s *KlinikService) GetTenagaMedisList(koperasiID uint64) ([]postgres.KlinikTenagaMedis, error) {
	return s.klinikRepo.GetTenagaMedisByKoperasi(koperasiID)
}

func (s *KlinikService) CreateKunjungan(req *CreateKunjunganRequest) (*postgres.KlinikKunjungan, error) {
	nomorKunjungan, err := s.generateNomorKunjungan(req.KoperasiID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nomor kunjungan: %v", err)
	}

	kunjungan := &postgres.KlinikKunjungan{
		KoperasiID:       req.KoperasiID,
		PasienID:         req.PasienID,
		DokterID:         req.DokterID,
		NomorKunjungan:   nomorKunjungan,
		TanggalKunjungan: time.Now(),
		KeluhanUtama:     req.KeluhanUtama,
		Anamnesis:        req.Anamnesis,
		PemeriksaanFisik: req.PemeriksaanFisik,
		Diagnosis:        req.Diagnosis,
		TerapiPengobatan: req.TerapiPengobatan,
		BiayaKonsultasi:  req.BiayaKonsultasi,
		BiayaTindakan:    req.BiayaTindakan,
		StatusPembayaran: "belum_bayar",
	}

	err = s.klinikRepo.CreateKunjungan(kunjungan)
	if err != nil {
		return nil, fmt.Errorf("failed to create kunjungan: %v", err)
	}

	if len(req.Reseps) > 0 {
		err = s.addResepToKunjungan(kunjungan.ID, req.Reseps)
		if err != nil {
			return nil, fmt.Errorf("failed to add resep: %v", err)
		}

		totalBiayaObat, err := s.calculateTotalBiayaObat(req.Reseps)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate biaya obat: %v", err)
		}

		kunjungan.BiayaObat = totalBiayaObat
		kunjungan.TotalBiaya = kunjungan.BiayaKonsultasi + kunjungan.BiayaTindakan + kunjungan.BiayaObat

		err = s.klinikRepo.UpdateKunjungan(kunjungan)
		if err != nil {
			return nil, fmt.Errorf("failed to update total biaya: %v", err)
		}
	}

	return kunjungan, nil
}

func (s *KlinikService) GetKunjunganByID(id uint64) (*postgres.KlinikKunjungan, error) {
	return s.klinikRepo.GetKunjunganByID(id)
}

func (s *KlinikService) GetKunjunganByPasien(pasienID uint64, page, limit int) ([]postgres.KlinikKunjungan, error) {
	offset := (page - 1) * limit
	return s.klinikRepo.GetKunjunganByPasien(pasienID, limit, offset)
}

func (s *KlinikService) CreateObat(req *CreateObatRequest) (*postgres.KlinikObat, error) {
	obat := &postgres.KlinikObat{
		KoperasiID:    req.KoperasiID,
		KodeObat:      req.KodeObat,
		NamaObat:      req.NamaObat,
		Kategori:      req.Kategori,
		BentukSediaan: req.BentukSediaan,
		Kekuatan:      req.Kekuatan,
		Satuan:        req.Satuan,
		StokMinimal:   req.StokMinimal,
		StokCurrent:   req.StokCurrent,
		HargaBeli:     req.HargaBeli,
		HargaJual:     req.HargaJual,
		IsAktif:       true,
	}

	err := s.klinikRepo.CreateObat(obat)
	if err != nil {
		return nil, fmt.Errorf("failed to create obat: %v", err)
	}

	return obat, nil
}

func (s *KlinikService) GetObatList(koperasiID uint64) ([]postgres.KlinikObat, error) {
	return s.klinikRepo.GetObatByKoperasi(koperasiID)
}

func (s *KlinikService) SearchObat(koperasiID uint64, search string) ([]postgres.KlinikObat, error) {
	return s.klinikRepo.SearchObat(koperasiID, search)
}

func (s *KlinikService) GetStatistik(koperasiID uint64) (*postgresRepo.KlinikStatistik, error) {
	return s.klinikRepo.GetStatistikKlinik(koperasiID)
}

func (s *KlinikService) GetObatStokRendah(koperasiID uint64) ([]postgres.KlinikObat, error) {
	return s.klinikRepo.GetObatStokRendah(koperasiID)
}

func (s *KlinikService) addResepToKunjungan(kunjunganID uint64, resepReqs []ResepRequest) error {
	var reseps []postgres.KlinikResep

	for _, req := range resepReqs {
		obat, err := s.klinikRepo.GetObatByID(req.ObatID)
		if err != nil {
			return fmt.Errorf("obat not found: %v", err)
		}

		if obat.StokCurrent < req.Jumlah {
			return fmt.Errorf("insufficient stock for %s", obat.NamaObat)
		}

		totalHarga := obat.HargaJual * float64(req.Jumlah)

		resep := postgres.KlinikResep{
			KunjunganID: kunjunganID,
			ObatID:      req.ObatID,
			Jumlah:      req.Jumlah,
			AturanPakai: req.AturanPakai,
			Keterangan:  req.Keterangan,
			HargaSatuan: obat.HargaJual,
			TotalHarga:  totalHarga,
		}

		reseps = append(reseps, resep)

		err = s.klinikRepo.UpdateStokObat(req.ObatID, req.Jumlah)
		if err != nil {
			return fmt.Errorf("failed to update stock: %v", err)
		}
	}

	return s.klinikRepo.CreateResep(reseps)
}

func (s *KlinikService) calculateTotalBiayaObat(resepReqs []ResepRequest) (float64, error) {
	var total float64

	for _, req := range resepReqs {
		obat, err := s.klinikRepo.GetObatByID(req.ObatID)
		if err != nil {
			return 0, err
		}
		total += obat.HargaJual * float64(req.Jumlah)
	}

	return total, nil
}

func (s *KlinikService) generateNomorRM(koperasiID uint64) (string, error) {
	number, err := s.sequenceService.GetNextNumber(1, koperasiID, "nomor_rm")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("RM%04d%06d", koperasiID, number), nil
}

func (s *KlinikService) generateNomorKunjungan(koperasiID uint64) (string, error) {
	number, err := s.sequenceService.GetNextNumber(1, koperasiID, "kunjungan")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("KUN%04d%08d", koperasiID, number), nil
}

type CreatePasienRequest struct {
	KoperasiID      uint64     `json:"koperasi_id" binding:"required"`
	NomorRM         string     `json:"nomor_rm"`
	NIK             string     `json:"nik"`
	NamaLengkap     string     `json:"nama_lengkap" binding:"required"`
	JenisKelamin    string     `json:"jenis_kelamin" binding:"required,oneof=L P"`
	TempatLahir     string     `json:"tempat_lahir"`
	TanggalLahir    *time.Time `json:"tanggal_lahir"`
	Alamat          string     `json:"alamat"`
	Telepon         string     `json:"telepon"`
	Email           string     `json:"email"`
	GolonganDarah   string     `json:"golongan_darah"`
	Alergi          string     `json:"alergi"`
	RiwayatPenyakit string     `json:"riwayat_penyakit"`
	AnggotaID       uint64     `json:"anggota_id"`
}

type CreateTenagaMedisRequest struct {
	KoperasiID      uint64  `json:"koperasi_id" binding:"required"`
	NIP             string  `json:"nip"`
	NamaLengkap     string  `json:"nama_lengkap" binding:"required"`
	JenisKelamin    string  `json:"jenis_kelamin" binding:"required,oneof=L P"`
	Spesialisasi    string  `json:"spesialisasi"`
	NoSTR           string  `json:"no_str"`
	NoSIP           string  `json:"no_sip"`
	Telepon         string  `json:"telepon"`
	Email           string  `json:"email"`
	JadwalPraktik   string  `json:"jadwal_praktik"`
	TarifKonsultasi float64 `json:"tarif_konsultasi"`
}

type CreateKunjunganRequest struct {
	KoperasiID       uint64         `json:"koperasi_id" binding:"required"`
	PasienID         uint64         `json:"pasien_id" binding:"required"`
	DokterID         uint64         `json:"dokter_id" binding:"required"`
	KeluhanUtama     string         `json:"keluhan_utama"`
	Anamnesis        string         `json:"anamnesis"`
	PemeriksaanFisik string         `json:"pemeriksaan_fisik"`
	Diagnosis        string         `json:"diagnosis"`
	TerapiPengobatan string         `json:"terapi_pengobatan"`
	BiayaKonsultasi  float64        `json:"biaya_konsultasi"`
	BiayaTindakan    float64        `json:"biaya_tindakan"`
	Reseps           []ResepRequest `json:"reseps"`
}

type ResepRequest struct {
	ObatID      uint64 `json:"obat_id" binding:"required"`
	Jumlah      int    `json:"jumlah" binding:"required,gt=0"`
	AturanPakai string `json:"aturan_pakai"`
	Keterangan  string `json:"keterangan"`
}

type CreateObatRequest struct {
	KoperasiID    uint64  `json:"koperasi_id" binding:"required"`
	KodeObat      string  `json:"kode_obat" binding:"required"`
	NamaObat      string  `json:"nama_obat" binding:"required"`
	Kategori      string  `json:"kategori"`
	BentukSediaan string  `json:"bentuk_sediaan"`
	Kekuatan      string  `json:"kekuatan"`
	Satuan        string  `json:"satuan"`
	StokMinimal   int     `json:"stok_minimal"`
	StokCurrent   int     `json:"stok_current"`
	HargaBeli     float64 `json:"harga_beli"`
	HargaJual     float64 `json:"harga_jual"`
}