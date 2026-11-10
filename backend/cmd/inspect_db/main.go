package main

import (
	"fmt"
	"go-api/config"
	"os"
)

func main() {
	if err := config.LoadEnv(".env"); err != nil {
		fmt.Printf("LoadEnv error: %v\n", err)
		return
	}
	config.ConnectDatabases()

	type CabangPos struct {
		NmCabang string `gorm:"column:nm_cabang"`
	}
	var cPos CabangPos
	err := config.DBMysqlPos.Table("mst_cabang").
		Select("nm_cabang").
		Where("id_cabang = ?", 42).
		Limit(1).Scan(&cPos).Error
	fmt.Printf("Query mst_cabang id_cabang=42: err=%v, res=%+v\n", err, cPos)
	os.Exit(0)
}
