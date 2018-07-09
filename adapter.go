package cache // import "github.com/novakit/cache"
import "errors"

var adapters = map[string]Adapter{}

// ErrKeyNotFound error will be returned if key is not found in cache
var ErrKeyNotFound = errors.New("cache: key not found")

// AdapterInstance a adapter instance
type AdapterInstance interface {
	// Get get the cached string by key, if key not found, ErrKeyNotFound will be returned
	Get(key string) (string, error)
	// Set set the cached string by key with expiration in seconds
	Set(key, val string, expires int) error
	// Del delete the cached string by key
	Del(key string) error
}

// Adapter the adapter interface
type Adapter interface {
	// Instance create a new instance based on options
	Instance(options string) (AdapterInstance, error)
}

// RegisterAdapter register a adapter
func RegisterAdapter(name string, adapter Adapter) {
	adapters[name] = adapter
}
