package edkvs

import "sync"

type Store struct {
	items        map[string]item
	itemsRWMutex sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		items: make(map[string]item),
	}
}

func (s *Store) Set(key, value []byte) error {
	s.itemsRWMutex.Lock()
	s.items[string(key)] = item{
		value: value,
	}
	s.itemsRWMutex.Unlock()
	return nil
}

func (s *Store) Get(key []byte) ([]byte, error) {
	s.itemsRWMutex.RLock()

	item, ok := s.items[string(key)]
	if ok {
		value := item.value
		s.itemsRWMutex.RUnlock()
		return value, nil
	}
	s.itemsRWMutex.RUnlock()
	return nil, nil
}

func (s *Store) Delete(key []byte) error {
	s.itemsRWMutex.Lock()
	delete(s.items, string(key))
	s.itemsRWMutex.Unlock()
	return nil
}

func (s *Store) Len() int {
	s.itemsRWMutex.RLock()
	length := len(s.items)
	s.itemsRWMutex.RUnlock()
	return length
}

type item struct {
	value []byte
}
