package main

import (
	"errors"
	"testing"
)

func Test_when_Put_KV_then_can_read_value(t *testing.T) {
	// Given
	testKey := "testKey"
	testValue := "testValue"

	// When
	Put(testKey, testValue)
	val, err := Get(testKey)

	// Then
	if err != nil {
		t.Error("Get shouldn't return error")
	}

	if val != testValue {
		t.Error("Get should should return proper value")
	}
}

func Test_when_twice_Put_KV_then_can_still_read_value(t *testing.T) {
	// Given
	testKey := "testKey"
	testValue := "testValue"

	// When
	Put(testKey, testValue)
	Put(testKey, testValue)
	val, err := Get(testKey)

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
	Put(testKey, testValue)

	// When
	val, err := Get(wrongKey)

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
	Put(testKey, testValue)

	// When
	delete_err := Delete(testKey)

	// Then
	if delete_err != nil {
		t.Error("Delete shouldn't return error")
	}

	val, err := Get(testKey)

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
