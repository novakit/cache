package cache // import "github.com/novakit/cache"
import (
	"sync"
	"time"
)

// MemoryAdapterName adapter name for memory
const MemoryAdapterName = "memory"

// MemoryAdapter a memory adapter
type MemoryAdapter struct{}

// Instance implements Adapter
func (m MemoryAdapter) Instance(options string) (i AdapterInstance, err error) {
	return &MemoryAdapterInstance{
		Values:  map[string]string{},
		Expires: map[string]time.Time{},
		Mutex:   &sync.Mutex{},
	}, nil
}

// MemoryAdapterInstance implements AdapterInstance
type MemoryAdapterInstance struct {
	// Values values storage
	Values map[string]string
	// Expires expires time storage
	Expires map[string]time.Time
	// RWMutex locks
	Mutex *sync.Mutex
}

// Get implements AdapterInstance
func (m *MemoryAdapterInstance) Get(key string) (val string, err error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	exp := m.Expires[key]
	if !exp.IsZero() && time.Now().After(exp) {
		delete(m.Expires, key)
		delete(m.Values, key)
		err = ErrKeyNotFound
		return
	}
	var ok bool
	if val, ok = m.Values[key]; !ok {
		err = ErrKeyNotFound
		return
	}
	return
}

// Set implements AdapterInstance
func (m *MemoryAdapterInstance) Set(key string, val string, expires int) (err error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.Values[key] = val
	m.Expires[key] = time.Now().Add(time.Second * time.Duration(expires))
	return
}

// Del implements AdapterInstance
func (m *MemoryAdapterInstance) Del(key string) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	delete(m.Expires, key)
	delete(m.Values, key)
	return nil
}

func init() {
	RegisterAdapter(MemoryAdapterName, MemoryAdapter{})
}
