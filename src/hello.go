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
	logger, err := CreateFileTransactionalLogger("transactional_log_file")
	if err != nil {
		log.Fatal("Can't create transactional logger")
	}

	contextVars.store, err = CreateKeyValueStore(logger)
	if err != nil {
		log.Fatal("Can't create key value store")
	}

	r := mux.NewRouter()
	r.HandleFunc("/v1/key/{key}", PutHandler).Methods("PUT")
	r.HandleFunc("/v1/key/{key}", GetHandler).Methods("GET")
	r.HandleFunc("/v1/key/{key}", DeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
