package main

import (
	"errors"
	"sync"
)

// Error showing that there is no such key in the store
var ErrorNoSuchKey = errors.New("no such key")

type KeyValueStore interface {
	Put(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

func CreateKeyValueStore(transactionalLogger TransactionalLogger) (*SimpleKeyValueStore, error) {
	return &SimpleKeyValueStore{m: make(map[string]string), transactionalLogger: transactionalLogger}, nil
}

type SimpleKeyValueStore struct {
	sync.RWMutex
	m                   map[string]string
	transactionalLogger TransactionalLogger
}

func (store *SimpleKeyValueStore) put(key string, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()

	return nil
}

// Put the value in the store
func (store *SimpleKeyValueStore) Put(key string, value string) error {
	err := store.put(key, value)

	if err != nil {
		return err
	}

	store.transactionalLogger.WritePut(key, value)

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

func (store *SimpleKeyValueStore) delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()

	return nil
}

// Delete the key value pair from the store
func (store *SimpleKeyValueStore) Delete(key string) error {
	err := store.delete(key)
	if err != nil {
		return err
	}

	store.transactionalLogger.WriteDelete(key)

	return nil
}

type KeyValueStorePersistance interface {
	RestorePersistedState() error
}

func (store *SimpleKeyValueStore) RestorePersistedState() error {
	events, errors := store.transactionalLogger.ReadAll()

	var err error = nil
	e, ok := Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				err = store.delete(e.Key)

			case EventPut:
				err = store.put(e.Key, e.Value)
			}
		}
	}

	return err
}
