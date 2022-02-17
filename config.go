package main

import (
	"strings"
)

var registry string
var stateManager string
var loggingPath string
var containerdPath string

func GetRegistry() string {
	return registry
}

func SetRegistry(reg string) {
	registry = reg
}

func GetStateManagerUrl() string {
	return strings.TrimRight(stateManager, "/")
}

func SetStateManagerUrl(url string) {
	stateManager = url
}

func GetLoggingPath() string {
	return loggingPath
}

func SetLoggingPath(path string) {
	loggingPath = path
}

func GetContainerdPath() string {
	return containerdPath
}

func SetContainerdPath(path string) {
	containerdPath = path
}
