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

	fmt.Println("=== Checking tables in db_gadaimulia ===")
	if config.DBSelada != nil {
		var columns []struct {
			ColumnName string `gorm:"column:column_name"`
			DataType   string `gorm:"column:data_type"`
		}
		config.DBSelada.Raw(`
			SELECT column_name, data_type 
			FROM information_schema.columns 
			WHERE table_name = 'tb_customer'
			ORDER BY ordinal_position
		`).Scan(&columns)
		fmt.Printf("=== tb_customer columns (%d columns) ===\n", len(columns))
		for _, c := range columns {
			fmt.Printf("  %s (%s)\n", c.ColumnName, c.DataType)
		}
	}
	os.Exit(0)
}
