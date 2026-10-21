package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	mysql "github.com/go-sql-driver/mysql"
)

func PostgresDSN(prefix, database string) (string, error) {
	for _, suffix := range []string{"HOST", "PORT", "USER", "NAME"} {
		if os.Getenv(prefix+"_"+suffix) == "" {
			return "", fmt.Errorf("%s_%s wajib diisi", prefix, suffix)
		}
	}
	if database == "" {
		database = os.Getenv(prefix + "_NAME")
	}
	u := url.URL{Scheme: "postgres", User: url.UserPassword(os.Getenv(prefix+"_USER"), os.Getenv(prefix+"_PASS")), Host: net.JoinHostPort(os.Getenv(prefix+"_HOST"), os.Getenv(prefix+"_PORT")), Path: database}
	q := u.Query()
	q.Set("sslmode", getEnv(prefix+"_SSLMODE", "disable"))
	q.Set("connect_timeout", "5")
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func SourceMySQLDSN() (string, error) {
	for _, suffix := range []string{"HOST", "PORT", "USER", "NAME"} {
		if os.Getenv("DB_MYSQL_"+suffix) == "" {
			return "", fmt.Errorf("DB_MYSQL_%s wajib diisi", suffix)
		}
	}
	c := mysql.NewConfig()
	c.User = os.Getenv("DB_MYSQL_USER")
	c.Passwd = os.Getenv("DB_MYSQL_PASS")
	c.Net = "tcp"
	c.Addr = net.JoinHostPort(os.Getenv("DB_MYSQL_HOST"), os.Getenv("DB_MYSQL_PORT"))
	c.DBName = os.Getenv("DB_MYSQL_NAME")
	c.Timeout = 5 * time.Second
	return c.FormatDSN(), nil
}
