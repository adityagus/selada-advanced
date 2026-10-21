package config

import (
	"bufio"
	"os"
	"strings"
)

// LoadEnv loads local defaults without overriding the process environment.
func LoadEnv(path string) error {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if ok {
			k = strings.TrimSpace(k)
			if _, exists := os.LookupEnv(k); !exists {
				os.Setenv(k, strings.Trim(strings.TrimSpace(v), "\"'"))
			}
		}
	}
	return s.Err()
}
