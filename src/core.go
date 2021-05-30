package main

import "errors"

var store = make(map[string]string)

// Put the value in the store
func Put(key string, value string) error {
	store[key] = value
	return nil
}

// Error showing that there is no such key in the store
var ErrorNoSuchKey = errors.New("No such key")

// Get the value from the store
func Get(key string) (string, error) {
	value, ok := store[key]

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

// Delete the key value pair from the store
func Delete(key string) error {
	delete(store, key)

	return nil
}
