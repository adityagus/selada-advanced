package item

// Motorcycle represents the motorcycle table in MySQL POS db
type Motorcycle struct {
	IDMotor    string `gorm:"column:id_motor;primaryKey" json:"id_motor"`
	NamaMotor  string `gorm:"column:nama_motor" json:"nama_motor"`
	IDCategory int    `gorm:"column:id_category" json:"id_category"`
}

func (Motorcycle) TableName() string {
	return "motorcycle"
}

// MotCol represents the motcol table in MySQL POS db
type MotCol struct {
	IDMotor      string `gorm:"column:id_motor" json:"id_motor"`
	IDIcon       string `gorm:"column:id_icon" json:"id_icon"`
	GambarColor  string `gorm:"column:gambar_color" json:"gambar_color"`
	DefaultColor string `gorm:"column:default_color" json:"default_color"`
}

func (MotCol) TableName() string {
	return "motcol"
}

// MotFitur represents the motfitur table in MySQL POS db
type MotFitur struct {
	IDMotor string  `gorm:"column:id_motor;primaryKey" json:"id_motor"`
	Harga   float64 `gorm:"column:harga" json:"harga"`
	Detail  string  `gorm:"column:detail" json:"detail"`
}

func (MotFitur) TableName() string {
	return "motfitur"
}

// MotFiturList represents the motfitur_list table in MySQL POS db
type MotFiturList struct {
	ID      int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDMotor string `gorm:"column:id_motor" json:"id_motor"`
	Fitur   string `gorm:"column:fitur" json:"fitur"`
}

func (MotFiturList) TableName() string {
	return "motfitur_list"
}

// FiturMotor represents the fiturmotor table in MySQL POS db
type FiturMotor struct {
	IDMotor string `gorm:"column:id_motor;primaryKey" json:"id_motor"`
	Brosur  string `gorm:"column:brosur" json:"brosur"`
}

func (FiturMotor) TableName() string {
	return "fiturmotor"
}
