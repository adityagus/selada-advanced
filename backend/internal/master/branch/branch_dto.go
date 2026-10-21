package branch

// BranchResponse represents standard output for branch data
type BranchResponse struct {
	IDCabang   int    `json:"id_cabang"`
	KodeCabang string `json:"kode_cabang"`
	NmCabang   string `json:"nm_cabang"`
}

// CurrentBranchResponse represents the active user's branch context
type CurrentBranchResponse struct {
	Cabang     string `json:"cabang"`
	BranchName string `json:"branch_name"`
}
