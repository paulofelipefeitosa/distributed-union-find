package utils

import "sync"

type Void struct {}

type ConcurrentSet struct {
	set     sync.Map
	nothing Void
	rwLock  sync.RWMutex
}

func (s *ConcurrentSet) Has(element interface{}) bool {
	s.rwLock.RLock()
	_, exists := s.set.Load(element)
	s.rwLock.RUnlock()
	return exists
}

func (s *ConcurrentSet) Insert(element interface{}) {
	s.rwLock.RLock()
	s.set.Store(element, s.nothing)
	s.rwLock.RUnlock()
}

func (s *ConcurrentSet) Remove(element interface{}) {
	s.rwLock.RLock()
	s.set.Delete(element)
	s.rwLock.RUnlock()
}

func (s ConcurrentSet) Keys() []interface{} {
	maps := make(map[interface{}]bool)

	s.rwLock.Lock()
	s.set.Range(func(key, value interface{}) bool {
		maps[key] = true
		return true
	})
	s.rwLock.Unlock()

	keys := make([]interface{}, 0, len(maps))
	for k := range maps {
		keys = append(keys, k)
	}
	return keys
}