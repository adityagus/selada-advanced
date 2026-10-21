package branch

// Branch maps to master.branch in PostgreSQL
type Branch struct {
	BranchID   int    `gorm:"column:branch_id;primaryKey" json:"branch_id"`
	BranchCode string `gorm:"column:branch_code" json:"branch_code"`
	BranchName string `gorm:"column:branch_name" json:"branch_name"`
	IsActive   bool   `gorm:"column:is_active;default:true" json:"is_active"`
}

func (Branch) TableName() string {
	return "master.branch"
}
