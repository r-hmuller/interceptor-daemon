package main

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/google/uuid"
	"os"
	"strconv"
	"time"
)

func GenerateSnapshot(service Snapshot) {
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

	for _, container := range containers {
		labels, _ := container.Labels(ctx)
		if labels["io.kubernetes.container.name"] == service.Container {
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

			//Here send to state manager
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
}
