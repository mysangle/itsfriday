package db

import (
	"errors"
	"fmt"

	"itsfriday/server/profile"
	"itsfriday/store"
	"itsfriday/store/db/sqlite"
)

func NewDBDriver(profile *profile.Profile) (store.Driver, error) {
	var driver store.Driver
	var err error

	switch profile.Driver {
	case "sqlite":
		driver, err = sqlite.NewDB(profile)
	default:
		return nil, errors.New("unknown db driver")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create db driver: %w", err)
	}
	return driver, nil
}
