package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"koperasi-merah-putih/internal/models/postgres"

	"gorm.io/gorm"
)

const wilayahCSVURL = "https://raw.githubusercontent.com/kodewilayah/permendagri-72-2019/refs/heads/main/dist/base.csv"

func seedWilayahFromCSV(db *gorm.DB) {
	fmt.Println("Downloading wilayah data from CSV...")

	// Download CSV
	resp, err := http.Get(wilayahCSVURL)
	if err != nil {
		log.Printf("Failed to download CSV: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to download CSV: HTTP %d", resp.StatusCode)
		return
	}

	reader := csv.NewReader(resp.Body)

	// Maps to track IDs
	provinsiMap := make(map[string]uint64)
	kabupatenMap := make(map[string]uint64)
	kecamatanMap := make(map[string]uint64)

	var (
		provinsiCount  int
		kabupatenCount int
		kecamatanCount int
		kelurahanCount int
	)

	fmt.Println("Parsing and inserting wilayah data...")

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading CSV: %v", err)
			continue
		}

		if len(record) < 2 {
			continue
		}

		kode := strings.TrimSpace(record[0])
		nama := strings.TrimSpace(record[1])

		if kode == "" || nama == "" {
			continue
		}

		parts := strings.Split(kode, ".")

		switch len(parts) {
		case 1: // Provinsi (2 digits)
			if len(kode) == 2 {
				id, _ := strconv.ParseUint(kode, 10, 64)
				provinsi := postgres.WilayahProvinsi{
					ID:   id,
					Kode: kode,
					Nama: nama,
				}
				db.FirstOrCreate(&provinsi, postgres.WilayahProvinsi{ID: id})
				provinsiMap[kode] = id
				provinsiCount++
			}

		case 2: // Kabupaten (XX.XX)
			if len(parts[0]) == 2 && len(parts[1]) == 2 {
				fullKode := strings.ReplaceAll(kode, ".", "")
				id, _ := strconv.ParseUint(fullKode, 10, 64)
				provinsiID, _ := strconv.ParseUint(parts[0], 10, 64)

				if _, exists := provinsiMap[parts[0]]; exists {
					kabupaten := postgres.WilayahKabupaten{
						ID:         id,
						Kode:       kode,
						Nama:       nama,
						ProvinsiID: provinsiID,
					}
					db.FirstOrCreate(&kabupaten, postgres.WilayahKabupaten{ID: id})
					kabupatenMap[kode] = id
					kabupatenCount++
				}
			}

		case 3: // Kecamatan (XX.XX.XX)
			if len(parts[0]) == 2 && len(parts[1]) == 2 && len(parts[2]) == 2 {
				fullKode := strings.ReplaceAll(kode, ".", "")
				id, _ := strconv.ParseUint(fullKode, 10, 64)
				kabupatenKode := parts[0] + "." + parts[1]
				kabupatenFullKode := parts[0] + parts[1]
				kabupatenID, _ := strconv.ParseUint(kabupatenFullKode, 10, 64)

				if _, exists := kabupatenMap[kabupatenKode]; exists {
					kecamatan := postgres.WilayahKecamatan{
						ID:          id,
						Kode:        kode,
						Nama:        nama,
						KabupatenID: kabupatenID,
					}
					db.FirstOrCreate(&kecamatan, postgres.WilayahKecamatan{ID: id})
					kecamatanMap[kode] = id
					kecamatanCount++
				}
			}

		case 4: // Kelurahan (XX.XX.XX.XXXX)
			if len(parts[0]) == 2 && len(parts[1]) == 2 && len(parts[2]) == 2 && len(parts[3]) == 4 {
				fullKode := strings.ReplaceAll(kode, ".", "")
				id, _ := strconv.ParseUint(fullKode, 10, 64)
				kecamatanKode := parts[0] + "." + parts[1] + "." + parts[2]
				kecamatanFullKode := parts[0] + parts[1] + parts[2]
				kecamatanID, _ := strconv.ParseUint(kecamatanFullKode, 10, 64)

				if _, exists := kecamatanMap[kecamatanKode]; exists {
					// Determine jenis (kelurahan or desa) based on name
					jenis := "desa"
					namaLower := strings.ToLower(nama)
					if strings.Contains(namaLower, "kelurahan") || strings.Contains(namaLower, "kel.") {
						jenis = "kelurahan"
					}

					kelurahan := postgres.WilayahKelurahan{
						ID:          id,
						Kode:        kode,
						Nama:        nama,
						KecamatanID: kecamatanID,
						Jenis:       jenis,
					}
					db.FirstOrCreate(&kelurahan, postgres.WilayahKelurahan{ID: id})
					kelurahanCount++
				}
			}
		}
	}

	fmt.Printf("✓ Seeded Wilayah from CSV:\n")
	fmt.Printf("  - Provinsi: %d\n", provinsiCount)
	fmt.Printf("  - Kabupaten: %d\n", kabupatenCount)
	fmt.Printf("  - Kecamatan: %d\n", kecamatanCount)
	fmt.Printf("  - Kelurahan: %d\n", kelurahanCount)
}
