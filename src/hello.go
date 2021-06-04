package main

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var contextVars struct {
	store KeyValueStore
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	err = contextVars.store.Put(key, string(value))

	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := contextVars.store.Get(key)

	if err != nil {
		if errors.Is(err, ErrorNoSuchKey) {
			http.Error(w,
				err.Error(),
				http.StatusNotFound)
		} else {
			http.Error(w,
				err.Error(),
				http.StatusInternalServerError)
		}
		return
	}

	w.Write([]byte(value))
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := contextVars.store.Delete(key)

	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
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

	r := mux.NewRouter()
	r.HandleFunc("/v1/key/{key}", PutHandler).Methods("PUT")
	r.HandleFunc("/v1/key/{key}", GetHandler).Methods("GET")
	r.HandleFunc("/v1/key/{key}", DeleteHandler).Methods("DELETE")

	log.Println("ready to run webserver")
	log.Fatal(http.ListenAndServe(":8080", r))
}
