package main

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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

	err = Put(key, string(value))

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

	value, err := Get(key)

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

	err := Delete(key)

	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/key/{key}", PutHandler).Methods("PUT")
	r.HandleFunc("/v1/key/{key}", GetHandler).Methods("GET")
	r.HandleFunc("/v1/key/{key}", DeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))

	// Put("test", "TEST VALUE")
	// val, error := Get("test")
	// if error != nil {
	// 	fmt.Println(error)
	// 	return
	// }
	// fmt.Println(val)
}
