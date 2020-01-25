package sh

import (
	"fmt"
	. "github.com/luckywinds/xsh/base"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/ssh"
	"sync"
	"time"
)

var _sshClientCache *sshClientCache

type sshClientCache struct {
	mu    sync.Mutex
	cache *cache.Cache
}

func (s *sshClientCache) delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache.Delete(key)
}

func (s *sshClientCache) get(key string) (client *ssh.Client, err error) {
	v, ok := s.cache.Get(key)
	if !ok {
		return nil, fmt.Errorf("key[%s] not exist", key)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	err = s.cache.Replace(key, v, time.Duration(XConfig.Cache.ExpirationS)*time.Second)

	return v.(*ssh.Client), err
}

func (s *sshClientCache) add(key string, client *ssh.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache.SetDefault(key, client)
}

func initSshClientCache() {
	_sshClientCache = &sshClientCache{
		cache: cache.New(time.Duration(XConfig.Cache.ExpirationS)*time.Second, time.Duration(XConfig.Cache.CleanupIntervalS)*time.Second),
	}

	go func() {
		t := time.NewTicker(time.Duration(XConfig.Cache.TickerIntervalS) * time.Second)
		defer t.Stop()

		for range t.C {
			for key, value := range _sshClientCache.cache.Items() {
				_, _, err := value.Object.(*ssh.Client).Conn.SendRequest("keepalive@xxx", true, nil)
				if err != nil {
					_sshClientCache.delete(key)
				}
			}
		}
	}()
}
