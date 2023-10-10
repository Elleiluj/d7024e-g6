package router

import (
	"fmt"
	"kademlia/server"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter_print(t *testing.T) {
	fmt.Print("\n--------------------\n router.go\n--------------------\n")
}

func TestNewRouter(t *testing.T) {
	mockKademlia := &server.Kademlia{}
	router := NewRouter(mockKademlia)
	if router.kademlia != mockKademlia {
		t.Errorf("Expected Router to have kademlia instance %v, but got %v", mockKademlia, router.kademlia)
	} else {
		fmt.Println("NewRouter \tPASS")
	}
}

func TestRouterHandlers(t *testing.T) {
	fail := false
	me := server.NewKademliaNode("127.0.0.1:8150")
	me.JoinNetwork(&me.Me)
	network := server.NewNetwork(&me)
	go network.Listen(me.Me.Address)
	router := NewRouter(&me)

	router.DefineHandleFunc()
	go router.StartHTTP()

	// Create a test HTTP server
	testServer := httptest.NewServer(router.router)
	defer testServer.Close()

	t.Run("RootHandler", func(t *testing.T) {
		resp, err := http.Get(testServer.URL)
		if err != nil {
			t.Errorf("Failed to send GET request to /: %v", err)
			fail = true
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
			fail = true
		}
	})

	t.Run("ObjectPostHandler", func(t *testing.T) {
		payload := "test"
		// post
		resp, err := http.Post(testServer.URL+"/objects", "application/json", strings.NewReader(payload))
		if err != nil {
			t.Errorf("Failed to send POST request to /objects: %v", err)
			fail = true
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status code %d, but got %d", http.StatusCreated, resp.StatusCode)
			fail = true
		}

		locationHeader := resp.Header.Get("Location")
		if locationHeader == "" {
			t.Errorf("Location header is empty")
		}

		// get
		respGet, errGet := http.Get(testServer.URL + "/objects")
		if errGet != nil {
			t.Errorf("Failed to send POST request to /objects: %v", errGet)
			fail = true
			return
		}
		defer respGet.Body.Close()

		if respGet.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, respGet.StatusCode)
			fail = true
		}

		locationHeaderGet := resp.Header.Get("Location")
		if locationHeaderGet == "" {
			t.Errorf("Location header is empty")
			fail = true
		}

		// invalid
		respInvld, errInvld := http.Head(testServer.URL + "/objects")
		if errInvld != nil {
			t.Errorf("Failed to send POST request to /objects: %v", errInvld)
			fail = true
			return
		}
		defer respInvld.Body.Close()

		if respInvld.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status code %d, but got %d", http.StatusMethodNotAllowed, respInvld.StatusCode)
			fail = true
		}

		locationHeaderInvld := resp.Header.Get("Location")
		if locationHeaderInvld == "" {
			t.Errorf("Location header is empty")
			fail = true
		}
	})

	t.Run("ObjectGetHandler", func(t *testing.T) {
		hash := server.CreateHash("test")
		resp, err := http.Get(testServer.URL + "/objects/" + hash)
		if err != nil {
			t.Errorf("Failed to send GET request to /objects/%s: %v", hash, err)
			fail = true
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
			fail = true
		}

		invldHash := server.CreateHash("invalid")

		invldResp, invldErr := http.Get(testServer.URL + "/objects/" + invldHash)
		if invldErr != nil {
			t.Errorf("Got invalid data /objects/%s", hash)
			fail = true
			return
		}

		defer invldResp.Body.Close()
	})

	if !fail {
		fmt.Println("RouterHandlers \tPASS")
	}
}
