package main

import (
	"os"
)

func GetRegistry() string {
	return os.Getenv("REGISTRY")
}

func GetStateManagerUrl() string {
	return os.Getenv("STATE_MANAGER")
}

func GetLoggingPath() string {
	return os.Getenv("LOGGING_PATH")
}

func GetContainerdPath() string {
	return os.Getenv("CONTAINERD_PATH")
}
