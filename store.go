package edkvs

import (
	"encoding/binary"
	"sync"
	"time"
)

type Store struct {
	items        map[string]*item
	itemsRWMutex sync.RWMutex
	state        *Set
	count        int
}

func NewStore() *Store {
	return &Store{
		items: make(map[string]*item),
		state: NewSet(),
		count: 0,
	}
}

func (s *Store) Set(key, value []byte) error {
	keyString := string(key)
	s.itemsRWMutex.Lock()
	if i, ok := s.items[keyString]; ok {
		s.state.Remove(stateItem(key, i.revision))
		i.value = value
		i.revision++
		if !i.deletedAt.IsZero() {
			i.deletedAt = time.Time{}
			s.count++
		}
		s.state.Insert(stateItem(key, i.revision))
	} else {
		s.items[keyString] = &item{value: value, revision: 0, deletedAt: time.Time{}}
		s.state.Insert(stateItem(key, 0))
		s.count++
	}
	s.itemsRWMutex.Unlock()
	return nil
}

func (s *Store) Get(key []byte) ([]byte, error) {
	s.itemsRWMutex.RLock()

	i, ok := s.items[string(key)]
	if ok {
		value := i.value
		s.itemsRWMutex.RUnlock()
		return value, nil
	}
	s.itemsRWMutex.RUnlock()
	return nil, nil
}

func (s *Store) Delete(key []byte) error {
	keyString := string(key)
	s.itemsRWMutex.Lock()
	if i, ok := s.items[keyString]; ok {
		s.state.Remove(stateItem(key, i.revision))
		i.value = nil
		i.deletedAt = time.Now()
		i.revision++
		s.state.Insert(stateItem(key, i.revision))
		s.count--
	}
	s.itemsRWMutex.Unlock()
	return nil
}

func (s *Store) Len() int {
	return s.count
}

func (s *Store) State() *Set {
	return s.state
}

type item struct {
	value     []byte
	revision  uint64
	deletedAt time.Time
}

func stateItem(key []byte, revision uint64) []byte {
	revisionBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(revisionBytes, revision)
	return append(key, revisionBytes...)
}
