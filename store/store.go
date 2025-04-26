package store

import (
	"sync"

    "itsfriday/server/profile"
)

type Store struct {
	Profile *profile.Profile
	driver  Driver

	userCache             sync.Map
	userSettingCache      sync.Map
}

func New(driver Driver, profile *profile.Profile) *Store {
	return &Store{
		driver:  driver,
		Profile: profile,
	}
}

func (s *Store) GetDriver() Driver {
	return s.driver
}

func (s *Store) Close() error {
	return s.driver.Close()
}
