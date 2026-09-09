package models

// Customer represents the tb_customer table in MySQL default db
type Customer struct {
	IDCustomer      int     `gorm:"column:id_customer;primaryKey;autoIncrement" json:"id_customer"`
	Nama            string  `gorm:"column:nama" json:"nama"`
	Hp              string  `gorm:"column:hp" json:"hp"`
	Alamat          string  `gorm:"column:alamat" json:"alamat"`
	IDLogin         string  `gorm:"column:id_login" json:"id_login"`
	Nickname        string  `gorm:"column:nickname" json:"nickname"`
	TmpLahir        string  `gorm:"column:tmp_lahir" json:"tmp_lahir"`
	Identitas       string  `gorm:"column:identitas" json:"identitas"`
	JenisKelamin    string  `gorm:"column:jenis_kelamin" json:"jenis_kelamin,omitempty"`
	Kota            string  `gorm:"column:kota" json:"kota,omitempty"`
	KodePos         string  `gorm:"column:kode_pos" json:"kode_pos,omitempty"`
	Telepon1        string  `gorm:"column:telepon_1" json:"telepon_1,omitempty"`
	Telepon2        string  `gorm:"column:telepon_2" json:"telepon_2,omitempty"`
	Fax             string  `gorm:"column:fax" json:"fax,omitempty"`
	Email           string  `gorm:"column:email" json:"email,omitempty"`
	NPWP            string  `gorm:"column:npwp" json:"npwp,omitempty"`
	NamaPajak       string  `gorm:"column:nama_pajak" json:"nama_pajak,omitempty"`
	AlamatPajak     string  `gorm:"column:alamat_pajak" json:"alamat_pajak,omitempty"`
	KotaPajak       string  `gorm:"column:kota_pajak" json:"kota_pajak,omitempty"`
	TaxPerson       int     `gorm:"column:tax_person" json:"tax_person,omitempty"`
	SalesPerson     string  `gorm:"column:sales_person" json:"sales_person,omitempty"`
	Plafon          float64 `gorm:"column:plafon" json:"plafon,omitempty"`
	IDBranch        int     `gorm:"column:id_branch" json:"id_branch,omitempty"`
	Overdue1        int     `gorm:"-" json:"overdue1,omitempty"`
	Overdue2        int     `gorm:"-" json:"overdue2,omitempty"`
	BillingDay      int     `gorm:"-" json:"billing_day,omitempty"`
	CollectionDay   int     `gorm:"-" json:"collection_day,omitempty"`
	NamaCusgroup    string  `gorm:"-" json:"nama_cusgroup,omitempty"`
	NamaCussubgroup string  `gorm:"-" json:"nama_cussubgroup,omitempty"`
	NamaCusjenis    string  `gorm:"-" json:"nama_cusjenis,omitempty"`
	NamaCusdivisi   string  `gorm:"-" json:"nama_cusdivisi,omitempty"`
	NamaCuswilayah  string  `gorm:"-" json:"nama_cuswilayah,omitempty"`
}

// TableName overrides the default table name for Customer
func (Customer) TableName() string {
	return "tb_customer"
}
