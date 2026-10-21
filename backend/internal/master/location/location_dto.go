package location

type KelurahanResponse struct {
	KdKelurahan int    `json:"kd_kelurahan"`
	NmKelurahan string `json:"nm_kelurahan"`
	KdKecamatan int    `json:"kd_kecamatan"`
	NmKecamatan string `json:"nm_kecamatan"`
	KdKota      int    `json:"kd_kota"`
	NmKota      string `json:"nm_kota"`
}
