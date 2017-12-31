package storage

import "gitlab.com/medakk/zyxdb/storage/chaining"

type StorageBackend interface {
	Insert(key, value string)
	Retrieve(key string) (string, bool)
}

var (
	m StorageBackend = chaining.New()
)

func Retrieve(key string) (string, bool) {
	value, ok := m.Retrieve(key)
	return value, ok
}

func Insert(key, value string) {
	m.Insert(key, value)
}
