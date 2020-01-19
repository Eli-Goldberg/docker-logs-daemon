package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/eli-goldberg/docker-logs-daemon/storage"

	"github.com/gorilla/websocket"
)

type listResponse struct {
	Streams []string `json:"streams"`
}

type logsRequest struct {
	StreamID string
}

// LogsAPIServer manages the api
type LogsAPIServer struct {
	containerSet *ContainerSet
	logStorage   *storage.Storage
}

// NewLogsAPIServer creates a new LosApiServer
func NewLogsAPIServer(set *ContainerSet, st *storage.Storage) *LogsAPIServer {
	return &LogsAPIServer{
		containerSet: set,
		logStorage:   st,
	}
}

func (LogsAPIServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	content := "{ \"ok\": true }"
	fmt.Fprintf(w, "%s", content)
}

func (s LogsAPIServer) handleList(w http.ResponseWriter, r *http.Request) {
	streams, err := (*s.logStorage).List()
	if err != nil {
		fmt.Printf("Error: Could not read log, %v", err)
		http.Error(w, "Could not send message", http.StatusInternalServerError)
		return
	}
	
	response, jsonErr := json.MarshalIndent(listResponse{Streams: streams}, "", "  ")
	if (jsonErr != nil) {
		fmt.Printf("Error: Unknown Error, %v", err)
		http.Error(w, "Uknown error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")

	fmt.Fprintf(w, string(response))
}

func (s LogsAPIServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	defer func() {
		conn.Close()
	}()

	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	lReq := logsRequest{}
	if requestErr := conn.ReadJSON(&lReq); requestErr != nil {
		fmt.Printf("Error: Could not read request")
		return
	}
	fmt.Printf("Log stream requested for container %v", lReq.StreamID)

	reader, readerErr := (*s.logStorage).Reader(lReq.StreamID)
	if readerErr != nil {
		fmt.Printf("Error: Could not read log, %v", readerErr)
		http.Error(w, "Could not send message", http.StatusInternalServerError)
		return
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		writeErr := conn.WriteMessage(websocket.TextMessage, scanner.Bytes())
		if writeErr != nil {
			fmt.Printf("Error: could not send message, %v", writeErr)
		}
	}
}

func (s LogsAPIServer) serve(host, port string) {

	http.HandleFunc("/status", s.handleStatus)
	http.HandleFunc("/list", s.handleList)
	http.HandleFunc("/ws", s.wsHandler)
	if err := http.ListenAndServe(host+":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
