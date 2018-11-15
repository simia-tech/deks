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
	containers        map[keyHash]*container
	containersRWMutex sync.RWMutex
	state             *Set
	count             int
	updateFn          func(keyHash, *container)
}

// NewStore returns a new store.
func NewStore() *Store {
	return &Store{
		containers: make(map[keyHash]*container),
		state:      NewSet(),
		count:      0,
	}
}

// Set sets the provided value at the provided key.
func (s *Store) Set(key, value []byte) error {
	kh := hashKey(key)
	s.containersRWMutex.Lock()
	if c, ok := s.containers[kh]; ok {
		s.state.Remove(stateItem(kh, c.revision))
		c.value = value
		c.revision++
		if c.isDeleted() {
			c.undelete()
			s.count++
		}
		s.state.Insert(stateItem(kh, c.revision))
		s.notify(kh, c)
	} else {
		c := &container{
			key:       key,
			value:     value,
			revision:  0,
			deletedAt: time.Time{},
		}
		s.containers[kh] = c
		s.state.Insert(stateItem(kh, 0))
		s.notify(kh, c)
		s.count++
	}
	s.containersRWMutex.Unlock()
	return nil
}

// Get returns the value at the provided key. If no value exists, nil is returned.
func (s *Store) Get(key []byte) ([]byte, error) {
	kh := hashKey(key)
	s.containersRWMutex.RLock()
	c, ok := s.containers[kh]
	if ok {
		if c.isDeleted() {
			s.containersRWMutex.RUnlock()
			return nil, nil
		}
		value := c.value
		s.containersRWMutex.RUnlock()
		return value, nil
	}
	s.containersRWMutex.RUnlock()
	return nil, nil
}

// Delete removes the value at the provided key.
func (s *Store) Delete(key []byte) error {
	hk := hashKey(key)
	s.containersRWMutex.Lock()
	if c, ok := s.containers[hk]; ok {
		if !c.isDeleted() {
			s.state.Remove(stateItem(hk, c.revision))
			c.value = nil
			c.delete()
			c.revision++
			s.state.Insert(stateItem(hk, c.revision))
			s.notify(hk, c)
			s.count--
		}
	}
	s.containersRWMutex.Unlock()
	return nil
}

// Each interates over all key-value-pairs.
func (s *Store) Each(fn func([]byte, []byte) error) (err error) {
	s.containersRWMutex.RLock()
	for _, c := range s.containers {
		if c.isDeleted() {
			continue
		}
		if err = fn(c.key, c.value); err != nil {
			break
		}
	}
	s.containersRWMutex.RUnlock()
	return
}

// Len returns the length of the store.
func (s *Store) Len() int {
	return s.count
}

// State returns a set containing all keys and revisions.
func (s *Store) State() *Set {
	return s.state
}

func (s *Store) setContainer(kh keyHash, bytes []byte) error {
	nc := &container{}
	if err := nc.UnmarshalBinary(bytes); err != nil {
		return errx.Annotatef(err, "unmarshal binary")
	}

	s.containersRWMutex.Lock()
	if c, ok := s.containers[kh]; ok {
		s.state.Remove(stateItem(kh, c.revision))
		s.containers[kh] = nc
		switch {
		case !c.isDeleted() && nc.isDeleted():
			s.count--
		case c.isDeleted() && !nc.isDeleted():
			s.count++
		}
		s.state.Insert(stateItem(kh, nc.revision))
	} else {
		s.containers[kh] = nc
		s.state.Insert(stateItem(kh, 0))
		if !nc.isDeleted() {
			s.count++
		}
	}
	s.containersRWMutex.Unlock()

	return nil
}

func (s *Store) getContainer(kh keyHash) ([]byte, error) {
	s.containersRWMutex.RLock()
	c, ok := s.containers[kh]
	if ok {
		bytes, err := c.MarshalBinary()
		if err != nil {
			s.containersRWMutex.RUnlock()
			return nil, errx.Annotatef(err, "marshal binary")
		}
		s.containersRWMutex.RUnlock()
		return bytes, nil
	}
	s.containersRWMutex.RUnlock()
	return nil, nil
}

func (s *Store) notify(kh keyHash, c *container) {
	if s.updateFn == nil {
		return
	}
	s.updateFn(kh, c)
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
