package edkvs

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"sync"
	"time"

	"github.com/simia-tech/errx"
)

const keyHashSize = 8

type keyHash [keyHashSize]byte

func newKeyHash(data []byte) keyHash {
	kh := keyHash{}
	copy(kh[:], data[:keyHashSize])
	return kh
}

func (kh keyHash) String() string {
	return hex.EncodeToString(kh[:])
}

// Store defines a key-value store.
type Store struct {
	items        map[keyHash]*item
	itemsRWMutex sync.RWMutex
	state        *Set
	count        int
	updateFn     func(keyHash, *item)
}

// NewStore returns a new store.
func NewStore() *Store {
	return &Store{
		items: make(map[keyHash]*item),
		state: NewSet(),
		count: 0,
	}
}

// Set sets the provided value at the provided key.
func (s *Store) Set(key, value []byte) error {
	kh := hashKey(key)
	s.itemsRWMutex.Lock()
	if i, ok := s.items[kh]; ok {
		s.state.Remove(stateItem(kh, i.revision))
		i.value = value
		i.revision++
		if !i.deletedAt.IsZero() {
			i.deletedAt = time.Time{}
			s.count++
		}
		s.state.Insert(stateItem(kh, i.revision))
		s.notify(kh, i)
	} else {
		i := &item{value: value, revision: 0, deletedAt: time.Time{}}
		s.items[kh] = i
		s.state.Insert(stateItem(kh, 0))
		s.notify(kh, i)
		s.count++
	}
	s.itemsRWMutex.Unlock()
	return nil
}

// Get returns the value at the provided key. If no value exists, nil is returned.
func (s *Store) Get(key []byte) ([]byte, error) {
	kh := hashKey(key)
	s.itemsRWMutex.RLock()
	i, ok := s.items[kh]
	if ok {
		if !i.deletedAt.IsZero() {
			s.itemsRWMutex.RUnlock()
			return nil, nil
		}
		value := i.value
		s.itemsRWMutex.RUnlock()
		return value, nil
	}
	s.itemsRWMutex.RUnlock()
	return nil, nil
}

// Delete removes the value at the provided key.
func (s *Store) Delete(key []byte) error {
	k := hashKey(key)
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

// Len returns the length of the store.
func (s *Store) Len() int {
	return s.count
}

// State returns a set containing all keys and revisions.
func (s *Store) State() *Set {
	return s.state
}

func (s *Store) setItem(kh keyHash, bytes []byte) error {
	ni := &item{}
	if err := ni.UnmarshalBinary(bytes); err != nil {
		return errx.Annotatef(err, "unmarshal binary")
	}

	s.itemsRWMutex.Lock()
	if i, ok := s.items[kh]; ok {
		s.state.Remove(stateItem(kh, i.revision))
		s.items[kh] = ni
		switch {
		case i.deletedAt.IsZero() && !ni.deletedAt.IsZero():
			s.count--
		case !i.deletedAt.IsZero() && ni.deletedAt.IsZero():
			s.count++
		}
		s.state.Insert(stateItem(kh, ni.revision))
	} else {
		s.items[kh] = ni
		s.state.Insert(stateItem(kh, 0))
		if ni.deletedAt.IsZero() {
			s.count++
		}
	}
	s.itemsRWMutex.Unlock()

	return nil
}

func (s *Store) getItem(kh keyHash) ([]byte, error) {
	s.itemsRWMutex.RLock()
	i, ok := s.items[kh]
	if ok {
		bytes, err := i.MarshalBinary()
		if err != nil {
			s.itemsRWMutex.RUnlock()
			return nil, errx.Annotatef(err, "marshal binary")
		}
		s.itemsRWMutex.RUnlock()
		return bytes, nil
	}
	s.itemsRWMutex.RUnlock()
	return nil, nil
}

func (s *Store) notify(kh keyHash, item *item) {
	if s.updateFn == nil {
		return
	}
	s.updateFn(kh, item)
}

type item struct {
	value     []byte
	revision  uint64
	deletedAt time.Time
}

func (i *item) MarshalBinary() ([]byte, error) {
	buffer := make([]byte, len(i.value)+16)
	binary.BigEndian.PutUint64(buffer[:8], i.revision)
	binary.BigEndian.PutUint64(buffer[8:16], uint64(i.deletedAt.Unix()))
	copy(buffer[16:], i.value)
	return buffer, nil
}

func (i *item) UnmarshalBinary(data []byte) error {
	if len(data) < 16 {
		return errx.Errorf("need at least 16 bytes")
	}
	i.revision = binary.BigEndian.Uint64(data[:8])
	i.deletedAt = time.Unix(int64(binary.BigEndian.Uint64(data[8:16])), 0)
	i.value = data[16:]
	return nil
}

func hashKey(k []byte) keyHash {
	hash := sha1.Sum(k)
	kh := keyHash{}
	copy(kh[:], hash[:keyHashSize])
	return kh
}

func stateItem(key keyHash, revision uint64) Item {
	item := Item{}
	copy(item[:keyHashSize], key[:])
	binary.BigEndian.PutUint64(item[keyHashSize:], revision)
	return item
}
