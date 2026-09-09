package models

import "time"

// Rencana represents the tb_rencana table in MySQL default db
type Rencana struct {
	IDRencana    int    `gorm:"column:id_rencana;primaryKey;autoIncrement"`
	IDCustomer   int    `gorm:"column:id_customer"`
	IDLogin      string `gorm:"column:id_login"`
	IDSumbercust string `gorm:"column:id_sumbercust"`
}

// TableName overrides the default table name for Rencana
func (Rencana) TableName() string {
	return "tb_rencana"
}

// RencanaDet represents the tb_rencana_det table in MySQL default db
type RencanaDet struct {
	IDRencanaDet int       `gorm:"column:id_rencana_det;primaryKey;autoIncrement"`
	IDRencana    int       `gorm:"column:id_rencana"`
	TglRencana   time.Time `gorm:"column:tgl_rencana"`
	KetRencana   string    `gorm:"column:ket_rencana"`
	IDHslaktiv   int       `gorm:"column:id_hslaktiv"`
	Status       int       `gorm:"column:status"`
}

type ApplicationSource struct {
	Code string `json:"code" gorm:"column:code"`
	Name string `json:"name" gorm:"column:name"`
}

// TableName overrides the default table name for RencanaDet
func (RencanaDet) TableName() string {
	return "tb_rencana_det"
}

// HslAktivitas represents the tb_hslaktivitas table in MySQL default db
type HslAktivitas struct {
	IDHslaktiv   int    `gorm:"column:id_hslaktiv;primaryKey"`
	NamaHslaktiv string `gorm:"column:nama_hslaktiv"`
}

// TableName overrides the default table name for HslAktivitas
func (HslAktivitas) TableName() string {
	return "tb_hslaktivitas"
}

// CronDetPointFn represents the tb_crondetpointfn table in MySQL default db
type CronDetPointFn struct {
	ID        int    `gorm:"column:id;primaryKey;autoIncrement"`
	Bln       int    `gorm:"column:bln"`
	Thn       int    `gorm:"column:thn"`
	IDSales   string `gorm:"column:id_sales"`
	Pleads    int    `gorm:"column:pleads"`
	Pfu       int    `gorm:"column:pfu"`
	Pdo       int    `gorm:"column:pdo"`
	Pakurdat  int    `gorm:"column:pakurdat"`
	Leadstodo int    `gorm:"column:leadstodo"`
	Leadspal  int    `gorm:"column:leadspal"`
}

// TableName overrides the default table name for CronDetPointFn
func (CronDetPointFn) TableName() string {
	return "tb_crondetpointfn"
}
