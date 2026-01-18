package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"time"
	"github.com/gorilla/mux"
	"github.com/hashicorp/raft"
	raftmdb "github.com/hashicorp/raft-mdb"
)

type kvFsm struct {
	db *map[string]string
}

type snapshotNoop struct{}

func (sn snapshotNoop) Persist(_ raft.SnapshotSink) error {
	return nil
}

func (sn snapshotNoop) Release() {
	return
}

func (kv *kvFsm) Apply(log *raft.Log) any {

	fmt.Println(log.Type)

	// Store key value pair in raft mdb
	return true
}

func (kv *kvFsm) Restore(rc io.ReadCloser) error {
	fmt.Println(rc)

	// Loop through and (re)store all key value pairs in raft mdb
	return nil
}

func (kv *kvFsm) Snapshot() (raft.FSMSnapshot, error) {
	return snapshotNoop{}, nil
}

func setupRaft(dir string, nodeId string, raftAddress string, kvFsm *kvFsm) (*raft.Raft, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeId)

	store, err := raftmdb.NewMDBStore(dir)
	if err != nil {
		return nil, fmt.Errorf("Could not create MDB store: %s", err)
	}

	snapshots, err := raft.NewFileSnapshotStore(path.Join(dir, "snapshot"), 2, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("Could not create snapshot store: %s", err)
	}

	tcpAddress, err := net.ResolveTCPAddr("tcp", raftAddress)
	if err != nil {
		return nil, fmt.Errorf("Could not resolve address: %s", err)
	}

	transport, err := raft.NewTCPTransport(raftAddress, tcpAddress, 10, time.Second*10, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("Could not create tcp transport: %s", err)
	}

	r, err := raft.NewRaft(config, kvFsm, store, store, snapshots, transport)
	if err != nil {
		return nil, fmt.Errorf("Could not create raft instance: %s", err)
	}

	return r, nil
}

type httpServer struct {
	r  *raft.Raft
	db *map[string]string
}

func (hs httpServer) get(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	fmt.Println(key)
}

func (hs httpServer) set(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Set key value pair")
}

func (hs httpServer) listNodes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List nodes")
}

func (hs httpServer) removeNode(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Println(id)
}

func main() {
	nodeId := os.Args[1]
	if nodeId == "" {
		log.Fatal("Missing required parameter: --node-id")
	}

	dataDir := "data"
	db := &map[string]string{}
	kv := &kvFsm{db}

	r, err := setupRaft(path.Join(dataDir, "raft" + nodeId), nodeId, "localhost:9090", kv)
	if err != nil {
		log.Fatal(err)
	}

	hs := httpServer{r, db}

	router := mux.NewRouter()
	router.HandleFunc("/{key}", hs.get).Methods("GET")
	router.HandleFunc("/", hs.set).Methods("PUT")
	router.HandleFunc("/nodes/list", hs.listNodes).Methods("GET")
	router.HandleFunc("/nodes/{id}", hs.removeNode).Methods("DELETE")

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
