package database

import (
	"fmt"
	"log"

	"koperasi-merah-putih/config"
	"koperasi-merah-putih/internal/models/postgres"
)

type DatabaseManager struct {
	Postgres  *PostgresDB
	Cassandra *CassandraDB
}

func NewDatabaseManager(cfg *config.Config) (*DatabaseManager, error) {
	postgresDB, err := NewPostgresConnection(&cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	cassandraDB, err := NewCassandraConnection(&cfg.Cassandra)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Cassandra: %v", err)
	}

	return &DatabaseManager{
		Postgres:  postgresDB,
		Cassandra: cassandraDB,
	}, nil
}

func (dm *DatabaseManager) AutoMigrate() error {
	err := dm.Postgres.DB.AutoMigrate(
		&postgres.Tenant{},
		&postgres.WilayahProvinsi{},
		&postgres.WilayahKabupaten{},
		&postgres.WilayahKecamatan{},
		&postgres.WilayahKelurahan{},
		&postgres.JenisKoperasi{},
		&postgres.BentukKoperasi{},
		&postgres.StatusKoperasi{},
		&postgres.KBLI{},
		&postgres.Koperasi{},
		&postgres.KoperasiAktivitasUsaha{},
		&postgres.JabatanKoperasi{},
		&postgres.AnggotaKoperasi{},
		&postgres.COAKategori{},
		&postgres.COAAkun{},
		&postgres.ModalKoperasi{},
		&postgres.JurnalUmum{},
		&postgres.JurnalDetail{},
		&postgres.ProdukSimpanPinjam{},
		&postgres.RekeningSimpanPinjam{},
		&postgres.TransaksiSimpanPinjam{},
		&postgres.PPOBKategori{},
		&postgres.PPOBProvider{},
		&postgres.PPOBProduk{},
		&postgres.PPOBTransaksi{},
		&postgres.PPOBPaymentConfig{},
		&postgres.PPOBSettlement{},
		&postgres.PPOBSettlementDetail{},
		&postgres.PaymentProvider{},
		&postgres.PaymentMethod{},
		&postgres.PaymentTransaction{},
		&postgres.PaymentCallback{},
		&postgres.User{},
		&postgres.Permission{},
		&postgres.RolePermission{},
		&postgres.UserRegistration{},
		&postgres.UserRegistrationLog{},
		&postgres.SimpananPokokConfig{},
		&postgres.SimpananPokokTransaksi{},
		&postgres.KlinikPasien{},
		&postgres.KlinikTenagaMedis{},
		&postgres.KlinikKunjungan{},
		&postgres.KlinikObat{},
		&postgres.KlinikResep{},
		&postgres.AuditLog{},
		&postgres.SystemSetting{},
		&postgres.SequenceNumber{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate PostgreSQL models: %v", err)
	}

	log.Println("PostgreSQL auto-migration completed successfully")
	return nil
}

func (dm *DatabaseManager) Close() {
	if dm.Cassandra != nil {
		dm.Cassandra.Close()
	}

	if dm.Postgres != nil && dm.Postgres.DB != nil {
		sqlDB, err := dm.Postgres.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}