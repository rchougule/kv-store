package kvstore

/*
	KVStore interface
	Author: rchougule
	LastUpdated: 2023-07-01
*/

type KVStore interface {
	// Get fetches the value associated with the particular key from the KVStore
	Get(key interface{}) (interface{}, error)
	// Put inserts new or updates an existing value associated with the given key
	Put(key interface{}, value interface{}) error
	// Keys returns the list of keys present in the KVStore
	Keys() []interface{}
	// Count returns the total entries present in the KVStore
	Count() int64
}
