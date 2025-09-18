package services

import (
	"fmt"
	"time"

	"koperasi-merah-putih/internal/models/postgres"
	postgresRepo "koperasi-merah-putih/internal/repository/postgres"
)

type KoperasiService struct {
	koperasiRepo *postgresRepo.KoperasiRepository
	anggotaRepo  *postgresRepo.AnggotaKoperasiRepository
	wilayahRepo  *postgresRepo.WilayahRepository
	sequenceService *SequenceService
}

func NewKoperasiService(
	koperasiRepo *postgresRepo.KoperasiRepository,
	anggotaRepo *postgresRepo.AnggotaKoperasiRepository,
	wilayahRepo *postgresRepo.WilayahRepository,
	sequenceService *SequenceService,
) *KoperasiService {
	return &KoperasiService{
		koperasiRepo: koperasiRepo,
		anggotaRepo:  anggotaRepo,
		wilayahRepo:  wilayahRepo,
		sequenceService: sequenceService,
	}
}

func (s *KoperasiService) CreateKoperasi(req *CreateKoperasiRequest) (*postgres.Koperasi, error) {
	existing, _ := s.koperasiRepo.GetByNIK(req.NIK)
	if existing != nil {
		return nil, fmt.Errorf("koperasi with NIK %d already exists", req.NIK)
	}

	existing, _ = s.koperasiRepo.GetByNomorSK(req.NomorSK)
	if existing != nil {
		return nil, fmt.Errorf("koperasi with nomor SK %s already exists", req.NomorSK)
	}

	koperasi := &postgres.Koperasi{
		TenantID:             req.TenantID,
		NomorSK:              req.NomorSK,
		NIK:                  req.NIK,
		NamaKoperasi:         req.NamaKoperasi,
		NamaSK:               req.NamaSK,
		JenisKoperasiID:      req.JenisKoperasiID,
		BentukKoperasiID:     req.BentukKoperasiID,
		StatusKoperasiID:     req.StatusKoperasiID,
		ProvinsiID:           req.ProvinsiID,
		KabupatenID:          req.KabupatenID,
		KecamatanID:          req.KecamatanID,
		KelurahanID:          req.KelurahanID,
		Alamat:               req.Alamat,
		RT:                   req.RT,
		RW:                   req.RW,
		KodePos:              req.KodePos,
		Email:                req.Email,
		Telepon:              req.Telepon,
		Website:              req.Website,
		TanggalBerdiri:       req.TanggalBerdiri,
		TanggalSK:            req.TanggalSK,
		TanggalPengesahan:    req.TanggalPengesahan,
		CreatedBy:            req.CreatedBy,
	}

	err := s.koperasiRepo.Create(koperasi)
	if err != nil {
		return nil, fmt.Errorf("failed to create koperasi: %v", err)
	}

	return koperasi, nil
}

func (s *KoperasiService) GetKoperasiByID(id uint64) (*postgres.Koperasi, error) {
	return s.koperasiRepo.GetByID(id)
}

func (s *KoperasiService) GetKoperasiByTenant(tenantID uint64) ([]postgres.Koperasi, error) {
	return s.koperasiRepo.GetByTenantID(tenantID)
}

func (s *KoperasiService) UpdateKoperasi(id uint64, req *UpdateKoperasiRequest) (*postgres.Koperasi, error) {
	koperasi, err := s.koperasiRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("koperasi not found: %v", err)
	}

	if req.NamaKoperasi != "" {
		koperasi.NamaKoperasi = req.NamaKoperasi
	}
	if req.NamaSK != "" {
		koperasi.NamaSK = req.NamaSK
	}
	if req.JenisKoperasiID != 0 {
		koperasi.JenisKoperasiID = req.JenisKoperasiID
	}
	if req.BentukKoperasiID != 0 {
		koperasi.BentukKoperasiID = req.BentukKoperasiID
	}
	if req.StatusKoperasiID != 0 {
		koperasi.StatusKoperasiID = req.StatusKoperasiID
	}
	if req.Alamat != "" {
		koperasi.Alamat = req.Alamat
	}
	if req.Email != "" {
		koperasi.Email = req.Email
	}
	if req.Telepon != "" {
		koperasi.Telepon = req.Telepon
	}
	if req.Website != "" {
		koperasi.Website = req.Website
	}

	koperasi.UpdatedBy = req.UpdatedBy

	err = s.koperasiRepo.Update(koperasi)
	if err != nil {
		return nil, fmt.Errorf("failed to update koperasi: %v", err)
	}

	return koperasi, nil
}

func (s *KoperasiService) DeleteKoperasi(id uint64) error {
	count, err := s.anggotaRepo.CountByKoperasiID(id)
	if err != nil {
		return fmt.Errorf("failed to check anggota count: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete koperasi with existing members")
	}

	return s.koperasiRepo.Delete(id)
}

func (s *KoperasiService) CreateAnggota(req *CreateAnggotaRequest) (*postgres.AnggotaKoperasi, error) {
	existing, _ := s.anggotaRepo.GetByNIAK(req.NIAK)
	if existing != nil {
		return nil, fmt.Errorf("anggota with NIAK %s already exists", req.NIAK)
	}

	if req.NIAK == "" {
		niak, err := s.generateNIAK(req.KoperasiID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate NIAK: %v", err)
		}
		req.NIAK = niak
	}

	anggota := &postgres.AnggotaKoperasi{
		KoperasiID:    req.KoperasiID,
		NIAK:          req.NIAK,
		NIK:           req.NIK,
		Nama:          req.Nama,
		JenisKelamin:  req.JenisKelamin,
		TempatLahir:   req.TempatLahir,
		TanggalLahir:  req.TanggalLahir,
		Alamat:        req.Alamat,
		RT:            req.RT,
		RW:            req.RW,
		KelurahanID:   req.KelurahanID,
		Telepon:       req.Telepon,
		Email:         req.Email,
		Posisi:        req.Posisi,
		JabatanID:     req.JabatanID,
		TanggalMasuk:  req.TanggalMasuk,
		StatusAnggota: "aktif",
		NPWP:          req.NPWP,
		Pekerjaan:     req.Pekerjaan,
		Pendidikan:    req.Pendidikan,
	}

	err := s.anggotaRepo.Create(anggota)
	if err != nil {
		return nil, fmt.Errorf("failed to create anggota: %v", err)
	}

	return anggota, nil
}

func (s *KoperasiService) GetAnggotaByID(id uint64) (*postgres.AnggotaKoperasi, error) {
	return s.anggotaRepo.GetByID(id)
}

func (s *KoperasiService) GetAnggotaByKoperasi(koperasiID uint64, page, limit int) ([]postgres.AnggotaKoperasi, error) {
	offset := (page - 1) * limit
	return s.anggotaRepo.GetByKoperasiID(koperasiID, limit, offset)
}

func (s *KoperasiService) UpdateAnggotaStatus(id uint64, status string) error {
	return s.anggotaRepo.UpdateStatus(id, status)
}

func (s *KoperasiService) GetProvinsiList() ([]postgres.WilayahProvinsi, error) {
	return s.wilayahRepo.GetProvinsiList()
}

func (s *KoperasiService) GetKabupatenByProvinsi(provinsiID uint64) ([]postgres.WilayahKabupaten, error) {
	return s.wilayahRepo.GetKabupatenByProvinsiID(provinsiID)
}

func (s *KoperasiService) GetKecamatanByKabupaten(kabupatenID uint64) ([]postgres.WilayahKecamatan, error) {
	return s.wilayahRepo.GetKecamatanByKabupatenID(kabupatenID)
}

func (s *KoperasiService) GetKelurahanByKecamatan(kecamatanID uint64) ([]postgres.WilayahKelurahan, error) {
	return s.wilayahRepo.GetKelurahanByKecamatanID(kecamatanID)
}

func (s *KoperasiService) generateNIAK(koperasiID uint64) (string, error) {
	number, err := s.sequenceService.GetNextNumber(1, koperasiID, "anggota")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("ANG%04d%06d", koperasiID, number), nil
}

type CreateKoperasiRequest struct {
	TenantID             uint64     `json:"tenant_id" binding:"required"`
	NomorSK              string     `json:"nomor_sk" binding:"required"`
	NIK                  uint64     `json:"nik" binding:"required"`
	NamaKoperasi         string     `json:"nama_koperasi" binding:"required"`
	NamaSK               string     `json:"nama_sk" binding:"required"`
	JenisKoperasiID      uint64     `json:"jenis_koperasi_id"`
	BentukKoperasiID     uint64     `json:"bentuk_koperasi_id"`
	StatusKoperasiID     uint64     `json:"status_koperasi_id"`
	ProvinsiID           uint64     `json:"provinsi_id"`
	KabupatenID          uint64     `json:"kabupaten_id"`
	KecamatanID          uint64     `json:"kecamatan_id"`
	KelurahanID          uint64     `json:"kelurahan_id"`
	Alamat               string     `json:"alamat"`
	RT                   string     `json:"rt"`
	RW                   string     `json:"rw"`
	KodePos              string     `json:"kode_pos"`
	Email                string     `json:"email"`
	Telepon              string     `json:"telepon"`
	Website              string     `json:"website"`
	TanggalBerdiri       *time.Time `json:"tanggal_berdiri"`
	TanggalSK            *time.Time `json:"tanggal_sk"`
	TanggalPengesahan    *time.Time `json:"tanggal_pengesahan"`
	CreatedBy            uint64     `json:"created_by"`
}

type UpdateKoperasiRequest struct {
	NamaKoperasi         string `json:"nama_koperasi"`
	NamaSK               string `json:"nama_sk"`
	JenisKoperasiID      uint64 `json:"jenis_koperasi_id"`
	BentukKoperasiID     uint64 `json:"bentuk_koperasi_id"`
	StatusKoperasiID     uint64 `json:"status_koperasi_id"`
	Alamat               string `json:"alamat"`
	Email                string `json:"email"`
	Telepon              string `json:"telepon"`
	Website              string `json:"website"`
	UpdatedBy            uint64 `json:"updated_by"`
}

type CreateAnggotaRequest struct {
	KoperasiID    uint64     `json:"koperasi_id" binding:"required"`
	NIAK          string     `json:"niak"`
	NIK           string     `json:"nik"`
	Nama          string     `json:"nama" binding:"required"`
	JenisKelamin  string     `json:"jenis_kelamin" binding:"required,oneof=L P"`
	TempatLahir   string     `json:"tempat_lahir"`
	TanggalLahir  *time.Time `json:"tanggal_lahir"`
	Alamat        string     `json:"alamat"`
	RT            string     `json:"rt"`
	RW            string     `json:"rw"`
	KelurahanID   uint64     `json:"kelurahan_id"`
	Telepon       string     `json:"telepon"`
	Email         string     `json:"email"`
	Posisi        string     `json:"posisi" binding:"oneof=pengurus pengawas anggota"`
	JabatanID     uint64     `json:"jabatan_id"`
	TanggalMasuk  *time.Time `json:"tanggal_masuk"`
	NPWP          string     `json:"npwp"`
	Pekerjaan     string     `json:"pekerjaan"`
	Pendidikan    string     `json:"pendidikan"`
}