package router

import (
	"fmt"
	"io"
	"kademlia/server"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	kademlia *server.Kademlia
	router   *mux.Router
}

func NewRouter(kademlia *server.Kademlia) *Router {
	return &Router{kademlia: kademlia, router: mux.NewRouter()}
}

func (router *Router) DefineHandleFunc() {

	router.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello! from node with ip: "+router.kademlia.Me.Address)
	})

	router.router.HandleFunc("/objects", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {
			data, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()
			error := router.kademlia.Store([]byte(data))
			if error != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error: Failed to store the data"))
			} else {
				w.Header().Set("Location", "/objects/"+server.CreateHash(string(data)))
				w.WriteHeader(http.StatusCreated)
				w.Write(data)
			}
		} else if r.Method == "GET" {
			fmt.Fprint(w, "Hello! from node with ip: "+router.kademlia.Me.Address)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	router.router.HandleFunc("/objects/{hash}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hash := vars["hash"]
		node, value := router.kademlia.LookupData(hash)
		if value == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to get data with hash: " + hash))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(value + " found at node: " + node.Address))
		}
	}).Methods("GET")
}

func (router *Router) StartHTTP() {
	http.Handle("/", router.router)
	http.ListenAndServe(":8080", nil)
}
