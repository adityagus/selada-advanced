package config

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBPostgres *gorm.DB
var DBSelada *gorm.DB
var DBMysqlPos *gorm.DB

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// tunePool recycles idle connections before firewalls / DB servers silently drop them,
// which otherwise makes the first requests after a long idle period fail.
func tunePool(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
}

// ensureCustomerColumns adds the aktivitas / jenis columns to tb_customer when missing.
func ensureCustomerColumns(db *gorm.DB) {
	for _, stmt := range []string{
		"ALTER TABLE tb_customer ADD COLUMN IF NOT EXISTS aktivitas VARCHAR(255) NOT NULL DEFAULT ''",
		"ALTER TABLE tb_customer ADD COLUMN IF NOT EXISTS jenis VARCHAR(255) NOT NULL DEFAULT ''",
	} {
		if err := db.Exec(stmt).Error; err != nil {
			fmt.Printf("Gagal menambah kolom tb_customer: %v\n", err)
		}
	}
}

func ConnectDatabases() {
	var err error

	// 1. PostgreSQL DB Connection (sam_live)
	pgHost := getEnv("DB_PG_HOST", "10.10.10.8")
	pgPort := getEnv("DB_PG_PORT", "5434")
	pgUser := getEnv("DB_PG_USER", "it_gms")
	pgPass := getEnv("DB_PG_PASS", "Gadai@passw0rd")
	pgName := getEnv("DB_PG_NAME", "db_gadaimulia")
	dsnPG := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", pgHost, pgUser, pgPass, pgName, pgPort)

	DBPostgres, err = gorm.Open(postgres.Open(dsnPG), &gorm.Config{})
	if err != nil {
		fmt.Printf("Gagal koneksi ke Postgres (sam_live): %v\n", err)
		DBPostgres = nil
	} else {
		tunePool(DBPostgres)
		fmt.Println("Koneksi Postgres (sam_live) Berhasil!")
	}

	// 2. PostgreSQL Selada connection; DB_MYSQL_* is only used by migration.
	dsnSelada, dsnErr := PostgresDSN("DB_PG2", "")
	if dsnErr != nil {
		fmt.Printf("Konfigurasi PostgreSQL (selada) tidak valid: %v\n", dsnErr)
		DBSelada = nil
	} else {
		DBSelada, err = gorm.Open(postgres.Open(dsnSelada), &gorm.Config{})
		if err != nil {
			fmt.Printf("Gagal koneksi ke PostgreSQL (selada): %v\n", err)
			DBSelada = nil
		} else {
			tunePool(DBSelada)
			ensureCustomerColumns(DBSelada)
			fmt.Println("Koneksi PostgreSQL (selada) Berhasil!")
		}
	}

	// 3. MySQL POS Connection (sam_pos)
	myPosHost := getEnv("DB_MYSQL_POS_HOST", "10.10.10.8")
	myPosPort := getEnv("DB_MYSQL_POS_PORT", "3307")
	myPosUser := getEnv("DB_MYSQL_POS_USER", "root")
	myPosPass := getEnv("DB_MYSQL_POS_PASS", "stagingdev@6177")
	myPosName := getEnv("DB_MYSQL_POS_NAME", "sam_pos")
	dsnMSPos := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", myPosUser, myPosPass, myPosHost, myPosPort, myPosName)

	DBMysqlPos, err = gorm.Open(mysql.Open(dsnMSPos), &gorm.Config{})
	if err != nil {
		fmt.Printf("Gagal koneksi ke MySQL (sam_pos): %v\n", err)
		DBMysqlPos = nil
	} else {
		tunePool(DBMysqlPos)
		fmt.Println("Koneksi MySQL (sam_pos) Berhasil!")
	}
}
