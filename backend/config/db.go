package config

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBPostgres *gorm.DB
var DBMysql *gorm.DB
var DBMysqlPos *gorm.DB

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
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
		fmt.Println("Koneksi Postgres (sam_live) Berhasil!")
	}

	// 2. MySQL Default Connection (selada)
	myHost := getEnv("DB_MYSQL_HOST", "10.10.10.8")
	myPort := getEnv("DB_MYSQL_PORT", "3307")
	myUser := getEnv("DB_MYSQL_USER", "root")
	myPass := getEnv("DB_MYSQL_PASS", "stagingdev@6177")
	myName := getEnv("DB_MYSQL_NAME", "selada_dev")
	dsnMS := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", myUser, myPass, myHost, myPort, myName)

	DBMysql, err = gorm.Open(mysql.Open(dsnMS), &gorm.Config{})
	if err != nil {
		fmt.Printf("Gagal koneksi ke MySQL (selada): %v\n", err)
		DBMysql = nil
	} else {
		fmt.Println("Koneksi MySQL (selada) Berhasil!")
	}

	// 3. MySQL POS Connection (sam_pos)
	myPosHost := getEnv("DB_MYSQL_POS_HOST", "10.10.10.8")
	myPosPort := getEnv("DB_MYSQL_POS_PORT", "3307")
	myPosUser := getEnv("DB_MYSQL_POS_USER", "root")
	myPosPass := getEnv("DB_MYSQL_POS_PASS", "stagingdev@6177")
	myPosName := getEnv("DB_MYSQL_POS_NAME", "sam_pos_dev")
	dsnMSPos := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", myPosUser, myPosPass, myPosHost, myPosPort, myPosName)

	DBMysqlPos, err = gorm.Open(mysql.Open(dsnMSPos), &gorm.Config{})
	if err != nil {
		fmt.Printf("Gagal koneksi ke MySQL (sam_pos): %v\n", err)
		DBMysqlPos = nil
	} else {
		fmt.Println("Koneksi MySQL (sam_pos) Berhasil!")
	}
}
