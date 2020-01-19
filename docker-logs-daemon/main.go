package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ahmetb/dlog"
	"github.com/eli-goldberg/docker-logs-daemon/storage"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

const daemonInterval = 5 * time.Second
const serverPort = "8080"
const defaultLabelFilter = "collect_logs=true"

func attachContainer(cli *client.Client, container types.Container, st *storage.Storage, done chan<- string) {
	defer func(c chan<- string, id string) {
		c <- id
	}(done, container.ID)

	reader, readerErr := cli.ContainerLogs(context.Background(), container.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})

	if readerErr != nil {
		fmt.Printf("Error reading logs from container %v: %v\n", container.ID, readerErr)
		return
	}

	defer func() {
		reader.Close()
	}()

	// The docker api adds extra "header" info at the beginning of each record
	// We use "dlog" to simply strip down those first few bytes, and get only the actual log message
	headerStrippedReader := dlog.NewReader(reader)
	writer, writerErr := (*st).Writer(container.ID)
	if writerErr != nil {
		fmt.Printf("Error writing to storage: %v\n", writerErr)
	}
	defer func() {
		writer.Close()
	}();

	io.Copy(writer, headerStrippedReader)
}


func createStorages() map[string]storage.Storage{
	storages := make(map[string]storage.Storage)

	fileStoragePath := os.Getenv("FILE_STORAGE");
	
	if (fileStoragePath == "") {
		fileStoragePath = storage.DefaultFileStoragePath
		fmt.Printf("No FILE_STORAGE env var specified, defaulting to \"%v\"\n", fileStoragePath)
	}

	if fileStorage, fileStorageErr := storage.NewFileStorage(fileStoragePath); fileStorageErr == nil {
		storages[storage.FileStorageType] = fileStorage
	}

	return storages
}

func main() {
	labelFilter := os.Getenv("LABEL_FILTER")
	if (labelFilter == "") {
		labelFilter = defaultLabelFilter
		fmt.Printf("No LABEL_FILTER env var specified, defaulting to \"%v\"\n", labelFilter)
	}

	storages := createStorages()

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	logStorage := storages[storage.FileStorageType]

	done := waitForTermination()
	set := NewContainerSet()
	go trackContainers(cli, set, labelFilter, &logStorage, daemonInterval)
	fmt.Println("Daemon running...")
	go serveAPI("", serverPort, set, &logStorage)
	fmt.Printf("Server listening on port %v...\n", serverPort)
	<-done
}

func serveAPI(serverHost string, serverPort string, set *ContainerSet, st *storage.Storage) {
	apiServer := NewLogsAPIServer(set, st)
	apiServer.serve(serverHost, serverPort)
}

func trackContainers(cli *client.Client, set *ContainerSet, labelToTrack string, st *storage.Storage, interval time.Duration) {
	filters := filters.NewArgs()
	filters.Add("label", labelToTrack)
	containerListOptions := types.ContainerListOptions{Filters: filters}

	for {
		containers, err := cli.ContainerList(context.Background(), containerListOptions)

		if err != nil {
			panic(err)
		}

		for _, container := range containers {
			if !set.Exists(container.ID) {
				fmt.Printf("Container detected: %v\n", container.ID)
				set.Add(container.ID)
				fmt.Printf("Attaching logs from container %v...\n", container.ID)

				done := make(chan string)
				go attachContainer(cli, container, st, done)
				go func(done <-chan string) {
					set.Remove(<-done)
					fmt.Printf("Container stopped: %v", container.ID)
				}(done)
			}
		}

		time.Sleep(interval)
	}
}

func waitForTermination() chan bool {
	done := make(chan bool)
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("User terminated")
		done <- true
	}()
	return done
}
