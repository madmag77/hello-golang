package main

import (
	"fmt"
	"os"
)

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
}

type FileTransactionalLogger struct {
	events chan<- Event
	errors <-chan error
	lastId uint64
	file   *os.File
}

func (l *FileTransactionalLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventDelete, Key: key}
}

func (l *FileTransactionalLogger) WritePut(key, value string) {
	l.events <- Event{EventType: EventDelete, Key: key, Value: value}
}

func (l *FileTransactionalLogger) Err() <-chan error {
	return l.errors
}

func CreateFileTransactionalLogger(filename string) (TransactionalLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open file for logs %w", err)
	}

	return &FileTransactionalLogger{file: file}, nil
}

func (l *FileTransactionalLogger) Run() {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		for e := range events {
			l.lastId++

			_, err := fmt.Fprintf(
				l.file,
				"%d\t%d\t%s\t%s\n",
				l.lastId, e.EventType, e.Key, e.Value)

			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (l *FileTransactionalLogger) ReadAll() (<-chan Event, error) {
	allEvents := make(chan Event)

	go func() {

	}()

	return allEvents, nil
}
