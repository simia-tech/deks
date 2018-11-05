package edkvs

import (
	"crypto/sha1"
	"encoding/binary"
	"sync"
	"time"
)

const keySize = 20

type key [keySize]byte

type Store struct {
	items        map[key]*item
	itemsRWMutex sync.RWMutex
	state        *Set
	count        int
}

func NewStore() *Store {
	return &Store{
		items: make(map[key]*item),
		state: NewSet(),
		count: 0,
	}
}

func (s *Store) Set(key, value []byte) error {
	k := mapKey(key)
	s.itemsRWMutex.Lock()
	if i, ok := s.items[k]; ok {
		s.state.Remove(stateItem(k, i.revision))
		i.value = value
		i.revision++
		if !i.deletedAt.IsZero() {
			i.deletedAt = time.Time{}
			s.count++
		}
		s.state.Insert(stateItem(k, i.revision))
	} else {
		s.items[k] = &item{value: value, revision: 0, deletedAt: time.Time{}}
		s.state.Insert(stateItem(k, 0))
		s.count++
	}
	s.itemsRWMutex.Unlock()
	return nil
}

func (s *Store) Get(key []byte) ([]byte, error) {
	k := mapKey(key)
	s.itemsRWMutex.RLock()
	i, ok := s.items[k]
	if ok {
		value := i.value
		s.itemsRWMutex.RUnlock()
		return value, nil
	}
	s.itemsRWMutex.RUnlock()
	return nil, nil
}

func (s *Store) Delete(key []byte) error {
	k := mapKey(key)
	s.itemsRWMutex.Lock()
	if i, ok := s.items[k]; ok {
		s.state.Remove(stateItem(k, i.revision))
		i.value = nil
		i.deletedAt = time.Now()
		i.revision++
		s.state.Insert(stateItem(k, i.revision))
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

func mapKey(k []byte) key {
	hash := sha1.Sum(k)
	result := key{}
	copy(result[:], hash[:keySize])
	return result
}

func stateItem(key key, revision uint64) Item {
	result := Item{}
	copy(result[:keySize], key[:])
	binary.BigEndian.PutUint64(result[keySize:], revision)
	return result
}
