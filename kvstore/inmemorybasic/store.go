package inmemorybasic

import "github.com/rchougule/kv-store/kvstore"

type Store struct {
	kv map[string]interface{}
}

func NewStore() *Store {
	return &Store{
		kv: make(map[string]interface{}),
	}
}

func (s *Store) Get(key interface{}) (interface{}, error) {
	return s.kv[key.(string)], nil
}

func (s *Store) Put(key interface{}, value interface{}) error {
	s.kv[key.(string)] = value
	return nil
}

func (s *Store) Keys() []interface{} {
	//TODO implement me
	panic("implement me")
}

func (s *Store) Count() int64 {
	//TODO implement me
	panic("implement me")
}

var _ kvstore.KVStore = &Store{}
