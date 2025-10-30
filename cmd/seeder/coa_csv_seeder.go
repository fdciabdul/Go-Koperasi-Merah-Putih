package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"koperasi-merah-putih/internal/models/postgres"

	"gorm.io/gorm"
)

func seedCOAFromCSV(db *gorm.DB) {
	fmt.Println("Loading COA from CSV...")

	// Get the first Koperasi ID
	var koperasi postgres.Koperasi
	if err := db.First(&koperasi).Error; err != nil {
		log.Printf("ERROR: No koperasi found in COA seeder: %v", err)
		// Let's check if any koperasi exist at all
		var count int64
		db.Model(&postgres.Koperasi{}).Count(&count)
		log.Printf("DEBUG: Total koperasi count: %d", count)
		return
	}
	fmt.Printf("DEBUG: Using Koperasi ID=%d, TenantID=%d\n", koperasi.ID, koperasi.TenantID)

	file, err := os.Open("cmd/seeder/csv/coa_koperasi.csv")
	if err != nil {
		log.Printf("Failed to open COA CSV: %v", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Failed to read COA CSV: %v", err)
		return
	}

	// Skip header row
	records = records[1:]

	// First, ensure COA categories exist
	kategoris := map[string]uint64{
		"ASET":       1,
		"KEWAJIBAN":  2,
		"EKUITAS":    3,
		"PENDAPATAN": 4,
		"BEBAN":      5,
	}

	kategoriData := []postgres.COAKategori{
		{ID: 1, Kode: "1", Nama: "ASET", Tipe: "aset", Urutan: 1},
		{ID: 2, Kode: "2", Nama: "KEWAJIBAN", Tipe: "kewajiban", Urutan: 2},
		{ID: 3, Kode: "3", Nama: "EKUITAS", Tipe: "ekuitas", Urutan: 3},
		{ID: 4, Kode: "4", Nama: "PENDAPATAN", Tipe: "pendapatan", Urutan: 4},
		{ID: 5, Kode: "5", Nama: "BEBAN", Tipe: "beban", Urutan: 5},
	}

	for _, kat := range kategoriData {
		db.FirstOrCreate(&kat, postgres.COAKategori{ID: kat.ID})
	}

	// Temporarily drop the self-referential foreign key constraint for seeding
	db.Exec("ALTER TABLE coa_akuns DROP CONSTRAINT IF EXISTS fk_coa_akuns_children")

	count := 0
	for _, record := range records {
		if len(record) < 6 {
			continue
		}

		kodeAkun := strings.TrimSpace(record[0])
		namaAkun := strings.TrimSpace(record[1])
		kategori := strings.TrimSpace(record[2])
		subKategori := strings.TrimSpace(record[3])
		tipe := strings.TrimSpace(record[4])
		deskripsi := strings.TrimSpace(record[5])

		if kodeAkun == "" || namaAkun == "" {
			continue
		}

		// Get kategori ID
		kategoriID, exists := kategoris[kategori]
		if !exists {
			continue
		}

		// Determine saldo normal based on tipe
		saldoNormal := "debit"
		if tipe == "Kredit" {
			saldoNormal = "kredit"
		}

		// Determine level based on kode_akun pattern
		levelAkun := 1
		parts := strings.Split(kodeAkun, "-")
		if len(parts) == 2 {
			if strings.HasSuffix(parts[1], "000") {
				levelAkun = 1 // Header
			} else if len(parts[1]) == 4 {
				levelAkun = 2 // Sub-account
			}
		}

		// Store deskripsi in a comment or just skip it (model doesn't have this field)
		_ = deskripsi
		_ = subKategori

		akun := postgres.COAAkun{
			TenantID:    koperasi.TenantID,
			KoperasiID:  koperasi.ID,
			KodeAkun:    kodeAkun,
			NamaAkun:    namaAkun,
			KategoriID:  kategoriID,
			ParentID:    0, // No parent for now
			LevelAkun:   levelAkun,
			SaldoNormal: saldoNormal,
			IsKas:       strings.Contains(strings.ToLower(namaAkun), "kas") || strings.Contains(strings.ToLower(namaAkun), "bank"),
			IsAktif:     true,
		}

		result := db.FirstOrCreate(&akun, postgres.COAAkun{KodeAkun: kodeAkun, KoperasiID: koperasi.ID})
		if result.Error != nil {
			if count == 0 {
				log.Printf("ERROR: Failed to create first COA account '%s': %v", namaAkun, result.Error)
			}
		} else {
			if count == 0 {
				fmt.Printf("DEBUG: Successfully created first COA account ID=%d, Name='%s'\n", akun.ID, akun.NamaAkun)
			}
		}
		count++
	}

	// Verify records were created
	var verifyCount int64
	db.Model(&postgres.COAAkun{}).Where("koperasi_id = ?", koperasi.ID).Count(&verifyCount)
	fmt.Printf("✓ Seeded %d COA Akun from CSV (Verified in DB: %d)\n", count, verifyCount)
}

func seedJurnalFromCSV(db *gorm.DB) {
	fmt.Println("Loading Jurnal from CSV...")

	// Get the first Koperasi ID
	var koperasi postgres.Koperasi
	if err := db.First(&koperasi).Error; err != nil {
		log.Printf("No koperasi found, skipping journal seeding: %v", err)
		return
	}

	file, err := os.Open("cmd/seeder/csv/jurnal_detail_koperasi.csv")
	if err != nil {
		log.Printf("Failed to open Jurnal CSV: %v", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Failed to read Jurnal CSV: %v", err)
		return
	}

	// Skip header row
	records = records[1:]

	// Build a map of COA accounts by kode_akun AND by name
	var coaAkuns []postgres.COAAkun
	db.Where("koperasi_id = ?", koperasi.ID).Find(&coaAkuns)

	coaMapByKode := make(map[string]uint64)
	coaMapByName := make(map[string]uint64)
	for _, akun := range coaAkuns {
		coaMapByKode[akun.KodeAkun] = akun.ID
		coaMapByName[akun.NamaAkun] = akun.ID
	}

	fmt.Printf("DEBUG: Loaded %d COA accounts for mapping\n", len(coaAkuns))
	if len(coaAkuns) > 0 {
		fmt.Printf("DEBUG: Sample account name: '%s'\n", coaAkuns[0].NamaAkun)
	}

	// Group records by No_Bukti to create JurnalUmum entries
	jurnalMap := make(map[string][]map[string]string)

	for _, record := range records {
		if len(record) < 7 {
			continue
		}

		tanggal := strings.TrimSpace(record[0])
		noBukti := strings.TrimSpace(record[1])
		keterangan := strings.TrimSpace(record[2])
		namaAkun := strings.TrimSpace(record[3]) // CSV "Kode_Akun" is actually the account name
		detailKeterangan := strings.TrimSpace(record[4]) // CSV "Nama_Akun" is actually the detail description
		debitStr := strings.TrimSpace(record[5])
		kreditStr := strings.TrimSpace(record[6])

		if noBukti == "" {
			continue
		}

		debit, _ := strconv.ParseFloat(debitStr, 64)
		kredit, _ := strconv.ParseFloat(kreditStr, 64)

		entry := map[string]string{
			"tanggal":          tanggal,
			"noBukti":          noBukti,
			"keterangan":       keterangan,
			"namaAkun":         namaAkun,
			"detailKeterangan": detailKeterangan,
			"debit":            fmt.Sprintf("%.0f", debit),
			"kredit":           fmt.Sprintf("%.0f", kredit),
		}

		jurnalMap[noBukti] = append(jurnalMap[noBukti], entry)
	}

	// Create JurnalUmum and JurnalDetail entries
	jurnalCount := 0
	detailCount := 0

	for noBukti, entries := range jurnalMap {
		if len(entries) == 0 {
			continue
		}

		// Get the date and keterangan from first entry
		firstEntry := entries[0]
		tanggalStr := firstEntry["tanggal"]

		// Parse date (format: 2025-01-02)
		tanggal, err := time.Parse("2006-01-02", tanggalStr)
		if err != nil {
			log.Printf("Failed to parse date %s: %v", tanggalStr, err)
			continue
		}

		// Calculate total debit and kredit
		var totalDebit, totalKredit float64
		var keterangan string
		for _, entry := range entries {
			debit, _ := strconv.ParseFloat(entry["debit"], 64)
			kredit, _ := strconv.ParseFloat(entry["kredit"], 64)
			totalDebit += debit
			totalKredit += kredit
			if keterangan == "" {
				keterangan = entry["keterangan"]
			}
		}

		// Create JurnalUmum
		nomorJurnal := fmt.Sprintf("JU%s%08d", tanggal.Format("20060102"), jurnalCount+1)

		jurnal := postgres.JurnalUmum{
			TenantID:         koperasi.TenantID,
			KoperasiID:       koperasi.ID,
			NomorJurnal:      nomorJurnal,
			TanggalTransaksi: tanggal,
			Referensi:        noBukti,
			Keterangan:       keterangan,
			TotalDebit:       totalDebit,
			TotalKredit:      totalKredit,
			Status:           "posted",
			CreatedBy:        1,
			PostedBy:         1,
			PostedAt:         &tanggal,
		}

		result := db.FirstOrCreate(&jurnal, postgres.JurnalUmum{Referensi: noBukti, KoperasiID: koperasi.ID})
		if result.Error != nil {
			log.Printf("Failed to create jurnal %s: %v", noBukti, result.Error)
			continue
		}

		jurnalCount++

		// Create JurnalDetail entries
		detailCreatedCount := 0
		for i, entry := range entries {
			namaAkun := entry["namaAkun"]
			detailKeterangan := entry["detailKeterangan"]

			if i == 0 && jurnalCount == 1 {
				fmt.Printf("DEBUG: First journal entry - looking for account: '%s'\n", namaAkun)
			}

			// Find account by name
			akunID, exists := coaMapByName[namaAkun]
			if !exists {
				if detailCreatedCount == 0 && i == 0 {
					fmt.Printf("DEBUG: Account not found. Sample from map: ")
					count := 0
					for name := range coaMapByName {
						if count < 3 {
							fmt.Printf("'%s', ", name)
							count++
						}
					}
					fmt.Println()
				}
				log.Printf("COA account not found for name: %s", namaAkun)
				continue
			}
			detailCreatedCount++

			debit, _ := strconv.ParseFloat(entry["debit"], 64)
			kredit, _ := strconv.ParseFloat(entry["kredit"], 64)

			detail := postgres.JurnalDetail{
				JurnalID:   jurnal.ID,
				AkunID:     akunID,
				Keterangan: detailKeterangan,
				Debit:      debit,
				Kredit:     kredit,
			}

			db.FirstOrCreate(&detail, postgres.JurnalDetail{JurnalID: jurnal.ID, AkunID: akunID})
			detailCount++
		}
	}

	fmt.Printf("✓ Seeded %d Jurnal Umum and %d Jurnal Details from CSV\n", jurnalCount, detailCount)
}
