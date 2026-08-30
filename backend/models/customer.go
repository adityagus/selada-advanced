package models

// Customer represents the tb_customer table in MySQL default db
type Customer struct {
	IDCustomer int    `gorm:"column:id_customer;primaryKey;autoIncrement"`
	Nama       string `gorm:"column:nama"`
	Hp         string `gorm:"column:hp"`
	Alamat     string `gorm:"column:alamat"`
	IDLogin    string `gorm:"column:id_login"`
}

// TableName overrides the default table name for Customer
func (Customer) TableName() string {
	return "tb_customer"
}
