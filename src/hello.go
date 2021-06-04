package main

import (
	"log"
)

var contextVars struct {
	store KeyValueStore
}

func main() {
	log.Println("start")
	use_file_logger := true

	var logger TransactionalLogger
	var err error

	if use_file_logger {
		logger, err = CreateFileTransactionalLogger("/usr/local/transactional_log_file")
	} else {
		logger, err = CreatePostgresTransactionalLogger(PostgresDbParams{
			dbName:   "postgres",
			host:     "localhost",
			user:     "postgres",
			password: "",
		})
	}

	if err != nil {
		log.Fatalf("Can't create transactional logger %w", err)
	}

	log.Println("transactinal logger created")
	store, err := CreateKeyValueStore(logger)
	if err != nil {
		log.Fatalf("Can't create key value store %w", err)
	}

	log.Println("key value store created")
	err = store.RestorePersistedState()
	if err != nil {
		log.Fatalf("Can't restore persistent state %w", err)
	}

	log.Println("persistent state restores")
	errorsChan := logger.Run()

	log.Println("transactional logger run")
	go func() {
		for loggerError := range errorsChan {
			log.Fatal(loggerError)
		}
	}()

	contextVars.store = store

	StartRestService(8080)
}
