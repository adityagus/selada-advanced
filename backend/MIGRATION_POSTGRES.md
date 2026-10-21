# Migrasi Selada ke PostgreSQL

Backend Selada menggunakan `config.DBSelada` dengan konfigurasi `DB_PG2_*`.
`DB_PG_*` tetap untuk sam_live dan `DB_MYSQL_POS_*` tetap untuk sam_pos.
`DB_MYSQL_*` hanya digunakan sebagai sumber migrasi Selada.

## Hasil eksekusi 20 September 2026

Migrasi berhasil commit: **23 tabel, 327.371 baris**. Jumlah baris setiap
tabel tujuan sama dengan snapshot sumber. MySQL sumber tidak diubah.

| Tabel berisi data | Baris |
| --- | ---: |
| log_daily_incentive | 1.176 |
| tb_aktivitas | 17 |
| tb_cronkpi | 19.198 |
| tb_customer | 79.987 |
| tb_grupksg | 7 |
| tb_hari | 46 |
| tb_hslaktivitas | 7 |
| tb_img | 756 |
| tb_katcust | 2 |
| tb_rencana | 88.891 |
| tb_rencana_det | 137.263 |
| tb_sumbercust | 9 |
| tb_target | 12 |

Sepuluh tabel lainnya kosong dan tetap dibuat. Primary key, indeks BTREE,
unique index, dan foreign key sumber disalin; sequence ID disesuaikan agar
insert berikutnya tidak bertabrakan dengan data lama.

## Menjalankan pada tujuan baru

Jalankan dari direktori `backend`. Isi `.env` mengikuti `.env-example`:
`DB_MYSQL_*` untuk sumber dan `DB_PG2_*` untuk tujuan PostgreSQL.
Pengaturan environment proses memiliki prioritas atas `.env`.

```powershell
# Pemeriksaan sumber dan tujuan, tanpa menulis data
go run ./cmd/migrate-selada

# Buat database tujuan bila belum ada, lalu salin schema dan data
go run ./cmd/migrate-selada --apply
```

Role tujuan perlu hak CREATE DATABASE jika database belum dibuat. Database
yang sudah ada harus memiliki schema public kosong. Skrip menolak menimpa
tabel yang sudah ada; menjalankan ulang pada tujuan yang sudah dimigrasikan
akan berhenti tanpa menghapus data.

Semua pembuatan tabel dan penyalinan data berada dalam satu transaksi
PostgreSQL. Kegagalan membatalkan transaksi tersebut. CREATE DATABASE berada
di luar transaksi, sehingga database kosong dapat tersisa setelah kegagalan.

Sumber dibaca dalam transaksi repeatable-read. Untuk cutover tanpa kehilangan
perubahan, hentikan penulisan aplikasi lama sebelum snapshot diambil, lalu
jalankan backend versi ini. Migrasi ini bukan sinkronisasi berkelanjutan:
penulisan baru ke MySQL setelah snapshot tidak otomatis masuk PostgreSQL.

## Validasi

```powershell
go test ./...
$env:SELADA_INTEGRATION = '1'
go test ./handlers -run TestSeladaPostgresIntegration -count=1 -v
Remove-Item Env:SELADA_INTEGRATION
```

Tes integrasi menggunakan database tujuan dari `.env`, menguji insert
customer/rencana/detail beserta pengembalian ID, pencarian customer dan
inquiry, riwayat customer, counts, konteks inquiry, dan penyimpanan SBG.
Semua record pengujian di-rollback; sequence PostgreSQL dapat tetap maju.

Query aplikasi telah disesuaikan: pencarian menggunakan ILIKE, timestamp
diubah menjadi text sebelum COALESCE dengan string kosong, cast ID memakai
TEXT, dan SET SESSION sql_mode MySQL dihapus. Riwayat SBG membaca kolom
`tb_rencana_det.no_do` yang tersedia pada schema sumber.

## Batasan

- Tipe/objek yang belum didukung, trigger, routine, non-InnoDB, prefix index,
  dan tanggal tahun nol membuat pemeriksaan berhenti; tidak dikonversi diam-diam.
- Default string/numerik wajib tanpa default diberi `''`/`0` untuk kompatibilitas
  dengan insert aplikasi yang sebelumnya memakai mode MySQL permisif.
- Collation MySQL tidak disalin. ILIKE mempertahankan pencarian tanpa membedakan
  huruf besar/kecil; kesetaraan, pengurutan, dan unique text mengikuti PostgreSQL.
- Tabel `tb_crondetpointfn` yang dirujuk fitur poin katalog memang tidak ada di
  sumber yang diperiksa; migrasi tidak membuat atau mengarang data tabel ini.
- Restart backend yang masih memakai binary lama agar koneksi Selada beralih.

Referensi tipe identity: https://www.postgresql.org/docs/18/ddl-identity-columns.html
