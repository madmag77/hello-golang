package main

import (
	"errors"
	"os"
	"testing"
)

var vars struct {
	store      KeyValueStore
	mockLogger MockTransactionalLogger
}

func TestMain(m *testing.M) {
	vars.mockLogger = MockTransactionalLogger{}
	vars.store, _ = CreateKeyValueStore(&vars.mockLogger)
	exitVal := m.Run()

	os.Exit(exitVal)
}

func Test_when_Put_KV_then_can_read_value(t *testing.T) {
	// Given
	testKey := "testKey"
	testValue := "testValue"

	// When
	vars.store.Put(testKey, testValue)
	val, err := vars.store.Get(testKey)

	// Then
	if err != nil {
		t.Error("Get shouldn't return error")
	}

	if val != testValue {
		t.Error("Get should should return proper value")
	}

	if !vars.mockLogger.putWasCalled {
		t.Error("Transaction logging of Put should be done")
	}
}

func Test_when_twice_Put_KV_then_can_still_read_value(t *testing.T) {
	// Given
	testKey := "testKey"
	testValue := "testValue"

	// When
	vars.store.Put(testKey, testValue)
	vars.store.Put(testKey, testValue)
	val, err := vars.store.Get(testKey)

	// Then
	if err != nil {
		t.Error("Get shouldn't return error")
	}

	if val != testValue {
		t.Error("Get should should return proper value")
	}
}

func Test_given_existed_KV_when_read_nonexisted_key_then_return_error_and_empty_value(t *testing.T) {
	// Given
	testKey := "testKey"
	wrongKey := "wrongKey"
	testValue := "testValue"
	vars.store.Put(testKey, testValue)

	// When
	val, err := vars.store.Get(wrongKey)

	// Then
	if err == nil {
		t.Error("Get should return error")
	}

	if val != "" {
		t.Error("Get should return empty value")
	}

	if !errors.Is(err, ErrorNoSuchKey) {
		t.Error("Get should return proper error")
	}
}

func Test_given_existed_KV_when_delete_key_and_read_deleted_key_then_return_error_and_empty_value(t *testing.T) {
	// Given
	testKey := "testKey"
	testValue := "testValue"
	vars.store.Put(testKey, testValue)

	// When
	delete_err := vars.store.Delete(testKey)

	// Then
	if delete_err != nil {
		t.Error("Delete shouldn't return error")
	}

	val, err := vars.store.Get(testKey)

	if err == nil {
		t.Error("Get should return error")
	}

	if val != "" {
		t.Error("Get should return empty value")
	}

	if !errors.Is(err, ErrorNoSuchKey) {
		t.Error("Get should return proper error")
	}

	if !vars.mockLogger.deleteWasCalled {
		t.Error("Transaction logging of Put should be done")
	}
}

type MockTransactionalLogger struct {
	deleteWasCalled bool
	putWasCalled    bool
}

func (l *MockTransactionalLogger) WriteDelete(key string) {
	l.deleteWasCalled = true
}

func (l *MockTransactionalLogger) WritePut(key, value string) {
	l.putWasCalled = true
}

func (l *MockTransactionalLogger) Run() <-chan error {
	return nil
}

func (l *MockTransactionalLogger) ReadAll() (<-chan Event, <-chan error) {
	return nil, nil
}
