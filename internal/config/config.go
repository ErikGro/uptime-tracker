package config

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

type Config struct {
	AdminUser  string
	AdminPass  string
	ListenAddr string
	DBPath     string
}

func Load() (*Config, error) {
	if err := loadDotEnv(".env"); err != nil {
		return nil, fmt.Errorf("load .env: %w", err)
	}
	cfg := &Config{
		AdminUser:  getenv("ADMIN_USER", "admin"),
		AdminPass:  os.Getenv("ADMIN_PASS"),
		ListenAddr: getenv("LISTEN_ADDR", ":8080"),
		DBPath:     getenv("DB_PATH", "uptime.db"),
	}
	if cfg.AdminPass == "" {
		return nil, errors.New("ADMIN_PASS is required (set it in .env)")
	}
	return cfg, nil
}

func loadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		v = strings.Trim(v, `"'`)
		if _, exists := os.LookupEnv(k); exists {
			continue
		}
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return sc.Err()
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
