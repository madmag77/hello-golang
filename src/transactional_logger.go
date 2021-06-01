package main

type EventType byte

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)

type Event struct {
	Id        uint64
	EventType EventType
	Key       string
	Value     string
}

type TransactionalLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)

	ReadAll() (<-chan Event, <-chan error)
	Run() <-chan error
}
