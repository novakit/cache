package redis // import "github.com/novakit/cache/redis"
import (
	"time"

	"github.com/go-redis/redis"
	"github.com/novakit/cache"
)

// AdapterName adapter name for redis
const AdapterName = "redis"

// Adapter redis adapter
type Adapter struct{}

// Instance implements cache.Adapter
func (Adapter) Instance(options string) (adi cache.AdapterInstance, err error) {
	var opt *redis.Options
	if len(options) == 0 {
		options = "redis://127.0.0.1:6379"
	}
	if opt, err = redis.ParseURL(options); err != nil {
		return
	}
	// github.com/go-redis/redis has a circuit breaker, so we assume it's ok not to do a ping test
	adi = &AdapterInstance{Client: redis.NewClient(opt)}
	return
}

// AdapterInstance implements cache.AdapterInstance
type AdapterInstance struct {
	Client *redis.Client
}

// Get implements cache.AdapterInstance
func (a *AdapterInstance) Get(key string) (s string, err error) {
	if s, err = a.Client.Get(key).Result(); len(s) == 0 {
		err = cache.ErrKeyNotFound
	}
	return
}

// Set implements cache.AdapterInstance
func (a *AdapterInstance) Set(key, val string, expires int) error {
	return a.Client.Set(key, val, time.Second*time.Duration(expires)).Err()
}

// Del implements cache.AdapterInstance
func (a *AdapterInstance) Del(key string) error {
	return a.Client.Del(key).Err()
}

func init() {
	cache.RegisterAdapter(AdapterName, Adapter{})
}
