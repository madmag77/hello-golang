package main

import (
	"errors"
	"sync"
)

// Error showing that there is no such key in the store
var ErrorNoSuchKey = errors.New("No such key")

type KeyValueStore interface {
	Put(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

func CreateKeyValueStore(transactionalLogger TransactionalLogger) (KeyValueStore, error) {
	return &SimpleKeyValueStore{m: make(map[string]string), transactionalLogger: transactionalLogger}, nil
}

type SimpleKeyValueStore struct {
	sync.RWMutex
	m                   map[string]string
	transactionalLogger TransactionalLogger
}

// Put the value in the store
func (store *SimpleKeyValueStore) Put(key string, value string) error {
	store.Lock()
	store.m[key] = value
	store.transactionalLogger.WritePut(key, value)
	store.Unlock()

	return nil
}

// Get the value from the store
func (store *SimpleKeyValueStore) Get(key string) (string, error) {
	store.RLock()
	value, ok := store.m[key]
	store.RUnlock()

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

// Delete the key value pair from the store
func (store *SimpleKeyValueStore) Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.transactionalLogger.WriteDelete(key)
	store.Unlock()

	return nil
}
