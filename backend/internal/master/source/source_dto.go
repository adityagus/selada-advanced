package source

type SumberCustResponse struct {
	KdRefCust string `json:"kd_ref_cust"`
	NmRefCust string `json:"nm_ref_cust"`
}

type ApplicationSourceResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
