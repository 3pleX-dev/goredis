package main

import (
	"fmt"
	"os"

	"github.com/3pleX-dev/goredis/store"

	"net/http"
)

func main() {
	store := store.NewStore();

	h1 := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}

		key := r.FormValue("key")
		if key == "" {
			http.Error(w, "key is required", http.StatusBadRequest)
			return
		}

		store.Set(key, r.FormValue("value"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "key is required", http.StatusBadRequest)
			return
		}

		value, ok := store.Get(key)
		if !ok {
			http.Error(w, "key not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(value+"\n"))
	}

	h3 := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "key not found", http.StatusBadRequest)
			return
		}

		store.Delete(key)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success\n"))
	}

	http.HandleFunc("/get", h2)
	http.HandleFunc("/set", h1)
	http.HandleFunc("/delete", h3)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintln(os.Stderr, "ERR starting HTTP server:", err)
	}
}