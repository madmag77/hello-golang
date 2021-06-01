package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresDbParams struct {
	dbName   string
	host     string
	user     string
	password string
}

type PostgresTransactionalLogger struct {
	events chan<- Event
	errors <-chan error
	db     *sql.DB
}

func (l *PostgresTransactionalLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventDelete, Key: key, Value: "_"}
}

func (l *PostgresTransactionalLogger) WritePut(key, value string) {
	l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

func (l *PostgresTransactionalLogger) Err() <-chan error {
	return l.errors
}

func CreatePostgresTransactionalLogger(dbParams PostgresDbParams) (TransactionalLogger, error) {
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s sslmode=disable password=%s", dbParams.host, dbParams.dbName, dbParams.user, dbParams.password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection %w", err)
	}

	logger := &PostgresTransactionalLogger{db: db}

	return logger, nil
}

func (l *PostgresTransactionalLogger) Run() <-chan error {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		query := "insert into transactions (event_type, key, value) values ($1, $2, $3)"
		for e := range events {

			_, err := l.db.Exec(query, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
				return
			}
		}
	}()

	return l.errors
}

func (l *PostgresTransactionalLogger) ReadAll() (<-chan Event, <-chan error) {

	allEvents := make(chan Event)
	errorChan := make(chan error, 1)

	go func() {
		var e Event

		defer close(allEvents)
		defer close(errorChan)

		query := "select id, event_type, key, value from transactions order by id"

		rows, err := l.db.Query(query)
		if err != nil {
			errorChan <- fmt.Errorf("sql query error %w", err)
			return
		}

		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&e.Id, &e.EventType, &e.Key, &e.Value)
			if err != nil {
				errorChan <- fmt.Errorf("sql reading error %w", err)
				return
			}

			allEvents <- e
		}

		err = rows.Err()
		if err != nil {
			errorChan <- fmt.Errorf("sql final reading error %w", err)
			return
		}
	}()

	return allEvents, errorChan
}
