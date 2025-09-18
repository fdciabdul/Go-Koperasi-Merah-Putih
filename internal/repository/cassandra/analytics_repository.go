package cassandra

import (
	"time"

	"github.com/gocql/gocql"
	"koperasi-merah-putih/internal/models/cassandra"
)

type AnalyticsRepository struct {
	session *gocql.Session
}

func NewAnalyticsRepository(session *gocql.Session) *AnalyticsRepository {
	return &AnalyticsRepository{session: session}
}

func (r *AnalyticsRepository) InsertTransactionLog(log *cassandra.TransactionLog) error {
	query := `INSERT INTO transaction_logs (id, tenant_id, koperasi_id, transaction_type, transaction_id,
			  user_id, amount, status, description, metadata, ip_address, user_agent, created_at, year, month, day)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return r.session.Query(query,
		log.ID, log.TenantID, log.KoperasiID, log.TransactionType, log.TransactionID,
		log.UserID, log.Amount, log.Status, log.Description, log.Metadata,
		log.IPAddress, log.UserAgent, log.CreatedAt, log.Year, log.Month, log.Day).Exec()
}

func (r *AnalyticsRepository) GetTransactionLogs(koperasiID uint64, year, month int, limit int) ([]cassandra.TransactionLog, error) {
	var logs []cassandra.TransactionLog

	query := `SELECT id, tenant_id, koperasi_id, transaction_type, transaction_id, user_id, amount, status,
			  description, metadata, ip_address, user_agent, created_at, year, month, day
			  FROM transaction_logs WHERE koperasi_id = ? AND year = ? AND month = ? LIMIT ?`

	iter := r.session.Query(query, koperasiID, year, month, limit).Iter()
	defer iter.Close()

	var log cassandra.TransactionLog
	for iter.Scan(&log.ID, &log.TenantID, &log.KoperasiID, &log.TransactionType, &log.TransactionID,
		&log.UserID, &log.Amount, &log.Status, &log.Description, &log.Metadata,
		&log.IPAddress, &log.UserAgent, &log.CreatedAt, &log.Year, &log.Month, &log.Day) {
		logs = append(logs, log)
	}

	return logs, iter.Close()
}

func (r *AnalyticsRepository) InsertPaymentAnalytics(analytics *cassandra.PaymentAnalytics) error {
	query := `INSERT INTO payment_analytics (id, tenant_id, koperasi_id, payment_provider, payment_method,
			  transaction_type, total_amount, admin_fee, transaction_count, success_count, failed_count,
			  success_rate, created_at, year, month, day)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return r.session.Query(query,
		analytics.ID, analytics.TenantID, analytics.KoperasiID, analytics.PaymentProvider,
		analytics.PaymentMethod, analytics.TransactionType, analytics.TotalAmount, analytics.AdminFee,
		analytics.TransactionCount, analytics.SuccessCount, analytics.FailedCount, analytics.SuccessRate,
		analytics.CreatedAt, analytics.Year, analytics.Month, analytics.Day).Exec()
}

func (r *AnalyticsRepository) InsertPPOBAnalytics(analytics *cassandra.PPOBAnalytics) error {
	query := `INSERT INTO ppob_analytics (id, tenant_id, koperasi_id, product_category, product_name,
			  provider, transaction_count, total_revenue, total_commission, success_count, failed_count,
			  success_rate, created_at, year, month, day)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return r.session.Query(query,
		analytics.ID, analytics.TenantID, analytics.KoperasiID, analytics.ProductCategory,
		analytics.ProductName, analytics.Provider, analytics.TransactionCount, analytics.TotalRevenue,
		analytics.TotalCommission, analytics.SuccessCount, analytics.FailedCount, analytics.SuccessRate,
		analytics.CreatedAt, analytics.Year, analytics.Month, analytics.Day).Exec()
}

func (r *AnalyticsRepository) GetPPOBAnalytics(koperasiID uint64, year, month int) ([]cassandra.PPOBAnalytics, error) {
	var analytics []cassandra.PPOBAnalytics

	query := `SELECT id, tenant_id, koperasi_id, product_category, product_name, provider,
			  transaction_count, total_revenue, total_commission, success_count, failed_count,
			  success_rate, created_at, year, month, day
			  FROM ppob_analytics WHERE koperasi_id = ? AND year = ? AND month = ?`

	iter := r.session.Query(query, koperasiID, year, month).Iter()
	defer iter.Close()

	var analytic cassandra.PPOBAnalytics
	for iter.Scan(&analytic.ID, &analytic.TenantID, &analytic.KoperasiID, &analytic.ProductCategory,
		&analytic.ProductName, &analytic.Provider, &analytic.TransactionCount, &analytic.TotalRevenue,
		&analytic.TotalCommission, &analytic.SuccessCount, &analytic.FailedCount, &analytic.SuccessRate,
		&analytic.CreatedAt, &analytic.Year, &analytic.Month, &analytic.Day) {
		analytics = append(analytics, analytic)
	}

	return analytics, iter.Close()
}

func (r *AnalyticsRepository) InsertUserActivityLog(log *cassandra.UserActivityLog) error {
	query := `INSERT INTO user_activity_logs (id, tenant_id, koperasi_id, user_id, username, activity,
			  module, description, ip_address, user_agent, session_id, created_at, year, month, day, hour)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return r.session.Query(query,
		log.ID, log.TenantID, log.KoperasiID, log.UserID, log.Username, log.Activity,
		log.Module, log.Description, log.IPAddress, log.UserAgent, log.SessionID,
		log.CreatedAt, log.Year, log.Month, log.Day, log.Hour).Exec()
}

func (r *AnalyticsRepository) InsertErrorLog(log *cassandra.ErrorLog) error {
	query := `INSERT INTO error_logs (id, tenant_id, koperasi_id, user_id, error_type, error_message,
			  stack_trace, request_data, response_data, module, function, ip_address, user_agent,
			  created_at, year, month, day)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return r.session.Query(query,
		log.ID, log.TenantID, log.KoperasiID, log.UserID, log.ErrorType, log.ErrorMessage,
		log.StackTrace, log.RequestData, log.ResponseData, log.Module, log.Function,
		log.IPAddress, log.UserAgent, log.CreatedAt, log.Year, log.Month, log.Day).Exec()
}

func (r *AnalyticsRepository) InsertPerformanceMetrics(metrics *cassandra.PerformanceMetrics) error {
	query := `INSERT INTO performance_metrics (id, tenant_id, koperasi_id, metric_name, metric_value,
			  metric_unit, endpoint_path, http_method, response_time, status_code, request_size,
			  response_size, created_at, year, month, day, hour, minute)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return r.session.Query(query,
		metrics.ID, metrics.TenantID, metrics.KoperasiID, metrics.MetricName, metrics.MetricValue,
		metrics.MetricUnit, metrics.EndpointPath, metrics.HttpMethod, metrics.ResponseTime,
		metrics.StatusCode, metrics.RequestSize, metrics.ResponseSize, metrics.CreatedAt,
		metrics.Year, metrics.Month, metrics.Day, metrics.Hour, metrics.Minute).Exec()
}

func (r *AnalyticsRepository) InsertFactKeuanganBulanan(fact *cassandra.FactKeuanganBulanan) error {
	query := `INSERT INTO fact_keuangan_bulanan (id, tenant_id, koperasi_id, periode_tahun, periode_bulan,
			  total_simpanan, total_pinjaman, total_angsuran, total_transaksi_ppob, total_fee_ppob,
			  jumlah_transaksi_ppob, total_pendapatan_klinik, jumlah_kunjungan_klinik, total_aset,
			  total_kewajiban, total_ekuitas, total_pendapatan, total_beban, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return r.session.Query(query,
		fact.ID, fact.TenantID, fact.KoperasiID, fact.PeriodeTahun, fact.PeriodeBulan,
		fact.TotalSimpanan, fact.TotalPinjaman, fact.TotalAngsuran, fact.TotalTransaksiPPOB,
		fact.TotalFeePPOB, fact.JumlahTransaksiPPOB, fact.TotalPendapatanKlinik, fact.JumlahKunjunganKlinik,
		fact.TotalAset, fact.TotalKewajiban, fact.TotalEkuitas, fact.TotalPendapatan, fact.TotalBeban,
		fact.CreatedAt).Exec()
}

func (r *AnalyticsRepository) GetFactKeuanganBulanan(koperasiID uint64, tahun int) ([]cassandra.FactKeuanganBulanan, error) {
	var facts []cassandra.FactKeuanganBulanan

	query := `SELECT id, tenant_id, koperasi_id, periode_tahun, periode_bulan, total_simpanan,
			  total_pinjaman, total_angsuran, total_transaksi_ppob, total_fee_ppob, jumlah_transaksi_ppob,
			  total_pendapatan_klinik, jumlah_kunjungan_klinik, total_aset, total_kewajiban, total_ekuitas,
			  total_pendapatan, total_beban, created_at
			  FROM fact_keuangan_bulanan WHERE koperasi_id = ? AND periode_tahun = ?`

	iter := r.session.Query(query, koperasiID, tahun).Iter()
	defer iter.Close()

	var fact cassandra.FactKeuanganBulanan
	for iter.Scan(&fact.ID, &fact.TenantID, &fact.KoperasiID, &fact.PeriodeTahun, &fact.PeriodeBulan,
		&fact.TotalSimpanan, &fact.TotalPinjaman, &fact.TotalAngsuran, &fact.TotalTransaksiPPOB,
		&fact.TotalFeePPOB, &fact.JumlahTransaksiPPOB, &fact.TotalPendapatanKlinik, &fact.JumlahKunjunganKlinik,
		&fact.TotalAset, &fact.TotalKewajiban, &fact.TotalEkuitas, &fact.TotalPendapatan, &fact.TotalBeban,
		&fact.CreatedAt) {
		facts = append(facts, fact)
	}

	return facts, iter.Close()
}