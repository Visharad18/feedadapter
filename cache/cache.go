package cache

// Cache stores key-value pair where the key has to be string and value could be of `any` type
type Cache interface {
	Store(string, interface{}) error
	Get(string) (interface{}, error)
	GetAll() map[string]interface{}
}
