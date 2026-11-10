-- ==============================================================================
-- MIGRATION: Table & Cronjob Reporting Paku Bulanan (Snapshot Akhir Bulan)
-- Target Database: PostgreSQL selada_pg
-- Tujuan: Membekukan (paku) data laporan performance & realisasi sales setiap akhir bulan
--         agar angka historis tidak berubah akibat perubahan di masa mendatang.
-- ==============================================================================

-- 1. Buat Tabel tb_report_paku_bulanan
CREATE TABLE IF NOT EXISTS tb_report_paku_bulanan (
    id SERIAL PRIMARY KEY,
    periode_tahun INT NOT NULL,                          -- Contoh: 2026
    periode_bulan INT NOT NULL,                          -- Contoh: 10
    periode_kode VARCHAR(7) NOT NULL,                    -- Format: 'YYYY-MM', misal '2026-10'
    id_cabang INT DEFAULT 0,                             -- ID Cabang
    nm_cabang VARCHAR(100) DEFAULT '',                   -- Nama Cabang
    id_sales VARCHAR(100) NOT NULL,                      -- ID / Username Sales
    nama_sales VARCHAR(150) DEFAULT '',                  -- Nama Sales Person
    total_leads INT DEFAULT 0,                           -- Total customer status Leads
    total_prospect INT DEFAULT 0,                        -- Total status Prospect
    total_hot_prospect INT DEFAULT 0,                    -- Total status Hot Prospect
    total_sbg_deal INT DEFAULT 0,                        -- Total deal pencairan SBG
    total_batal INT DEFAULT 0,                           -- Total status Batal
    total_booking INT DEFAULT 0,                         -- Total aktivitas rencana/booking
    estimasi_up_cair NUMERIC(18,2) DEFAULT 0,            -- Estimasi Uang Pinjaman terkumpul
    conversion_rate NUMERIC(5,2) DEFAULT 0,              -- Tingkat konversi (%)
    is_locked BOOLEAN DEFAULT TRUE,                      -- Status paku (terkunci)
    snapshot_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_report_paku_periode_sales UNIQUE (periode_tahun, periode_bulan, id_sales, id_cabang)
);

-- Buat Index untuk pencarian cepat
CREATE INDEX IF NOT EXISTS idx_paku_periode ON tb_report_paku_bulanan (periode_kode);
CREATE INDEX IF NOT EXISTS idx_paku_sales ON tb_report_paku_bulanan (id_sales);
CREATE INDEX IF NOT EXISTS idx_paku_cabang ON tb_report_paku_bulanan (id_cabang);

-- ==============================================================================
-- 2. Stored Procedure / Function untuk Mengisi Data Paku Otomatis
-- ==============================================================================
CREATE OR REPLACE FUNCTION sp_generate_report_paku_bulanan(p_tahun INT, p_bulan INT)
RETURNS void AS $$
DECLARE
    v_periode_kode VARCHAR(7);
BEGIN
    v_periode_kode := p_tahun || '-' || LPAD(p_bulan::TEXT, 2, '0');

    -- Insert snapshot dari tabel tb_rencana & tb_rencana_det
    INSERT INTO tb_report_paku_bulanan (
        periode_tahun,
        periode_bulan,
        periode_kode,
        id_cabang,
        nm_cabang,
        id_sales,
        nama_sales,
        total_leads,
        total_prospect,
        total_hot_prospect,
        total_sbg_deal,
        total_batal,
        total_booking,
        estimasi_up_cair,
        conversion_rate,
        is_locked,
        snapshot_at
    )
    SELECT
        p_tahun AS periode_tahun,
        p_bulan AS periode_bulan,
        v_periode_kode AS periode_kode,
        COALESCE(c.id_cabang, 0) AS id_cabang,
        COALESCE(c.nm_cabang, '-') AS nm_cabang,
        r.id_login AS id_sales,
        COALESCE(l.nama, r.id_login) AS nama_sales,
        COUNT(CASE WHEN LOWER(h.nama_hslaktiv) LIKE '%lead%' THEN 1 END) AS total_leads,
        COUNT(CASE WHEN LOWER(h.nama_hslaktiv) LIKE '%prospect%' AND LOWER(h.nama_hslaktiv) NOT LIKE '%hot%' THEN 1 END) AS total_prospect,
        COUNT(CASE WHEN LOWER(h.nama_hslaktiv) LIKE '%hot%' THEN 1 END) AS total_hot_prospect,
        COUNT(CASE WHEN LOWER(h.nama_hslaktiv) LIKE '%sbg%' THEN 1 END) AS total_sbg_deal,
        COUNT(CASE WHEN LOWER(h.nama_hslaktiv) LIKE '%batal%' THEN 1 END) AS total_batal,
        COUNT(rd.id_rencana_det) AS total_booking,
        (COUNT(CASE WHEN LOWER(h.nama_hslaktiv) LIKE '%sbg%' THEN 1 END) * 5500000.00) AS estimasi_up_cair,
        CASE
            WHEN COUNT(rd.id_rencana_det) > 0 THEN
                ROUND((COUNT(CASE WHEN LOWER(h.nama_hslaktiv) LIKE '%sbg%' THEN 1 END)::NUMERIC / COUNT(rd.id_rencana_det)::NUMERIC) * 100, 2)
            ELSE 0.00
        END AS conversion_rate,
        TRUE AS is_locked,
        CURRENT_TIMESTAMP AS snapshot_at
    FROM tb_rencana r
    JOIN tb_rencana_det rd ON rd.id_rencana = r.id_rencana
    LEFT JOIN tb_hslaktivitas h ON h.id_hslaktiv = rd.id_hslaktiv
    LEFT JOIN tb_login l ON l.id_login = r.id_login
    LEFT JOIN (
        -- Referensi mapping cabang jika ada di tb_customer atau profile
        SELECT DISTINCT id_login, 42 AS id_cabang, 'SAM Head Office' AS nm_cabang FROM tb_customer
    ) c ON c.id_login = r.id_login
    WHERE EXTRACT(YEAR FROM rd.tgl_rencana) = p_tahun
      AND EXTRACT(MONTH FROM rd.tgl_rencana) = p_bulan
    GROUP BY c.id_cabang, c.nm_cabang, r.id_login, l.nama
    ON CONFLICT (periode_tahun, periode_bulan, id_sales, id_cabang)
    DO UPDATE SET
        total_leads = EXCLUDED.total_leads,
        total_prospect = EXCLUDED.total_prospect,
        total_hot_prospect = EXCLUDED.total_hot_prospect,
        total_sbg_deal = EXCLUDED.total_sbg_deal,
        total_batal = EXCLUDED.total_batal,
        total_booking = EXCLUDED.total_booking,
        estimasi_up_cair = EXCLUDED.estimasi_up_cair,
        conversion_rate = EXCLUDED.conversion_rate,
        is_locked = TRUE,
        snapshot_at = CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql;

-- ==============================================================================
-- 3. FORMAT CRONJOB (Setiap Akhir Bulan Pukul 23:59)
-- ==============================================================================
-- Format Standar Linux Crontab:
-- Jalankan command ini untuk membuka crontab: `crontab -e`
--
-- 59 23 28-31 * * [ $(date -d tomorrow +\%d) -eq 01 ] && PGPASSWORD="stagingdev@6177" psql -h 10.10.10.8 -p 5433 -U ervan -d selada_pg -c "SELECT sp_generate_report_paku_bulanan(EXTRACT(YEAR FROM CURRENT_DATE)::INT, EXTRACT(MONTH FROM CURRENT_DATE)::INT);"
--
-- Penjelasan Format Cron:
-- • '59 23 28-31 * *' : Berjalan pada jam 23:59 pada tanggal 28 hingga 31 setiap bulan.
-- • '[ $(date -d tomorrow +\%d) -eq 01 ]' : Memastikan HANYA berjalan jika besok adalah tanggal 1 (yaitu benar-benar hari terakhir di bulan tsb).
-- • Memanggil function sp_generate_report_paku_bulanan() dengan parameter tahun & bulan berjalan.
-- ==============================================================================
