package crud

// Store is a simplified interface to a key-blob store supporting CRUD operations.
type Store interface {
	List() ([]string, error)
	Save(name string, data []byte) error
	Read(name string) ([]byte, error)
	Delete(name string) error
}

// HasConnect indicates that a struct must be initialized using the Connect
// method before the interface's methods are called.
type HasConnect interface {
	Connect() error
}

// HasClose indicates that a struct must be cleaned up using the Close
// method before the interface's methods are called.
type HasClose interface {
	Close() error
}
