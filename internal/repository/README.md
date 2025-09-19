# Repository Layer - Data Access Implementation

## Purpose
Repository layer bertanggung jawab untuk semua operasi database. Layer ini menyediakan abstraksi untuk data access dan mengimplementasikan persistence logic menggunakan GORM untuk PostgreSQL.

## Structure Pattern
```
repository/
├── postgres/
│   ├── user_repository.go           # User data operations
│   ├── koperasi_repository.go       # Koperasi CRUD operations
│   ├── financial_repository.go      # Financial data management
│   ├── simpan_pinjam_repository.go  # Savings/loans data
│   └── klinik_repository.go         # Healthcare data
└── cassandra/                       # Future analytics data
    └── audit_repository.go          # Audit logs
```

## Repository Architecture

### Repository Structure
```go
type EntityRepository struct {
    db *gorm.DB
}

func NewEntityRepository(db *gorm.DB) *EntityRepository {
    return &EntityRepository{db: db}
}
```

### CRUD Operation Patterns

#### Create Operations
```go
func (r *EntityRepository) CreateEntity(entity *models.Entity) error {
    return r.db.Create(entity).Error
}

func (r *EntityRepository) BulkCreateEntities(entities []models.Entity) error {
    return r.db.CreateInBatches(entities, 100).Error
}
```

#### Read Operations
```go
// Single record
func (r *EntityRepository) GetEntityByID(id uint64) (*models.Entity, error) {
    var entity models.Entity
    err := r.db.First(&entity, id).Error
    if err != nil {
        return nil, err
    }
    return &entity, nil
}

// With preloading
func (r *EntityRepository) GetEntityWithRelations(id uint64) (*models.Entity, error) {
    var entity models.Entity
    err := r.db.Preload("Relation1").Preload("Relation2").First(&entity, id).Error
    if err != nil {
        return nil, err
    }
    return &entity, nil
}

// List with filters
func (r *EntityRepository) GetEntitiesByFilters(
    tenantID uint64,
    filters EntityFilters,
    limit, offset int,
) ([]models.Entity, error) {
    var entities []models.Entity

    query := r.db.Where("tenant_id = ?", tenantID)

    if filters.Name != "" {
        query = query.Where("name ILIKE ?", "%"+filters.Name+"%")
    }

    if filters.Status != "" {
        query = query.Where("status = ?", filters.Status)
    }

    if !filters.CreatedAfter.IsZero() {
        query = query.Where("created_at >= ?", filters.CreatedAfter)
    }

    err := query.Limit(limit).Offset(offset).Find(&entities).Error
    return entities, err
}
```

#### Update Operations
```go
func (r *EntityRepository) UpdateEntity(entity *models.Entity) error {
    return r.db.Save(entity).Error
}

func (r *EntityRepository) UpdateEntityFields(id uint64, updates map[string]interface{}) error {
    return r.db.Model(&models.Entity{}).Where("id = ?", id).Updates(updates).Error
}
```

#### Delete Operations
```go
// Soft delete
func (r *EntityRepository) DeleteEntity(id uint64) error {
    return r.db.Delete(&models.Entity{}, id).Error
}

// Hard delete
func (r *EntityRepository) HardDeleteEntity(id uint64) error {
    return r.db.Unscoped().Delete(&models.Entity{}, id).Error
}
```

## Advanced Query Patterns

### Complex Joins
```go
func (r *KoperasiRepository) GetKoperasiWithMembers(koperasiID uint64) (*models.Koperasi, error) {
    var koperasi models.Koperasi
    err := r.db.
        Preload("Anggota", "status = ?", "aktif").
        Preload("Provinsi").
        Preload("Kabupaten").
        First(&koperasi, koperasiID).Error

    if err != nil {
        return nil, err
    }
    return &koperasi, nil
}
```

### Aggregations
```go
func (r *FinancialRepository) GetAccountBalanceSummary(koperasiID uint64) (*BalanceSummary, error) {
    var summary BalanceSummary

    err := r.db.Table("jurnal_detail jd").
        Select(`
            SUM(CASE WHEN ca.saldo_normal = 'debit' THEN jd.debit - jd.kredit ELSE jd.kredit - jd.debit END) as total_balance,
            COUNT(DISTINCT ca.id) as account_count
        `).
        Joins("JOIN coa_akun ca ON jd.akun_id = ca.id").
        Joins("JOIN jurnal_umum ju ON jd.jurnal_id = ju.id").
        Where("ca.koperasi_id = ? AND ju.status = 'posted'", koperasiID).
        Scan(&summary).Error

    return &summary, err
}
```

### Raw SQL Queries
```go
func (r *FinancialRepository) GetFinancialReport(
    koperasiID uint64,
    startDate, endDate time.Time,
) ([]ReportItem, error) {
    var items []ReportItem

    query := `
        SELECT
            ca.kode_akun,
            ca.nama_akun,
            SUM(jd.debit) as total_debit,
            SUM(jd.kredit) as total_kredit,
            SUM(
                CASE WHEN ca.saldo_normal = 'debit'
                THEN jd.debit - jd.kredit
                ELSE jd.kredit - jd.debit
                END
            ) as saldo
        FROM coa_akun ca
        LEFT JOIN jurnal_detail jd ON ca.id = jd.akun_id
        LEFT JOIN jurnal_umum ju ON jd.jurnal_id = ju.id
        WHERE ca.koperasi_id = ?
            AND ju.status = 'posted'
            AND ju.tanggal_transaksi BETWEEN ? AND ?
        GROUP BY ca.id, ca.kode_akun, ca.nama_akun
        ORDER BY ca.kode_akun
    `

    err := r.db.Raw(query, koperasiID, startDate, endDate).Scan(&items).Error
    return items, err
}
```

## Transaction Management

### Using GORM Transactions
```go
func (r *FinancialRepository) CreateJournalWithDetails(
    journal *models.JurnalUmum,
    details []models.JurnalDetail,
) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // Create journal header
        if err := tx.Create(journal).Error; err != nil {
            return err
        }

        // Set journal ID for details
        for i := range details {
            details[i].JurnalID = journal.ID
        }

        // Create journal details
        if err := tx.CreateInBatches(details, 50).Error; err != nil {
            return err
        }

        return nil
    })
}
```

### External Transaction
```go
func (r *Repository) CreateWithExternalTx(tx *gorm.DB, entity *models.Entity) error {
    return tx.Create(entity).Error
}

// Usage in service
func (s *Service) ComplexOperation() error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        if err := s.repo1.CreateWithExternalTx(tx, entity1); err != nil {
            return err
        }
        if err := s.repo2.CreateWithExternalTx(tx, entity2); err != nil {
            return err
        }
        return nil
    })
}
```

## Performance Optimization

### Preloading Strategies
```go
// Eager loading dengan select
func (r *Repository) GetEntityWithOptimizedPreload(id uint64) (*models.Entity, error) {
    var entity models.Entity
    err := r.db.
        Select("id, name, status, created_at"). // Only needed fields
        Preload("Relations", func(db *gorm.DB) *gorm.DB {
            return db.Select("id, entity_id, name") // Only needed relation fields
        }).
        First(&entity, id).Error

    return &entity, err
}
```

### Batch Operations
```go
func (r *Repository) BulkUpdateStatus(ids []uint64, status string) error {
    return r.db.Model(&models.Entity{}).
        Where("id IN ?", ids).
        Update("status", status).Error
}
```

### Indexing Hints
```go
func (r *Repository) GetEntitiesWithIndex(tenantID uint64) ([]models.Entity, error) {
    var entities []models.Entity

    // Use specific index for performance
    err := r.db.
        Set("gorm:query_hint", "USE INDEX (idx_tenant_status)").
        Where("tenant_id = ? AND status = ?", tenantID, "active").
        Find(&entities).Error

    return entities, err
}
```

## Multi-tenant Patterns

### Tenant Isolation
```go
func (r *Repository) withTenantScope(tenantID uint64) *gorm.DB {
    return r.db.Where("tenant_id = ?", tenantID)
}

func (r *Repository) GetEntitiesByTenant(tenantID uint64) ([]models.Entity, error) {
    var entities []models.Entity
    err := r.withTenantScope(tenantID).Find(&entities).Error
    return entities, err
}
```

### Global Scopes
```go
func (r *Repository) addCommonScopes(query *gorm.DB) *gorm.DB {
    return query.
        Where("deleted_at IS NULL").
        Where("is_active = ?", true)
}
```

## Koperasi-Specific Patterns

### NIAK Validation
```go
func (r *KoperasiRepository) IsNIAKExists(niak string) (bool, error) {
    var count int64
    err := r.db.Model(&models.Koperasi{}).Where("niak = ?", niak).Count(&count).Error
    return count > 0, err
}
```

### Member Queries
```go
func (r *KoperasiRepository) GetActiveMembers(koperasiID uint64) ([]models.AnggotaKoperasi, error) {
    var members []models.AnggotaKoperasi
    err := r.db.
        Where("koperasi_id = ? AND status = ?", koperasiID, "aktif").
        Order("nomor_anggota ASC").
        Find(&members).Error
    return members, err
}
```

### Financial Calculations
```go
func (r *FinancialRepository) CalculateAccountBalance(
    accountID uint64,
    upToDate time.Time,
) (float64, error) {
    var result struct {
        Balance float64
    }

    err := r.db.Table("jurnal_detail jd").
        Select("SUM(jd.debit - jd.kredit) as balance").
        Joins("JOIN jurnal_umum ju ON jd.jurnal_id = ju.id").
        Where("jd.akun_id = ? AND ju.status = 'posted' AND ju.tanggal_transaksi <= ?",
            accountID, upToDate).
        Scan(&result).Error

    return result.Balance, err
}
```

## Error Handling Patterns

### GORM Error Handling
```go
func (r *Repository) GetEntityByID(id uint64) (*models.Entity, error) {
    var entity models.Entity
    err := r.db.First(&entity, id).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("entity with ID %d not found", id)
        }
        return nil, fmt.Errorf("database error: %v", err)
    }

    return &entity, nil
}
```

### Duplicate Key Handling
```go
func (r *Repository) CreateEntitySafe(entity *models.Entity) error {
    err := r.db.Create(entity).Error
    if err != nil {
        if strings.Contains(err.Error(), "duplicate key") {
            return fmt.Errorf("entity already exists")
        }
        return fmt.Errorf("failed to create entity: %v", err)
    }
    return nil
}
```

## Testing Patterns

### Repository Testing with Test Database
```go
func TestCreateEntity(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    repo := NewEntityRepository(db)

    // Test data
    entity := &models.Entity{
        Name: "Test Entity",
        Status: "active",
    }

    // Execute
    err := repo.CreateEntity(entity)

    // Assert
    assert.NoError(t, err)
    assert.NotZero(t, entity.ID)

    // Verify in database
    var saved models.Entity
    db.First(&saved, entity.ID)
    assert.Equal(t, entity.Name, saved.Name)
}
```

### Mock Repository Pattern
```go
type MockEntityRepository struct {
    mock.Mock
}

func (m *MockEntityRepository) CreateEntity(entity *models.Entity) error {
    args := m.Called(entity)
    return args.Error(0)
}

func (m *MockEntityRepository) GetEntityByID(id uint64) (*models.Entity, error) {
    args := m.Called(id)
    return args.Get(0).(*models.Entity), args.Error(1)
}
```

## Security Considerations

### SQL Injection Prevention
```go
// Good - using parameterized queries
func (r *Repository) GetEntitiesByName(name string) ([]models.Entity, error) {
    var entities []models.Entity
    err := r.db.Where("name ILIKE ?", "%"+name+"%").Find(&entities).Error
    return entities, err
}

// Bad - vulnerable to SQL injection
func (r *Repository) GetEntitiesByNameUnsafe(name string) ([]models.Entity, error) {
    var entities []models.Entity
    query := fmt.Sprintf("SELECT * FROM entities WHERE name LIKE '%%%s%%'", name)
    err := r.db.Raw(query).Scan(&entities).Error
    return entities, err
}
```

### Data Sanitization
```go
func (r *Repository) sanitizeInput(input string) string {
    // Remove potentially dangerous characters
    input = strings.ReplaceAll(input, "'", "''")
    input = strings.TrimSpace(input)
    return input
}
```

## Performance Monitoring

### Query Logging
```go
func (r *Repository) enableQueryLogging() {
    r.db = r.db.Debug() // Enable SQL logging in development
}
```

### Slow Query Detection
```go
func (r *Repository) withTimeout(timeout time.Duration) *gorm.DB {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    return r.db.WithContext(ctx)
}
```

## Best Practices

1. **Always handle GORM errors properly**
2. **Use transactions for related operations**
3. **Implement proper indexing strategies**
4. **Use preloading to avoid N+1 queries**
5. **Implement pagination for large datasets**
6. **Use batch operations for bulk data**
7. **Sanitize all inputs**
8. **Test with realistic data volumes**
9. **Monitor query performance**
10. **Implement proper logging**