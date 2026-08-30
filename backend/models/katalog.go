package models

// Motorcycle represents the motorcycle table in MySQL POS db
type Motorcycle struct {
	IDMotor    string `gorm:"column:id_motor;primaryKey"`
	NamaMotor  string `gorm:"column:nama_motor"`
	IDCategory int    `gorm:"column:id_category"`
}

// TableName overrides the default table name for Motorcycle
func (Motorcycle) TableName() string {
	return "motorcycle"
}

// MotCol represents the motcol table in MySQL POS db
type MotCol struct {
	IDMotor      string `gorm:"column:id_motor"`
	IDIcon       string `gorm:"column:id_icon"`
	GambarColor  string `gorm:"column:gambar_color"`
	DefaultColor string `gorm:"column:default_color"`
}

// TableName overrides the default table name for MotCol
func (MotCol) TableName() string {
	return "motcol"
}

// MasterColor represents the master_color table in MySQL POS db
type MasterColor struct {
	IDMaster  string `gorm:"column:id_master;primaryKey"`
	ColorName string `gorm:"column:color_name"`
}

// TableName overrides the default table name for MasterColor
func (MasterColor) TableName() string {
	return "master_color"
}

// MotFitur represents the motfitur table in MySQL POS db
type MotFitur struct {
	IDMotor string  `gorm:"column:id_motor;primaryKey"`
	Harga   float64 `gorm:"column:harga"`
	Detail  string  `gorm:"column:detail"`
}

// TableName overrides the default table name for MotFitur
func (MotFitur) TableName() string {
	return "motfitur"
}

// MotFiturList represents the motfitur_list table in MySQL POS db
type MotFiturList struct {
	ID      int    `gorm:"column:id;primaryKey;autoIncrement"`
	IDMotor string `gorm:"column:id_motor"`
	Fitur   string `gorm:"column:fitur"`
}

// TableName overrides the default table name for MotFiturList
func (MotFiturList) TableName() string {
	return "motfitur_list"
}

// FiturMotor represents the fiturmotor table in MySQL POS db
type FiturMotor struct {
	IDMotor string `gorm:"column:id_motor;primaryKey"`
	Brosur  string `gorm:"column:brosur"`
}

// TableName overrides the default table name for FiturMotor
func (FiturMotor) TableName() string {
	return "fiturmotor"
}
