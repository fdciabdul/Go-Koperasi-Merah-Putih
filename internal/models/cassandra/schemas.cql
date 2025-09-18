-- Cassandra CQL Schema for Koperasi Merah Putih Analytics

CREATE KEYSPACE IF NOT EXISTS koperasi_analytics
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3};

USE koperasi_analytics;

-- Fact table for monthly financial data
CREATE TABLE IF NOT EXISTS fact_keuangan_bulanan (
    id UUID PRIMARY KEY,
    tenant_id BIGINT,
    koperasi_id BIGINT,
    periode_tahun INT,
    periode_bulan INT,
    total_simpanan DECIMAL,
    total_pinjaman DECIMAL,
    total_angsuran DECIMAL,
    total_transaksi_ppob DECIMAL,
    total_fee_ppob DECIMAL,
    jumlah_transaksi_ppob INT,
    total_pendapatan_klinik DECIMAL,
    jumlah_kunjungan_klinik INT,
    total_aset DECIMAL,
    total_kewajiban DECIMAL,
    total_ekuitas DECIMAL,
    total_pendapatan DECIMAL,
    total_beban DECIMAL,
    created_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS ON fact_keuangan_bulanan (koperasi_id);
CREATE INDEX IF NOT EXISTS ON fact_keuangan_bulanan (periode_tahun);
CREATE INDEX IF NOT EXISTS ON fact_keuangan_bulanan (periode_bulan);

-- Fact table for monthly member data
CREATE TABLE IF NOT EXISTS fact_anggota_bulanan (
    id UUID PRIMARY KEY,
    koperasi_id BIGINT,
    periode_tahun INT,
    periode_bulan INT,
    jumlah_anggota_aktif INT,
    jumlah_anggota_baru INT,
    jumlah_anggota_keluar INT,
    jumlah_pengurus INT,
    jumlah_pengawas INT,
    created_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS ON fact_anggota_bulanan (koperasi_id);
CREATE INDEX IF NOT EXISTS ON fact_anggota_bulanan (periode_tahun);

-- Transaction logs for real-time analytics
CREATE TABLE IF NOT EXISTS transaction_logs (
    id UUID,
    tenant_id BIGINT,
    koperasi_id BIGINT,
    transaction_type TEXT,
    transaction_id BIGINT,
    user_id BIGINT,
    amount DECIMAL,
    status TEXT,
    description TEXT,
    metadata TEXT,
    ip_address TEXT,
    user_agent TEXT,
    created_at TIMESTAMP,
    year INT,
    month INT,
    day INT,
    PRIMARY KEY ((koperasi_id, year, month), day, created_at, id)
) WITH CLUSTERING ORDER BY (day ASC, created_at DESC, id ASC);

CREATE INDEX IF NOT EXISTS ON transaction_logs (transaction_type);
CREATE INDEX IF NOT EXISTS ON transaction_logs (status);
CREATE INDEX IF NOT EXISTS ON transaction_logs (user_id);

-- Payment analytics
CREATE TABLE IF NOT EXISTS payment_analytics (
    id UUID,
    tenant_id BIGINT,
    koperasi_id BIGINT,
    payment_provider TEXT,
    payment_method TEXT,
    transaction_type TEXT,
    total_amount DECIMAL,
    admin_fee DECIMAL,
    transaction_count INT,
    success_count INT,
    failed_count INT,
    success_rate DECIMAL,
    created_at TIMESTAMP,
    year INT,
    month INT,
    day INT,
    PRIMARY KEY ((koperasi_id, year, month), payment_provider, day, created_at, id)
) WITH CLUSTERING ORDER BY (payment_provider ASC, day ASC, created_at DESC, id ASC);

CREATE INDEX IF NOT EXISTS ON payment_analytics (payment_method);
CREATE INDEX IF NOT EXISTS ON payment_analytics (transaction_type);

-- PPOB analytics
CREATE TABLE IF NOT EXISTS ppob_analytics (
    id UUID,
    tenant_id BIGINT,
    koperasi_id BIGINT,
    product_category TEXT,
    product_name TEXT,
    provider TEXT,
    transaction_count INT,
    total_revenue DECIMAL,
    total_commission DECIMAL,
    success_count INT,
    failed_count INT,
    success_rate DECIMAL,
    created_at TIMESTAMP,
    year INT,
    month INT,
    day INT,
    PRIMARY KEY ((koperasi_id, year, month), product_category, day, created_at, id)
) WITH CLUSTERING ORDER BY (product_category ASC, day ASC, created_at DESC, id ASC);

CREATE INDEX IF NOT EXISTS ON ppob_analytics (provider);
CREATE INDEX IF NOT EXISTS ON ppob_analytics (product_name);

-- Clinic analytics
CREATE TABLE IF NOT EXISTS klinik_analytics (
    id UUID,
    tenant_id BIGINT,
    koperasi_id BIGINT,
    dokter_id BIGINT,
    dokter_nama TEXT,
    spesialisasi TEXT,
    jumlah_kunjungan INT,
    total_pendapatan DECIMAL,
    rata_rata_biaya_kunjungan DECIMAL,
    created_at TIMESTAMP,
    year INT,
    month INT,
    day INT,
    PRIMARY KEY ((koperasi_id, year, month), dokter_id, day, created_at, id)
) WITH CLUSTERING ORDER BY (dokter_id ASC, day ASC, created_at DESC, id ASC);

CREATE INDEX IF NOT EXISTS ON klinik_analytics (spesialisasi);

-- User activity logs
CREATE TABLE IF NOT EXISTS user_activity_logs (
    id UUID,
    tenant_id BIGINT,
    koperasi_id BIGINT,
    user_id BIGINT,
    username TEXT,
    activity TEXT,
    module TEXT,
    description TEXT,
    ip_address TEXT,
    user_agent TEXT,
    session_id TEXT,
    created_at TIMESTAMP,
    year INT,
    month INT,
    day INT,
    hour INT,
    PRIMARY KEY ((koperasi_id, year, month), day, hour, created_at, id)
) WITH CLUSTERING ORDER BY (day ASC, hour ASC, created_at DESC, id ASC);

CREATE INDEX IF NOT EXISTS ON user_activity_logs (user_id);
CREATE INDEX IF NOT EXISTS ON user_activity_logs (activity);
CREATE INDEX IF NOT EXISTS ON user_activity_logs (module);

-- Error logs
CREATE TABLE IF NOT EXISTS error_logs (
    id UUID,
    tenant_id BIGINT,
    koperasi_id BIGINT,
    user_id BIGINT,
    error_type TEXT,
    error_message TEXT,
    stack_trace TEXT,
    request_data TEXT,
    response_data TEXT,
    module TEXT,
    function TEXT,
    ip_address TEXT,
    user_agent TEXT,
    created_at TIMESTAMP,
    year INT,
    month INT,
    day INT,
    PRIMARY KEY ((koperasi_id, year, month), error_type, day, created_at, id)
) WITH CLUSTERING ORDER BY (error_type ASC, day ASC, created_at DESC, id ASC);

CREATE INDEX IF NOT EXISTS ON error_logs (module);
CREATE INDEX IF NOT EXISTS ON error_logs (user_id);

-- Performance metrics
CREATE TABLE IF NOT EXISTS performance_metrics (
    id UUID,
    tenant_id BIGINT,
    koperasi_id BIGINT,
    metric_name TEXT,
    metric_value DECIMAL,
    metric_unit TEXT,
    endpoint_path TEXT,
    http_method TEXT,
    response_time DECIMAL,
    status_code INT,
    request_size BIGINT,
    response_size BIGINT,
    created_at TIMESTAMP,
    year INT,
    month INT,
    day INT,
    hour INT,
    minute INT,
    PRIMARY KEY ((koperasi_id, year, month), day, hour, minute, created_at, id)
) WITH CLUSTERING ORDER BY (day ASC, hour ASC, minute ASC, created_at DESC, id ASC);

CREATE INDEX IF NOT EXISTS ON performance_metrics (endpoint_path);
CREATE INDEX IF NOT EXISTS ON performance_metrics (http_method);
CREATE INDEX IF NOT EXISTS ON performance_metrics (status_code);