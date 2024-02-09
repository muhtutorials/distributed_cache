package cache

import "time"

type Cacher interface {
	Set([]byte, []byte, time.Duration) error
	Get([]byte) ([]byte, error)
	Has([]byte) bool
	Delete([]byte) error
}
