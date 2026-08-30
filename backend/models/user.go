package models

// UserPos represents the user table in the MySQL POS database (sam_pos.tbl_user)
type UserPos struct {
	Username  string `gorm:"column:username;primaryKey"`
	Nama      string `gorm:"column:nama"`
	Password  string `gorm:"column:password"`
	IDCabang  string `gorm:"column:id_cabang"`
	IDJabatan int    `gorm:"column:id_jabatan"`
	NmJabatan string `gorm:"column:nm_jabatan"`
	Void      int    `gorm:"column:void"`
}

// TableName overrides the default table name for UserPos
func (UserPos) TableName() string {
	return "tbl_user"
}

// UserPg represents the user reference in the PostgreSQL database
type UserPg struct {
	Username     string  `gorm:"column:username;primaryKey"`
	EmployeeID   *string `gorm:"column:employee_id"`
	EmployeeCode *string `gorm:"column:fk_karyawan"`
	IsActive     *bool   `gorm:"column:is_active"`
}

// TableName overrides the default table name for UserPg
func (UserPg) TableName() string {
	return "auth.users"
}
