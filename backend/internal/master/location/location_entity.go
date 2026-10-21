package location

type VillageResult struct {
	KdKelurahan int    `json:"kd_kelurahan" gorm:"column:kd_kelurahan"`
	NmKelurahan string `json:"nm_kelurahan" gorm:"column:nm_kelurahan"`
	KdKecamatan int    `json:"kd_kecamatan" gorm:"column:kd_kecamatan"`
	NmKecamatan string `json:"nm_kecamatan" gorm:"column:nm_kecamatan"`
	KdKota      int    `json:"kd_kota" gorm:"column:kd_kota"`
	NmKota      string `json:"nm_kota" gorm:"column:nm_kota"`
}
