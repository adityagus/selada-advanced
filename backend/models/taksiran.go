package models

// JenisBarang represents tbljenis_barang in PostgreSQL
type JenisBarang struct {
	KdJenisBarang string `gorm:"column:kd_jenis_barang;primaryKey"`
	NmJenisBarang string `gorm:"column:nm_jenis_barang"`
}

// TableName overrides the default table name for JenisBarang
func (JenisBarang) TableName() string {
	return "tbljenis_barang"
}

// Barang represents tblbarang in PostgreSQL
type Barang struct {
	KdBarang     string `gorm:"column:kd_barang;primaryKey"`
	NmBarang     string `gorm:"column:nm_barang"`
	BarangActive string `gorm:"column:barang_active"` // 'true' / 'false'
	FkJenisBarang string `gorm:"column:fk_jenis_barang"`
	FkGrade      string `gorm:"column:fk_grade"`
	StatusAntam  string `gorm:"column:status_antam"`
}

// TableName overrides the default table name for Barang
func (Barang) TableName() string {
	return "tblbarang"
}

// BarangDetail represents tblbarang_detail in PostgreSQL
type BarangDetail struct {
	FkBarang  string  `gorm:"column:fk_barang"`
	FkWilayah string  `gorm:"column:fk_wilayah"`
	HargaTaks float64 `gorm:"column:harga_taks"` // example detail field
}

// TableName overrides the default table name for BarangDetail
func (BarangDetail) TableName() string {
	return "tblbarang_detail"
}

// Cabang represents tblcabang in PostgreSQL
type Cabang struct {
	KdCabang           string  `gorm:"column:kd_cabang;primaryKey"`
	HargaStleAntam     float64 `gorm:"column:harga_stle_antam"`
	HargaStleNonAntam  float64 `gorm:"column:harga_stle_non_antam"`
	HargaStlePerhiasan float64 `gorm:"column:harga_stle_perhiasan"`
}

// TableName overrides the default table name for Cabang
func (Cabang) TableName() string {
	return "tblcabang"
}

// Produk represents tblproduk in PostgreSQL
type Produk struct {
	KdProduk string `gorm:"column:kd_produk;primaryKey"`
	NmProduk string `gorm:"column:nm_produk"`
	Active   bool   `gorm:"column:active"`
}

// TableName overrides the default table name for Produk
func (Produk) TableName() string {
	return "tblproduk"
}

// ProdukDetailBiayaAdmin represents tblproduk_detail_biaya_admin in PostgreSQL
type ProdukDetailBiayaAdmin struct {
	FkProduk string  `gorm:"column:fk_produk"`
	Dari     float64 `gorm:"column:dari"`
	Ke       float64 `gorm:"column:ke"`
	Nilai    float64 `gorm:"column:nilai"`
}

// TableName overrides the default table name for ProdukDetailBiayaAdmin
func (ProdukDetailBiayaAdmin) TableName() string {
	return "tblproduk_detail_biaya_admin"
}
