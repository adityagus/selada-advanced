package sales

type SalesResponse struct {
	Username   string  `json:"username"`
	Nama       string  `json:"nama"`
	IDJabatan  int     `json:"id_jabatan"`
	NmJabatan  string  `json:"nm_jabatan"`
	FKUser     string  `json:"fkuser"`
	TotalOmzet float64 `json:"total_omzet"`
}
