package source

// ReferenceSource maps to master.reference_sources (sumber customer)
type ReferenceSource struct {
	Code     string `gorm:"column:code;primaryKey" json:"code"`
	Name     string `gorm:"column:name" json:"name"`
	IsActive bool   `gorm:"column:is_active;default:true" json:"is_active"`
}

func (ReferenceSource) TableName() string {
	return "master.reference_sources"
}

// ApplicationSource maps to master.application_sources
type ApplicationSource struct {
	Code     string `gorm:"column:code;primaryKey" json:"code"`
	Name     string `gorm:"column:name" json:"name"`
	IsActive bool   `gorm:"column:is_active;default:true" json:"is_active"`
}

func (ApplicationSource) TableName() string {
	return "master.application_sources"
}
