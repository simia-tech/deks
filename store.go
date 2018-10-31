package edkvs

type Store struct {
	items map[string]item
}

func NewStore() *Store {
	return &Store{
		items: make(map[string]item),
	}
}

func (s *Store) Set(key, value []byte) error {
	s.items[string(key)] = item{
		value: value,
	}
	return nil
}

func (s *Store) Get(key []byte) ([]byte, error) {
	item, ok := s.items[string(key)]
	if ok {
		return item.value, nil
	}
	return nil, nil
}

func (s *Store) Delete(key []byte) error {
	delete(s.items, string(key))
	return nil
}

func (s *Store) Len() int {
	return len(s.items)
}

type item struct {
	value []byte
}
