package cachefile

import (
	"sync"

	"github.com/metacubex/bbolt"
	"github.com/metacubex/mihomo/component/smart"
	"github.com/metacubex/mihomo/log"
)

var (
	smartInitMutex sync.Mutex
	smartInitDone  bool
	smartStore     *smart.Store
)

type SmartStore struct {
	store *smart.Store
}

func NewSmartStore(cache *CacheFile) *SmartStore {
	if cache == nil || cache.DB == nil {
		return nil
	}

	smartInitMutex.Lock()
	defer smartInitMutex.Unlock()

	if !smartInitDone {
		err := cache.DB.Update(func(tx *bbolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(bucketSmartStats)
			return err
		})
		if err != nil {
			log.Warnln("[SmartStore] Failed to create bucket: %v", err)
			return nil
		}
		smart.InitGlobalParams()
		smartStore = smart.NewStore(cache.DB)
		smartInitDone = true
	}

	if smartStore == nil {
		return nil
	}

	return &SmartStore{
		store: smartStore,
	}
}

func (s *SmartStore) GetStore() *smart.Store {
	return s.store
}