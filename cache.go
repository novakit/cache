package cache // import "github.com/novakit/cache"
import (
	"encoding/json"
	"fmt"

	"github.com/novakit/nova"
)

// ContextKey context key for cache module
const ContextKey = "_cache"

// Options options for cache module
type Options struct {
	// Adapter adapter name
	Adapter string
	// AdapterOptions adapter config string
	AdapterOptions string
}

func sanitizeOptions(opts ...Options) (opt Options) {
	if len(opts) > 0 {
		opt = opts[0]
	}
	if len(opt.Adapter) == 0 {
		opt.Adapter = MemoryAdapterName
	}
	return
}

// Cache the cache struct
type Cache struct {
	AdapterInstance
}

// GetJSON get a cache string and unmarshal to json object
func (c *Cache) GetJSON(key string, out interface{}) (err error) {
	var s string
	if s, err = c.AdapterInstance.Get(key); err != nil {
		return
	}
	err = json.Unmarshal([]byte(s), out)
	return
}

// SetJSON marshal a object and put into cache
func (c *Cache) SetJSON(key string, in interface{}, expires int) (err error) {
	var buf []byte
	if buf, err = json.Marshal(in); err != nil {
		return
	}
	err = c.Set(key, string(buf), expires)
	return
}

// Handler create a handler injects Cache
func Handler(opts ...Options) nova.HandlerFunc {
	opt := sanitizeOptions(opts...)
	// find adapter
	adp := adapters[opt.Adapter]
	if adp == nil {
		panic(fmt.Errorf("cache: can not find a adapter named: %s", opt.Adapter))
	}
	// create instance
	var ins AdapterInstance
	var err error
	if ins, err = adp.Instance(opt.AdapterOptions); err != nil {
		panic(err)
	}
	return func(c *nova.Context) (err error) {
		c.Values[ContextKey] = &Cache{AdapterInstance: ins}
		c.Next()
		return
	}
}

// Extract extract a Cache from injected nova.Context
func Extract(c *nova.Context) (ca *Cache) {
	ca, _ = c.Values[ContextKey].(*Cache)
	return
}
