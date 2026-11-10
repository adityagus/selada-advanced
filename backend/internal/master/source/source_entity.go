package source

// ReferenceSource maps to master.reference_sources (sumber customer / aktivitas)
type ReferenceSource struct {
	ID       int    `gorm:"column:id;primaryKey" json:"id"`
	Code     string `gorm:"column:code" json:"code"`
	Name     string `gorm:"column:name" json:"name"`
	IsActive bool   `gorm:"column:is_active;default:true" json:"is_active"`
}

func (ReferenceSource) TableName() string {
	return "master.reference_sources"
}

// ApplicationSource maps to master.application_sources (SOA)
type ApplicationSource struct {
	ID       int    `gorm:"column:id;primaryKey" json:"id"`
	Code     string `gorm:"column:code" json:"code"`
	Name     string `gorm:"column:name" json:"name"`
	IsActive bool   `gorm:"column:is_active;default:true" json:"is_active"`
}

func (ApplicationSource) TableName() string {
	return "master.application_sources"
}

// ApplicationType maps to master.application_types (jenis aplikasi)
type ApplicationType struct {
	ID       int    `gorm:"column:id;primaryKey" json:"id"`
	Code     string `gorm:"column:code" json:"code"`
	Name     string `gorm:"column:name" json:"name"`
	IsActive bool   `gorm:"column:is_active;default:true" json:"is_active"`
}

func (ApplicationType) TableName() string {
	return "master.application_types"
}
