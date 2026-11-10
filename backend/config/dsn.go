package config

import (
	"fmt"
	"net"
	"net/url"
	"time"

	mysql "github.com/go-sql-driver/mysql"
)

func PostgresDSN(prefix, database string) (string, error) {
	for _, suffix := range []string{"HOST", "PORT", "USER", "NAME"} {
		if getEnv(prefix+"_"+suffix, "") == "" {
			return "", fmt.Errorf("%s_%s wajib diisi", prefix, suffix)
		}
	}
	if database == "" {
		database = getEnv(prefix+"_NAME", "")
	}
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(getEnv(prefix+"_USER", ""), getEnv(prefix+"_PASS", "")),
		Host:   net.JoinHostPort(getEnv(prefix+"_HOST", ""), getEnv(prefix+"_PORT", "")),
		Path:   database,
	}
	q := u.Query()
	q.Set("sslmode", getEnv(prefix+"_SSLMODE", "disable"))
	q.Set("connect_timeout", "5")
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func SourceMySQLDSN() (string, error) {
	for _, suffix := range []string{"HOST", "PORT", "USER", "NAME"} {
		if getEnv("DB_MYSQL_"+suffix, "") == "" {
			return "", fmt.Errorf("DB_MYSQL_%s wajib diisi", suffix)
		}
	}
	c := mysql.NewConfig()
	c.User = getEnv("DB_MYSQL_USER", "")
	c.Passwd = getEnv("DB_MYSQL_PASS", "")
	c.Net = "tcp"
	c.Addr = net.JoinHostPort(getEnv("DB_MYSQL_HOST", ""), getEnv("DB_MYSQL_PORT", ""))
	c.DBName = getEnv("DB_MYSQL_NAME", "")
	c.Timeout = 5 * time.Second
	return c.FormatDSN(), nil
}
