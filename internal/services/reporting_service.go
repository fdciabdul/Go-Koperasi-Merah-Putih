package services

import (
	"fmt"
	"time"

	"koperasi-merah-putih/internal/cache"
	"koperasi-merah-putih/internal/models/postgres"
	"koperasi-merah-putih/internal/repository"
)

type ReportingService struct {
	koperasiRepo     *repository.KoperasiRepository
	produkRepo       *repository.ProdukRepository
	simpanPinjamRepo *repository.SimpanPinjamRepository
	financialRepo    *repository.FinancialRepository
	klinikRepo       *repository.KlinikRepository
	cache            *cache.RedisCache
}

func NewReportingService(
	koperasiRepo *repository.KoperasiRepository,
	produkRepo *repository.ProdukRepository,
	simpanPinjamRepo *repository.SimpanPinjamRepository,
	financialRepo *repository.FinancialRepository,
	klinikRepo *repository.KlinikRepository,
	cache *cache.RedisCache,
) *ReportingService {
	return &ReportingService{
		koperasiRepo:     koperasiRepo,
		produkRepo:       produkRepo,
		simpanPinjamRepo: simpanPinjamRepo,
		financialRepo:    financialRepo,
		klinikRepo:       klinikRepo,
		cache:            cache,
	}
}

// Dashboard Analytics
type DashboardData struct {
	KoperasiID       uint64                 `json:"koperasi_id"`
	Period           string                 `json:"period"`
	GeneratedAt      time.Time              `json:"generated_at"`
	Summary          DashboardSummary       `json:"summary"`
	Revenue          RevenueAnalytics       `json:"revenue"`
	ProductMetrics   ProductMetrics         `json:"product_metrics"`
	MemberMetrics    MemberMetrics          `json:"member_metrics"`
	FinancialMetrics FinancialMetrics       `json:"financial_metrics"`
	Charts           map[string]interface{} `json:"charts"`
}

type DashboardSummary struct {
	TotalRevenue         float64 `json:"total_revenue"`
	TotalExpenses        float64 `json:"total_expenses"`
	NetProfit            float64 `json:"net_profit"`
	TotalMembers         int     `json:"total_members"`
	ActiveMembers        int     `json:"active_members"`
	TotalProducts        int     `json:"total_products"`
	LowStockProducts     int     `json:"low_stock_products"`
	TotalTransactions    int     `json:"total_transactions"`
	OutstandingLoans     float64 `json:"outstanding_loans"`
	TotalSavings         float64 `json:"total_savings"`
	PendingApprovals     int     `json:"pending_approvals"`
	RevenueGrowth        float64 `json:"revenue_growth_percentage"`
	MemberGrowth         float64 `json:"member_growth_percentage"`
}

type RevenueAnalytics struct {
	DailySales      []DailySales      `json:"daily_sales"`
	MonthlySales    []MonthlySales    `json:"monthly_sales"`
	TopProducts     []TopProduct      `json:"top_products"`
	TopCategories   []TopCategory     `json:"top_categories"`
	PaymentMethods  []PaymentMethod   `json:"payment_methods"`
}

type ProductMetrics struct {
	TotalProducts      int                `json:"total_products"`
	ActiveProducts     int                `json:"active_products"`
	LowStockAlert      []LowStockProduct  `json:"low_stock_alert"`
	ExpiringProducts   []ExpiringProduct  `json:"expiring_products"`
	FastMovingProducts []FastMovingProduct `json:"fast_moving_products"`
	SlowMovingProducts []SlowMovingProduct `json:"slow_moving_products"`
	StockValue         float64            `json:"total_stock_value"`
	CategoryBreakdown  []CategoryStock    `json:"category_breakdown"`
}

type MemberMetrics struct {
	TotalMembers       int              `json:"total_members"`
	ActiveMembers      int              `json:"active_members"`
	NewMembersThisMonth int             `json:"new_members_this_month"`
	MembersByStatus    []MemberStatus   `json:"members_by_status"`
	TopMembers         []TopMember      `json:"top_members"`
	MemberRetention    float64          `json:"retention_rate"`
	AverageTransaction float64          `json:"average_transaction"`
}

type FinancialMetrics struct {
	CashFlow         CashFlow         `json:"cash_flow"`
	ProfitLoss       ProfitLoss       `json:"profit_loss"`
	BalanceSheet     BalanceSheetSummary `json:"balance_sheet"`
	LoanMetrics      LoanMetrics      `json:"loan_metrics"`
	SavingsMetrics   SavingsMetrics   `json:"savings_metrics"`
}

// Supporting structs
type DailySales struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
	Count  int     `json:"count"`
}

type MonthlySales struct {
	Month  string  `json:"month"`
	Year   int     `json:"year"`
	Amount float64 `json:"amount"`
	Count  int     `json:"count"`
}

type TopProduct struct {
	ProductID   uint64  `json:"product_id"`
	ProductName string  `json:"product_name"`
	TotalSold   int     `json:"total_sold"`
	Revenue     float64 `json:"revenue"`
}

type TopCategory struct {
	CategoryID   uint64  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalSold    int     `json:"total_sold"`
	Revenue      float64 `json:"revenue"`
}

type PaymentMethod struct {
	Method      string  `json:"method"`
	Count       int     `json:"count"`
	TotalAmount float64 `json:"total_amount"`
	Percentage  float64 `json:"percentage"`
}

type LowStockProduct struct {
	ProductID    uint64 `json:"product_id"`
	ProductName  string `json:"product_name"`
	CurrentStock int    `json:"current_stock"`
	MinimalStock int    `json:"minimal_stock"`
	Status       string `json:"status"`
}

type ExpiringProduct struct {
	ProductID      uint64    `json:"product_id"`
	ProductName    string    `json:"product_name"`
	ExpiryDate     time.Time `json:"expiry_date"`
	DaysRemaining  int       `json:"days_remaining"`
	CurrentStock   int       `json:"current_stock"`
}

type FastMovingProduct struct {
	ProductID    uint64  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	TurnoverRate float64 `json:"turnover_rate"`
	TotalSold    int     `json:"total_sold"`
}

type SlowMovingProduct struct {
	ProductID     uint64  `json:"product_id"`
	ProductName   string  `json:"product_name"`
	DaysSinceLastSold int `json:"days_since_last_sold"`
	CurrentStock  int     `json:"current_stock"`
	StockValue    float64 `json:"stock_value"`
}

type CategoryStock struct {
	CategoryID   uint64  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalItems   int     `json:"total_items"`
	StockValue   float64 `json:"stock_value"`
	Percentage   float64 `json:"percentage"`
}

type MemberStatus struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type TopMember struct {
	MemberID         uint64  `json:"member_id"`
	MemberName       string  `json:"member_name"`
	TotalTransaction float64 `json:"total_transaction"`
	TransactionCount int     `json:"transaction_count"`
}

type CashFlow struct {
	Period      string  `json:"period"`
	CashInflow  float64 `json:"cash_inflow"`
	CashOutflow float64 `json:"cash_outflow"`
	NetCashFlow float64 `json:"net_cash_flow"`
}

type ProfitLoss struct {
	Revenue       float64 `json:"revenue"`
	CostOfGoods   float64 `json:"cost_of_goods"`
	GrossProfit   float64 `json:"gross_profit"`
	OperatingExp  float64 `json:"operating_expenses"`
	NetProfit     float64 `json:"net_profit"`
	ProfitMargin  float64 `json:"profit_margin"`
}

type BalanceSheetSummary struct {
	TotalAssets      float64 `json:"total_assets"`
	TotalLiabilities float64 `json:"total_liabilities"`
	TotalEquity      float64 `json:"total_equity"`
	CurrentRatio     float64 `json:"current_ratio"`
	DebtToEquity     float64 `json:"debt_to_equity"`
}

type LoanMetrics struct {
	TotalLoans          float64 `json:"total_loans"`
	OutstandingLoans    float64 `json:"outstanding_loans"`
	OverdueLoans        float64 `json:"overdue_loans"`
	CollectionRate      float64 `json:"collection_rate"`
	AverageInterestRate float64 `json:"average_interest_rate"`
	NPLRatio            float64 `json:"npl_ratio"`
}

type SavingsMetrics struct {
	TotalSavings        float64 `json:"total_savings"`
	AverageSavings      float64 `json:"average_savings"`
	SavingsGrowth       float64 `json:"savings_growth"`
	ActiveSaversCount   int     `json:"active_savers_count"`
}

// Main dashboard function
func (s *ReportingService) GetDashboard(koperasiID uint64, period string) (*DashboardData, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("dashboard:%d:%s", koperasiID, period)
	var cachedData DashboardData
	if err := s.cache.Get(cacheKey, &cachedData); err == nil {
		return &cachedData, nil
	}

	dashboard := &DashboardData{
		KoperasiID:  koperasiID,
		Period:      period,
		GeneratedAt: time.Now(),
	}

	// Load all metrics (in production, these would be actual DB queries)
	if err := s.loadSummary(dashboard); err != nil {
		return nil, err
	}

	if err := s.loadRevenueAnalytics(dashboard); err != nil {
		return nil, err
	}

	if err := s.loadProductMetrics(dashboard); err != nil {
		return nil, err
	}

	if err := s.loadMemberMetrics(dashboard); err != nil {
		return nil, err
	}

	if err := s.loadFinancialMetrics(dashboard); err != nil {
		return nil, err
	}

	if err := s.loadCharts(dashboard); err != nil {
		return nil, err
	}

	// Cache the result
	s.cache.Set(cacheKey, dashboard, 5*time.Minute)

	return dashboard, nil
}

func (s *ReportingService) loadSummary(dashboard *DashboardData) error {
	// This would be actual database queries in production
	dashboard.Summary = DashboardSummary{
		TotalRevenue:      1250000000,
		TotalExpenses:     980000000,
		NetProfit:         270000000,
		TotalMembers:      1250,
		ActiveMembers:     980,
		TotalProducts:     450,
		LowStockProducts:  12,
		TotalTransactions: 8500,
		OutstandingLoans:  450000000,
		TotalSavings:      780000000,
		PendingApprovals:  5,
		RevenueGrowth:     12.5,
		MemberGrowth:      8.3,
	}
	return nil
}

func (s *ReportingService) loadRevenueAnalytics(dashboard *DashboardData) error {
	// Sample data - would be actual queries in production
	dashboard.Revenue = RevenueAnalytics{
		DailySales: []DailySales{
			{Date: "2024-01-01", Amount: 45000000, Count: 120},
			{Date: "2024-01-02", Amount: 38000000, Count: 105},
			{Date: "2024-01-03", Amount: 52000000, Count: 145},
		},
		MonthlySales: []MonthlySales{
			{Month: "January", Year: 2024, Amount: 380000000, Count: 3200},
			{Month: "February", Year: 2024, Amount: 420000000, Count: 3500},
			{Month: "March", Year: 2024, Amount: 450000000, Count: 3800},
		},
		TopProducts: []TopProduct{
			{ProductID: 1, ProductName: "Beras Premium 5kg", TotalSold: 2500, Revenue: 130000000},
			{ProductID: 2, ProductName: "Minyak Goreng 2L", TotalSold: 1800, Revenue: 90000000},
		},
		TopCategories: []TopCategory{
			{CategoryID: 1, CategoryName: "Sembako", TotalSold: 8500, Revenue: 450000000},
			{CategoryID: 2, CategoryName: "Sayuran", TotalSold: 5200, Revenue: 180000000},
		},
		PaymentMethods: []PaymentMethod{
			{Method: "cash", Count: 5200, TotalAmount: 650000000, Percentage: 52},
			{Method: "transfer", Count: 3200, TotalAmount: 450000000, Percentage: 36},
			{Method: "simpanan", Count: 1100, TotalAmount: 150000000, Percentage: 12},
		},
	}
	return nil
}

func (s *ReportingService) loadProductMetrics(dashboard *DashboardData) error {
	products, err := s.produkRepo.GetStokReport(dashboard.KoperasiID)
	if err != nil {
		return err
	}

	lowStock, _ := s.produkRepo.GetProdukStokRendah(dashboard.KoperasiID)
	expiring, _ := s.produkRepo.GetProdukExpiringSoon(dashboard.KoperasiID, 30)

	dashboard.ProductMetrics = ProductMetrics{
		TotalProducts:  len(products),
		ActiveProducts: len(products),
		StockValue:     calculateStockValue(products),
		LowStockAlert:  mapLowStockProducts(lowStock),
		ExpiringProducts: mapExpiringProducts(expiring),
	}
	return nil
}

func (s *ReportingService) loadMemberMetrics(dashboard *DashboardData) error {
	// Would be actual queries
	dashboard.MemberMetrics = MemberMetrics{
		TotalMembers:        1250,
		ActiveMembers:       980,
		NewMembersThisMonth: 45,
		MemberRetention:     78.4,
		AverageTransaction:  450000,
		MembersByStatus: []MemberStatus{
			{Status: "active", Count: 980},
			{Status: "inactive", Count: 220},
			{Status: "suspended", Count: 50},
		},
	}
	return nil
}

func (s *ReportingService) loadFinancialMetrics(dashboard *DashboardData) error {
	dashboard.FinancialMetrics = FinancialMetrics{
		CashFlow: CashFlow{
			Period:      dashboard.Period,
			CashInflow:  1250000000,
			CashOutflow: 980000000,
			NetCashFlow: 270000000,
		},
		ProfitLoss: ProfitLoss{
			Revenue:      1250000000,
			CostOfGoods:  750000000,
			GrossProfit:  500000000,
			OperatingExp: 230000000,
			NetProfit:    270000000,
			ProfitMargin: 21.6,
		},
		BalanceSheet: BalanceSheetSummary{
			TotalAssets:      5500000000,
			TotalLiabilities: 2200000000,
			TotalEquity:      3300000000,
			CurrentRatio:     2.1,
			DebtToEquity:     0.67,
		},
		LoanMetrics: LoanMetrics{
			TotalLoans:          850000000,
			OutstandingLoans:    450000000,
			OverdueLoans:        35000000,
			CollectionRate:      92.5,
			AverageInterestRate: 12.5,
			NPLRatio:            4.1,
		},
		SavingsMetrics: SavingsMetrics{
			TotalSavings:      780000000,
			AverageSavings:    795918,
			SavingsGrowth:     8.5,
			ActiveSaversCount: 980,
		},
	}
	return nil
}

func (s *ReportingService) loadCharts(dashboard *DashboardData) error {
	dashboard.Charts = map[string]interface{}{
		"revenue_trend": []map[string]interface{}{
			{"month": "Jan", "revenue": 380000000},
			{"month": "Feb", "revenue": 420000000},
			{"month": "Mar", "revenue": 450000000},
		},
		"member_growth": []map[string]interface{}{
			{"month": "Jan", "count": 1180},
			{"month": "Feb", "count": 1205},
			{"month": "Mar", "count": 1250},
		},
		"category_distribution": []map[string]interface{}{
			{"category": "Sembako", "value": 45},
			{"category": "Sayuran", "value": 18},
			{"category": "Minuman", "value": 15},
			{"category": "Daging", "value": 12},
			{"category": "Others", "value": 10},
		},
	}
	return nil
}

// Report generation functions
func (s *ReportingService) GenerateSalesReport(koperasiID uint64, startDate, endDate time.Time) (interface{}, error) {
	cacheKey := fmt.Sprintf("report:sales:%d:%s:%s", koperasiID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	var report interface{}
	if err := s.cache.Get(cacheKey, &report); err == nil {
		return report, nil
	}

	// Generate actual report
	salesData, err := s.produkRepo.GetPenjualansByKoperasi(koperasiID, startDate, endDate, 1000, 0)
	if err != nil {
		return nil, err
	}

	report = map[string]interface{}{
		"period": map[string]string{
			"start": startDate.Format("2006-01-02"),
			"end":   endDate.Format("2006-01-02"),
		},
		"summary": map[string]interface{}{
			"total_sales":       calculateTotalSales(salesData),
			"total_transactions": len(salesData),
			"average_transaction": calculateAverageTransaction(salesData),
		},
		"details": salesData,
	}

	s.cache.Set(cacheKey, report, 10*time.Minute)
	return report, nil
}

func (s *ReportingService) GenerateInventoryReport(koperasiID uint64) (interface{}, error) {
	products, err := s.produkRepo.GetStokReport(koperasiID)
	if err != nil {
		return nil, err
	}

	lowStock, _ := s.produkRepo.GetProdukStokRendah(koperasiID)
	expiring, _ := s.produkRepo.GetProdukExpiringSoon(koperasiID, 30)

	report := map[string]interface{}{
		"generated_at": time.Now(),
		"summary": map[string]interface{}{
			"total_products": len(products),
			"total_value":    calculateStockValue(products),
			"low_stock_count": len(lowStock),
			"expiring_count": len(expiring),
		},
		"products": products,
		"alerts": map[string]interface{}{
			"low_stock": lowStock,
			"expiring":  expiring,
		},
	}

	return report, nil
}

func (s *ReportingService) GenerateFinancialReport(koperasiID uint64, reportType string, period string) (interface{}, error) {
	switch reportType {
	case "balance_sheet":
		return s.financialRepo.GetNeraca(koperasiID, parsePeriod(period))
	case "profit_loss":
		return s.financialRepo.GetLabaRugi(koperasiID, parsePeriod(period))
	case "cash_flow":
		return s.generateCashFlowReport(koperasiID, period)
	default:
		return nil, fmt.Errorf("unknown report type: %s", reportType)
	}
}

func (s *ReportingService) GenerateMemberReport(koperasiID uint64) (interface{}, error) {
	members, err := s.koperasiRepo.GetAnggotaList(koperasiID, 10000, 0)
	if err != nil {
		return nil, err
	}

	activeCount := 0
	for _, member := range members {
		if member.StatusKeanggotaan == "active" {
			activeCount++
		}
	}

	report := map[string]interface{}{
		"generated_at": time.Now(),
		"summary": map[string]interface{}{
			"total_members":  len(members),
			"active_members": activeCount,
			"inactive_members": len(members) - activeCount,
		},
		"members": members,
	}

	return report, nil
}

func (s *ReportingService) generateCashFlowReport(koperasiID uint64, period string) (interface{}, error) {
	// Simplified cash flow report
	return map[string]interface{}{
		"period": period,
		"cash_flow": map[string]interface{}{
			"operating_activities": map[string]float64{
				"cash_from_sales":    1250000000,
				"cash_to_suppliers": -750000000,
				"cash_to_employees": -120000000,
				"net_operating":      380000000,
			},
			"investing_activities": map[string]float64{
				"equipment_purchase": -50000000,
				"net_investing":      -50000000,
			},
			"financing_activities": map[string]float64{
				"loan_proceeds":   200000000,
				"loan_repayment": -150000000,
				"net_financing":   50000000,
			},
			"net_cash_flow": 380000000,
		},
	}, nil
}

// Helper functions
func calculateStockValue(products []postgres.Produk) float64 {
	total := 0.0
	for _, p := range products {
		total += float64(p.StokCurrent) * p.HargaBeli
	}
	return total
}

func mapLowStockProducts(products []postgres.Produk) []LowStockProduct {
	var result []LowStockProduct
	for _, p := range products {
		status := "warning"
		if p.StokCurrent == 0 {
			status = "critical"
		}
		result = append(result, LowStockProduct{
			ProductID:    p.ID,
			ProductName:  p.NamaProduk,
			CurrentStock: p.StokCurrent,
			MinimalStock: p.StokMinimal,
			Status:       status,
		})
	}
	return result
}

func mapExpiringProducts(products []postgres.Produk) []ExpiringProduct {
	var result []ExpiringProduct
	now := time.Now()
	for _, p := range products {
		if p.TanggalExpired != nil {
			daysRemaining := int(p.TanggalExpired.Sub(now).Hours() / 24)
			result = append(result, ExpiringProduct{
				ProductID:     p.ID,
				ProductName:   p.NamaProduk,
				ExpiryDate:    *p.TanggalExpired,
				DaysRemaining: daysRemaining,
				CurrentStock:  p.StokCurrent,
			})
		}
	}
	return result
}

func calculateTotalSales(sales []postgres.PenjualanHeader) float64 {
	total := 0.0
	for _, s := range sales {
		total += s.GrandTotal
	}
	return total
}

func calculateAverageTransaction(sales []postgres.PenjualanHeader) float64 {
	if len(sales) == 0 {
		return 0
	}
	return calculateTotalSales(sales) / float64(len(sales))
}

func parsePeriod(period string) time.Time {
	// Parse period string to time.Time
	// Format expected: "2024-01", "2024-Q1", "2024"
	t, _ := time.Parse("2006-01", period)
	return t
}