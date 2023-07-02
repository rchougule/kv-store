package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/rchougule/kv-store/kvstore"
	"github.com/rchougule/kv-store/kvstore/inmemorybasic"
)

var inMemKvStore kvstore.KVStore = inmemorybasic.NewStore()

var pathAndHandlerMap = map[string]func(w http.ResponseWriter, r *http.Request){
	"/get": handleGet,
	"/put": handlePut,

	"/v1/get": handleGet,
	"/v1/put": handlePut,
}

func main() {
	mux := &http.ServeMux{}
	for path, handler := range pathAndHandlerMap {
		mux.HandleFunc(path, handler)
	}

	ctx := context.Background()

	server := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(listener net.Listener) context.Context {
			ctx = context.WithValue(ctx, "serveraddr", listener.Addr().String())
			return ctx
		},
	}

	fmt.Printf("listening on port: %s\n", server.Addr)
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("error: server closed")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	value, err := getKVStore(r.URL.Path).Get(key)
	if err != nil {
		w.WriteHeader(404)
		if _, err = io.WriteString(w, "error fetching the key's value"); err != nil {
			fmt.Printf("error responding the error string: %s\n", err)
		}
		return
	}

	marshalledValue, err := json.Marshal(value)
	if err != nil {
		w.WriteHeader(500)
		if _, err = io.WriteString(w, "failed marshalling the value"); err != nil {
			fmt.Printf("error responding the error string: %s\n", err)
		}
		return
	}

	if _, err = io.WriteString(w, string(marshalledValue)); err != nil {
		w.WriteHeader(500)
		return
	}
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		if _, err = io.WriteString(w, "could not read body"); err != nil {
			fmt.Printf("error responding the error string: %s\n", err)
		}
		return
	}

	var unmarshalledBody interface{}
	if err = json.Unmarshal(body, &unmarshalledBody); err != nil {
		w.WriteHeader(400)
		if _, err = io.WriteString(w, "could not unmarshall body"); err != nil {
			fmt.Printf("error responding the error string: %s\n", err)
		}
		return
	}

	for key, value := range unmarshalledBody.(map[string]interface{}) {
		if err = getKVStore(r.URL.Path).Put(key, value); err != nil {
			w.WriteHeader(500)
			if _, err = io.WriteString(w, "could not store the key value"); err != nil {
				fmt.Printf("error responding the error string: %s\n", err)
			}
			return
		}
	}

	w.WriteHeader(200)
}

func getKVStore(path string) kvstore.KVStore {
	version := strings.Split(path, "/")[0]
	switch version {
	case "v1":
		return inMemKvStore
	default:
		return inMemKvStore
	}
}
