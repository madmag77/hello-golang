package main

import (
	"bufio"
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

	ReadAll() (<-chan Event, <-chan error)
	Run() <-chan error
}

type FileTransactionalLogger struct {
	events chan<- Event
	errors <-chan error
	lastId uint64
	file   *os.File
}

func (l *FileTransactionalLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventDelete, Key: key, Value: "_"}
}

func (l *FileTransactionalLogger) WritePut(key, value string) {
	l.events <- Event{EventType: EventPut, Key: key, Value: value}
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

func (l *FileTransactionalLogger) Run() <-chan error {
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

	return l.errors
}

func (l *FileTransactionalLogger) ReadAll() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	allEvents := make(chan Event)
	errorChan := make(chan error, 1)

	go func() {
		var e Event

		defer close(allEvents)
		defer close(errorChan)

		for scanner.Scan() {
			line := scanner.Text()

			if _, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s\n", &e.Id, &e.EventType, &e.Key, &e.Value); err != nil {
				errorChan <- fmt.Errorf("inout parse error: %w", err)
				return
			}

			if l.lastId >= e.Id {
				errorChan <- fmt.Errorf("ids are out of order")
				return
			}

			l.lastId = e.Id

			allEvents <- e
		}

		if err := scanner.Err(); err != nil {
			errorChan <- fmt.Errorf("transactional log read failure: %w", err)
			return
		}
	}()

	return allEvents, errorChan
}
