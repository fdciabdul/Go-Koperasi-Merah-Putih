# Services Layer - Business Logic Implementation

## Purpose
Service layer adalah inti dari business logic aplikasi. Layer ini mengimplementasikan semua domain rules, validations, calculations, dan orchestrations yang diperlukan untuk operasional koperasi.

## Structure Pattern
```
services/
├── user_service.go           # User management dan authentication
├── koperasi_service.go       # Koperasi business operations
├── financial_service.go      # Financial calculations dan reporting
├── simpan_pinjam_service.go  # Savings/loans business logic
├── klinik_service.go         # Healthcare services logic
└── sequence_service.go       # Auto-numbering sequences
```

## Service Architecture

### Service Structure
```go
type EntityService struct {
    entityRepo      *repository.EntityRepository
    otherService    *OtherService
    // External dependencies
    paymentGateway  PaymentGatewayInterface
}

func NewEntityService(
    entityRepo *repository.EntityRepository,
    otherService *OtherService,
) *EntityService {
    return &EntityService{
        entityRepo:   entityRepo,
        otherService: otherService,
    }
}
```

### Business Method Pattern
```go
func (s *EntityService) CreateEntity(req *CreateEntityRequest) (*models.Entity, error) {
    // 1. Input validation
    if err := s.validateCreateRequest(req); err != nil {
        return nil, fmt.Errorf("validation failed: %v", err)
    }

    // 2. Business rules enforcement
    if err := s.checkBusinessRules(req); err != nil {
        return nil, fmt.Errorf("business rule violation: %v", err)
    }

    // 3. Data transformation
    entity := s.transformRequestToModel(req)

    // 4. External service calls (if needed)
    if err := s.notifyExternalService(entity); err != nil {
        return nil, fmt.Errorf("external service error: %v", err)
    }

    // 5. Repository call
    if err := s.entityRepo.CreateEntity(entity); err != nil {
        return nil, fmt.Errorf("failed to create entity: %v", err)
    }

    // 6. Post-processing
    s.performPostCreationTasks(entity)

    return entity, nil
}
```

## Request/Response Patterns

### Request Structs
```go
type CreateKoperasiRequest struct {
    TenantID         uint64 `json:"tenant_id" binding:"required"`
    Nama             string `json:"nama" binding:"required,min=3,max=255"`
    Email            string `json:"email" binding:"required,email"`
    Telepon          string `json:"telepon" binding:"required"`
    Alamat           string `json:"alamat" binding:"required"`
    ProvinsiID       uint64 `json:"provinsi_id" binding:"required"`
    KabupatenID      uint64 `json:"kabupaten_id" binding:"required"`
    JenisKoperasiID  uint64 `json:"jenis_koperasi_id" binding:"required"`
    BentukKoperasiID uint64 `json:"bentuk_koperasi_id" binding:"required"`
    KBLIID           uint64 `json:"kbli_id" binding:"required"`
}

type UpdateKoperasiRequest struct {
    Nama      string `json:"nama" binding:"required,min=3,max=255"`
    Email     string `json:"email" binding:"required,email"`
    Telepon   string `json:"telepon" binding:"required"`
    // Only updatable fields
}
```

## Indonesian Business Logic

### NIAK Generation
```go
func (s *KoperasiService) generateNIAK(provinsiID uint64, jenisKoperasiID uint64) (string, error) {
    // Format: PPJJNNNNNNNNNN
    // PP = Provinsi code (2 digits)
    // JJ = Jenis koperasi code (2 digits)
    // NNNNNNNNNN = Sequential number (10 digits)

    provinsi, err := s.wilayahRepo.GetProvinsiByID(provinsiID)
    if err != nil {
        return "", err
    }

    sequence, err := s.sequenceService.GetNextNumber(1, 0, "niak")
    if err != nil {
        return "", err
    }

    niak := fmt.Sprintf("%02s%02d%010d",
        provinsi.Kode,
        jenisKoperasiID,
        sequence)

    return niak, nil
}
```

### NIK Validation
```go
func (s *UserService) validateNIK(nik string) error {
    if len(nik) != 16 {
        return fmt.Errorf("NIK harus 16 digit")
    }

    if !isNumeric(nik) {
        return fmt.Errorf("NIK hanya boleh angka")
    }

    // Additional validation logic for Indonesian NIK format
    provinsiCode := nik[:2]
    if !s.isValidProvinsiCode(provinsiCode) {
        return fmt.Errorf("kode provinsi tidak valid")
    }

    return nil
}
```

## Financial Business Logic

### Interest Calculation
```go
func (s *SimpanPinjamService) calculateInterest(principal float64, rate float64, months int) float64 {
    // Simple interest calculation for Indonesian koperasi
    return (principal * rate * float64(months)) / (12 * 100)
}

func (s *SimpanPinjamService) calculateAngsuran(principal float64, rate float64, months int) float64 {
    // Calculate monthly installment (angsuran)
    monthlyRate := rate / (12 * 100)

    if monthlyRate == 0 {
        return principal / float64(months)
    }

    return principal * (monthlyRate * math.Pow(1+monthlyRate, float64(months))) /
           (math.Pow(1+monthlyRate, float64(months)) - 1)
}
```

### Journal Validation
```go
func (s *FinancialService) validateJournalEntry(req *CreateJurnalRequest) error {
    if len(req.Details) < 2 {
        return fmt.Errorf("jurnal harus memiliki minimal 2 detail")
    }

    var totalDebit, totalKredit float64
    for _, detail := range req.Details {
        if detail.Debit < 0 || detail.Kredit < 0 {
            return fmt.Errorf("nilai debit/kredit tidak boleh negatif")
        }

        if detail.Debit > 0 && detail.Kredit > 0 {
            return fmt.Errorf("detail tidak boleh memiliki debit dan kredit bersamaan")
        }

        totalDebit += detail.Debit
        totalKredit += detail.Kredit
    }

    if totalDebit != totalKredit {
        return fmt.Errorf("total debit (%.2f) harus sama dengan total kredit (%.2f)",
            totalDebit, totalKredit)
    }

    return nil
}
```

## Healthcare Business Logic

### Medical Record Number Generation
```go
func (s *KlinikService) generateNomorRM(koperasiID uint64) (string, error) {
    sequence, err := s.sequenceService.GetNextNumber(1, koperasiID, "nomor_rm")
    if err != nil {
        return "", err
    }

    // Format: RM + KoperasiID(4 digits) + Sequence(6 digits)
    return fmt.Sprintf("RM%04d%06d", koperasiID, sequence), nil
}
```

### Prescription Validation
```go
func (s *KlinikService) validateResep(reseps []ResepRequest) error {
    for _, resep := range reseps {
        // Check medicine availability
        obat, err := s.klinikRepo.GetObatByID(resep.ObatID)
        if err != nil {
            return fmt.Errorf("obat ID %d tidak ditemukan", resep.ObatID)
        }

        if !obat.IsAktif {
            return fmt.Errorf("obat %s tidak aktif", obat.NamaObat)
        }

        if obat.StokCurrent < resep.Jumlah {
            return fmt.Errorf("stok obat %s tidak mencukupi", obat.NamaObat)
        }

        if resep.Jumlah <= 0 {
            return fmt.Errorf("jumlah obat harus lebih dari 0")
        }
    }

    return nil
}
```

## Error Handling Patterns

### Business Error Types
```go
type BusinessError struct {
    Code    string
    Message string
}

func (e BusinessError) Error() string {
    return e.Message
}

var (
    ErrInsufficientBalance = BusinessError{"INSUFFICIENT_BALANCE", "Saldo tidak mencukupi"}
    ErrDuplicateEntry     = BusinessError{"DUPLICATE_ENTRY", "Data sudah ada"}
    ErrInvalidStatus      = BusinessError{"INVALID_STATUS", "Status tidak valid"}
)
```

### Error Wrapping
```go
func (s *Service) SomeOperation(req *Request) error {
    result, err := s.repository.GetData(req.ID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return fmt.Errorf("data dengan ID %d tidak ditemukan", req.ID)
        }
        return fmt.Errorf("failed to get data: %v", err)
    }

    // Business logic...
    return nil
}
```

## Transaction Management

### Database Transactions
```go
func (s *FinancialService) ProcessComplexTransaction(req *ComplexRequest) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Multiple operations that must all succeed
        if err := s.createJournal(tx, req.Journal); err != nil {
            return err
        }

        if err := s.updateBalances(tx, req.Balances); err != nil {
            return err
        }

        if err := s.notifyExternalSystem(req); err != nil {
            return err
        }

        return nil
    })
}
```

## Validation Patterns

### Input Validation
```go
func (s *Service) validateCreateRequest(req *CreateRequest) error {
    if req.Name == "" {
        return fmt.Errorf("nama tidak boleh kosong")
    }

    if len(req.Name) > 255 {
        return fmt.Errorf("nama maksimal 255 karakter")
    }

    if req.Email != "" && !isValidEmail(req.Email) {
        return fmt.Errorf("format email tidak valid")
    }

    return nil
}
```

### Business Rules Validation
```go
func (s *KoperasiService) validateKoperasiRules(req *CreateKoperasiRequest) error {
    // Check if NIAK already exists
    existing, _ := s.koperasiRepo.GetByNIAK(req.NIAK)
    if existing != nil {
        return fmt.Errorf("NIAK %s sudah digunakan", req.NIAK)
    }

    // Validate minimum members for koperasi type
    if req.JenisKoperasiID == PRIMER_TYPE && req.MinMembers < 20 {
        return fmt.Errorf("koperasi primer minimal 20 anggota")
    }

    return nil
}
```

## Calculation Helpers

### Financial Calculations
```go
func (s *FinancialService) calculateCompoundInterest(
    principal, rate float64,
    compoundingFreq, years int,
) float64 {
    return principal * math.Pow(
        1+(rate/(float64(compoundingFreq)*100)),
        float64(compoundingFreq*years),
    )
}

func (s *FinancialService) calculateTax(amount float64, taxRate float64) float64 {
    return amount * (taxRate / 100)
}
```

## External Service Integration

### Payment Gateway
```go
func (s *PaymentService) ProcessPayment(req *PaymentRequest) (*PaymentResult, error) {
    // Choose payment gateway based on business rules
    gateway := s.selectPaymentGateway(req.Amount, req.Method)

    // Process payment
    result, err := gateway.ProcessPayment(req)
    if err != nil {
        return nil, fmt.Errorf("payment failed: %v", err)
    }

    // Store transaction record
    if err := s.savePaymentRecord(req, result); err != nil {
        // Log error but don't fail - payment already processed
        log.Printf("Failed to save payment record: %v", err)
    }

    return result, nil
}
```

## Testing Guidelines

### Service Testing Pattern
```go
func TestCreateKoperasi(t *testing.T) {
    // Setup mocks
    mockRepo := &mocks.KoperasiRepository{}
    mockSequence := &mocks.SequenceService{}

    service := NewKoperasiService(mockRepo, mockSequence)

    // Test data
    req := &CreateKoperasiRequest{
        Nama: "Test Koperasi",
        // ... other fields
    }

    // Setup expectations
    mockRepo.On("CreateKoperasi", mock.AnythingOfType("*models.Koperasi")).Return(nil)
    mockSequence.On("GetNextNumber", mock.Anything, mock.Anything, "niak").Return(uint64(1), nil)

    // Execute
    result, err := service.CreateKoperasi(req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, req.Nama, result.Nama)
}
```

## Performance Considerations

- Use pagination dalam service methods
- Implement caching untuk frequently accessed data
- Optimize database queries di repository calls
- Use background processing untuk heavy operations
- Implement rate limiting untuk external API calls

## Security Guidelines

- Validate all inputs thoroughly
- Implement proper authorization checks
- Never expose sensitive data dalam responses
- Log security-relevant operations
- Use encryption untuk sensitive data storage