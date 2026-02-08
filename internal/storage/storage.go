package storage

// Storage interface for context items
type Storage interface {
	Load() error
	Save() error
	GetAll() []interface{}
	Add(item interface{}) error
	Remove(id string) error
}
