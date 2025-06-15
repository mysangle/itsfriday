package profile

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Profile struct {
    // Mode can be "prod" or "dev"
	Mode string
	// Addr is the binding address for server
	Addr string
	// Port is the binding port for server
	Port int
	// Data is the data directory
	Data string
	// DSN points to where itsfriday stores its own data
	DSN string
	// Driver is the database driver: sqlite
	Driver string
	// Version is the current version of server
	Version string
}

func (p *Profile) IsDev() bool {
	return p.Mode != "prod"
}

func (p *Profile) Validate() error {
	dataDir, err := checkDataDir(p.Data)
	if err != nil {
		slog.Error("failed to check dsn", slog.String("data", dataDir), slog.String("error", err.Error()))
		return err
	}

	p.Data = dataDir
	if p.Driver == "sqlite" && p.DSN == "" {
		dbFile := fmt.Sprintf("itsfriday_%s.db", p.Mode)
		p.DSN = filepath.Join(dataDir, dbFile)
	}

    return nil
}

func checkDataDir(dataDir string) (string, error) {
	// Convert to absolute path if relative path is supplied.
	if !filepath.IsAbs(dataDir) {
		relativeDir := filepath.Join(filepath.Dir(os.Args[0]), dataDir)
		absDir, err := filepath.Abs(relativeDir)
		if err != nil {
			return "", err
		}
		dataDir = absDir
	}

	// Trim trailing \ or / in case user supplies
	dataDir = strings.TrimRight(dataDir, "\\/")
	if _, err := os.Stat(dataDir); err != nil {
		return "", fmt.Errorf("unable to access data folder %s: %w", dataDir, err)
	}

	return dataDir, nil
}
