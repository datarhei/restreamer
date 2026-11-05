package store

type Store interface {
	// Load data from the store
	Load() (StoreData, error)

	// Save data to the store
	Store(data StoreData) error
}
