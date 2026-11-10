package source

type SumberCustResponse struct {
	ID        int    `json:"id"`
	KdRefCust string `json:"kd_ref_cust"`
	NmRefCust string `json:"nm_ref_cust"`
}

type ApplicationSourceResponse struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type ApplicationTypeResponse struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type SoaOptionsResponse struct {
	SOA           []ApplicationSourceResponse `json:"soa"`
	Aktifitas     []SumberCustResponse        `json:"aktifitas"`
	JenisAplikasi []ApplicationTypeResponse   `json:"jenis_aplikasi"`
}
