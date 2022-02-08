package main

import (
	"os"
	"strings"
)

func GetRegistry() string {
	return os.Getenv("REGISTRY")
}

func GetStateManagerUrl() string {
	return strings.TrimRight(os.Getenv("STATE_MANAGER"), "/")
}

func GetLoggingPath() string {
	return os.Getenv("LOGGING_PATH")
}

func GetContainerdPath() string {
	return os.Getenv("CONTAINERD_PATH")
}
