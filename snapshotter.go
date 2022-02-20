package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GenerateSnapshot(service Snapshot) string {
	startTime := time.Now().Unix()
	// create a new client connected to the default socket path for containerd
	client, err := containerd.New(GetContainerdPath())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// create a new context with an "example" namespace
	ctx := namespaces.WithNamespace(context.Background(), service.Namespace)

	containers, err := client.Containers(ctx)
	if err != nil {
		panic(err)
	}

	port := ""
	for _, container := range containers {
		labels, _ := container.Labels(ctx)
		if labels["io.kubernetes.pod.name"] == service.Container {
			task, err := container.Task(ctx, cio.NewAttach(cio.WithStdio))
			if err != nil {
				panic(err)
			}
			checkpoint, err := task.Checkpoint(ctx)
			if err != nil {
				panic(err)
			}

			registry := GetRegistry()
			containerSnapshotVersion := registry + ":" + uuid.NewString()
			err = client.Push(ctx, containerSnapshotVersion, checkpoint.Target())
			if err != nil {
				panic(err)
			}

			postBody, _ := json.Marshal(map[string]string{
				"Service": service.Service,
			})
			responseBody := bytes.NewBuffer(postBody)
			//Leverage Go's HTTP Post function to make request
			fullUrl := GetStateManagerUrl() + "/"
			resp, err := http.Post(fullUrl, "application/json", responseBody)
			//Handle Error
			if err != nil {
				log.Fatalf("An Error Occured %v", err)
			}
			defer resp.Body.Close()

			port, err = DeployNewContainer(service.YamString, containerSnapshotVersion, service.Service)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	endTime := time.Now().Unix()
	deltaTime := endTime - startTime
	loggingPath := GetLoggingPath()
	f, err := os.OpenFile(loggingPath+"/checkpoint.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	if _, err = f.WriteString(strconv.FormatInt(deltaTime, 10)); err != nil {
		panic(err)
	}

	return port
}
